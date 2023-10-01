<h1 align="center">
  <img alt="shinkro logo" src=".github/images/logo.png" width="160px"/><br/>
  shinkro
</h1>

<p align="center">An application to sync Plex watch status to myanimelist.</p>

<p align="center"><img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/varoOP/shinkro?style=for-the-badge">&nbsp;<img alt="GitHub all releases" src="https://img.shields.io/github/downloads/varoOP/shinkro/total?style=for-the-badge">&nbsp;<img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/varoOP/shinkro/release.yml?style=for-the-badge"></p>

## Documentation

Installation guide and documentation can be found at https://docs.shinkro.com

## Key features

- Support for multiple metadata agents, including default Plex agent, HAMA, and MyAnimeList.bundle.
- Live updates to myanimelist as soon as you watch or rate in Plex.
- Powerful anime-id mapping support, make custom maps or use the community mapping.
- Built on Go making shinkro lightweight and perfect for supporting multiple platforms (Linux, FreeBSD,
  Windows, macOS) on different architectures. (e.g. x86, ARM)
- Discord Notifications.
- Base path / Subfolder (and subdomain) support for convenient reverse-proxy support.

Available methods to use shinkro

- Official Plex Webhook (Requires Plex Pass) (recommended)
- Tautulli (Limitation: cannot sync ratings)

## What is shinkro?

If you use both Plex and Myanimelist to watch and track your anime, you know how mundande and boring it is to have to update your myanimelist manually after watching an anime on your Plex server.

shinkro enables you to sync your Plex ratings and watch status for anime to myanimelist.net.

## Installation

Full installation guide and documentation can be found at https://docs.shinkro.com

### Windows

Check the windows setup guide [here](https://autobrr.com/installation/windows)

### Linux generic

Download the latest release,

```bash
wget $(curl -s https://api.github.com/repos/varoOP/shinkro/releases/latest | grep download | grep linux_x86_64 | cut -d\" -f4)
```

#### Unpack

Run with `root` or `sudo`. If you do not have root, or are on a shared system, place the binaries somewhere in your home
directory like `~/.bin`.

```bash
tar -C /usr/local/bin -xzf shinkro*.tar.gz
```

This will extract `shinkro` to `/usr/local/bin`.
Note: If the command fails, prefix it with `sudo ` and re-run again.

#### Systemd (Recommended)

On Linux-based systems, it is recommended to run shinkro as a sort of service with auto-restarting capabilities, in
order to account for potential downtime. The most common way is to do it via systemd.

You will need to create a service file in `/etc/systemd/system/` called `shinkro.service`.

```bash
touch /etc/systemd/system/shinkro@.service
```

Then place the following content inside the file (e.g. via nano/vim/ed):

```systemd title="/etc/systemd/system/shinkro@.service"
[Unit]
Description=shinkro service for %i
After=syslog.target network-online.target

[Service]
Type=simple
User=%i
Group=%i
ExecStart=/usr/bin/shinkro --config=/home/%i/.config/shinkro

[Install]
WantedBy=multi-user.target
```

Start the service. Enable will make it startup on reboot.

```bash
systemctl enable -q --now --user shinkro@$USER
```

By default, the configuration is set to listen on `127.0.0.1`. While shinkro works fine as is exposed to the internet,
it is recommended to use a reverse proxy
like [nginx](https://docs.shinkro.com/installation/linux#nginx)

If you are not running a reverse proxy change `Host` in the `config.toml` to `0.0.0.0`.

## Community

Come join us on [Discord](https://discord.gg/ZkYdfNgbAT)!

## License

* [MIT](https://mit-license.org/)
* Copyright 2022-2023