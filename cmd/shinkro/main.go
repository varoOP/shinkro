package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/pflag"
	"golang.org/x/term"

	"github.com/varoOP/shinkro/internal/anime"
	"github.com/varoOP/shinkro/internal/animeupdate"
	"github.com/varoOP/shinkro/internal/api"
	"github.com/varoOP/shinkro/internal/auth"
	"github.com/varoOP/shinkro/internal/config"
	"github.com/varoOP/shinkro/internal/database"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/internal/filesystem"
	"github.com/varoOP/shinkro/internal/http"
	"github.com/varoOP/shinkro/internal/logger"
	"github.com/varoOP/shinkro/internal/malauth"
	"github.com/varoOP/shinkro/internal/mapping"
	"github.com/varoOP/shinkro/internal/notification"
	"github.com/varoOP/shinkro/internal/plex"
	"github.com/varoOP/shinkro/internal/plexsettings"
	"github.com/varoOP/shinkro/internal/plexstatus"
	"github.com/varoOP/shinkro/internal/server"
	"github.com/varoOP/shinkro/internal/user"
)

const usage = `shinkro
Sync your Anime watch status in Plex to myanimelist.net!
Usage:
    shinkro --config <path to shinkro configuration directory>           Run shinkro
    shinkro setup [--config <dir>]                                      Setup new config, DB, and admin user
    shinkro --config=<dir> change-password <username>                   Change password for user
    shinkro version                                                     Print version info
    shinkro help                                                        Show this help message
`

func init() {
	pflag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), usage)
	}
}

var (
	version = "dev"
	commit  = ""
	date    = ""
)

func main() {
	var configPath string

	d, err := homedir.Dir()
	if err != nil {
		fmt.Fprint(flag.CommandLine.Output(), "FATAL: Unable to get home directory")
		os.Exit(1)
	}

	d = filepath.Join(d, ".config", "shinkro")
	pflag.StringVar(&configPath, "config", d, "path to configuration")
	pflag.Parse()
	configPath, err = homedir.Expand(configPath)
	if err != nil {
		fmt.Fprint(flag.CommandLine.Output(), "FATAL: Unable to expand configuration path")
		os.Exit(1)
	}

	switch cmd := pflag.Arg(0); cmd {
	case "":
		cfg := config.NewConfig(configPath, version)
		log := logger.NewLogger(cfg.Config).Logger
		go cfg.DynamicReload(log)
		db := database.NewDB(configPath, &log)

		log.Info().Msg("Starting shinkro")
		log.Info().Msgf("Version: %s", version)
		log.Info().Msgf("Commit: %s", commit)
		log.Info().Msgf("Build date: %s", date)
		log.Info().Msgf("Base URL: %s", cfg.Config.BaseUrl)
		log.Info().Msgf("Log-level: %s", cfg.Config.LogLevel)

		err = db.Migrate()
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}

		var animeRepo = database.NewAnimeRepo(log, db)
		var animeUpdateRepo = database.NewAnimeUpdateRepo(log, db)
		var plexRepo = database.NewPlexRepo(log, db)
		var plexSettingsRepo = database.NewPlexSettingsRepo(log, db)
		var malauthRepo = database.NewMalAuthRepo(log, db)
		var userRepo = database.NewUserRepo(log, db)
		var apiRepo = database.NewAPIRepo(log, db)
		var mappingRepo = database.NewMappingRepo(log, db)
		var notificationRepo = database.NewNotificationRepo(log, db)
		var plexStatusRepo = database.NewPlexStatusRepo(log, db)

		var animeService = anime.NewService(log, animeRepo)
		var malauthService = malauth.NewService(cfg.Config, log, malauthRepo)
		var mapService = mapping.NewService(log, mappingRepo)
		var animeUpdateService = animeupdate.NewService(log, animeUpdateRepo, animeService, mapService, malauthService)
		var plexSettingsService = plexsettings.NewService(cfg.Config, log, plexSettingsRepo)
		var notificationService = notification.NewService(log, notificationRepo)
		var plexStatusService = plexstatus.NewService(log, plexStatusRepo)
		var plexService = plex.NewService(log, plexSettingsService, plexRepo, animeService, mapService, malauthService, animeUpdateService, notificationService, plexStatusService)
		var userService = user.NewService(userRepo, log)
		var authService = auth.NewService(log, userService)
		var apiService = api.NewService(log, apiRepo)
		var fsService = filesystem.NewService(cfg.Config, log)

		srv := server.NewServer(log, cfg.Config, animeService, mapService)
		if err := srv.Start(); err != nil {
			log.Fatal().Stack().Err(err).Msg("could not start server")
			return
		}

		errorChannel := make(chan error)

		go func() {
			httpServer := http.NewServer(
				log,
				cfg,
				db,
				version,
				commit,
				date,
				plexService,
				plexSettingsService,
				malauthService,
				apiService,
				authService,
				mapService,
				fsService,
				notificationService,
				animeUpdateService,
			)
			errorChannel <- httpServer.Open()
		}()

		sigchnl := make(chan os.Signal, 1)
		signal.Notify(sigchnl, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

		for sig := range sigchnl {
			log.Info().Msgf("received signal: %v, shutting down server.", sig)
			if err := db.Close(); err != nil {
				log.Error().Err(err).Msg("failed to close the database connection properly")
				os.Exit(1)
			}
			os.Exit(0)
		}

	case "setup":
		fmt.Println("--- Shinkro Setup ---")
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Config directory [%s]: ", configPath)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			configPath = input
		}
		cfg := config.NewConfig(configPath, version)
		log := logger.NewLogger(cfg.Config).Logger
		if err := cfg.WriteConfig(configPath, "config.toml"); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write config: %v\n", err)
			os.Exit(1)
		}
		db := database.NewDB(configPath, &log)
		if err := db.Migrate(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to migrate DB: %v\n", err)
			os.Exit(1)
		}
		userRepo := database.NewUserRepo(log, db)
		userService := user.NewService(userRepo, log)
		authService := auth.NewService(log, userService)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		fmt.Print("Enter username: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)
		var password, password2 string
		for {
			fmt.Print("Enter password: ")
			pw1, _ := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			fmt.Print("Confirm password: ")
			pw2, _ := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			password = strings.TrimSpace(string(pw1))
			password2 = strings.TrimSpace(string(pw2))
			if password == password2 && password != "" {
				break
			}
			fmt.Println("Passwords do not match or are empty. Try again.")
		}
		if err := authService.CreateUser(ctx, domain.CreateUserRequest{Username: username, Password: password}); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create user: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Setup complete! Config, DB, and admin user created.")
		os.Exit(0)

	case "change-password":
		if pflag.NArg() < 2 {
			fmt.Fprintln(os.Stderr, "Usage: shinkro --config=<dir> change-password <username>")
			os.Exit(1)
		}
		username := pflag.Arg(1)
		cfg := config.NewConfig(configPath, version)
		log := logger.NewLogger(cfg.Config).Logger
		db := database.NewDB(configPath, &log)
		userRepo := database.NewUserRepo(log, db)
		userService := user.NewService(userRepo, log)
		authService := auth.NewService(log, userService)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		fmt.Printf("Changing password for user '%s'\n", username)
		fmt.Print("Enter new password: ")
		pwNew, _ := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		fmt.Print("Confirm new password: ")
		pwNew2, _ := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if strings.TrimSpace(string(pwNew)) != strings.TrimSpace(string(pwNew2)) || strings.TrimSpace(string(pwNew)) == "" {
			fmt.Fprintln(os.Stderr, "New passwords do not match or are empty.")
			os.Exit(1)
		}
		newPassword := strings.TrimSpace(string(pwNew))
		if err := authService.ResetPassword(ctx, username, newPassword); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to change password: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Password updated successfully.")
		os.Exit(0)

	case "version":
		fmt.Fprintln(flag.CommandLine.Output(), "Version:", version)
		fmt.Fprintln(flag.CommandLine.Output(), "Commit:", commit)
		fmt.Fprintln(flag.CommandLine.Output(), "Build Date:", date)

	default:
		pflag.Usage()
		if cmd != "help" {
			os.Exit(0)
		}
	}
}
