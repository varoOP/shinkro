package domain

type CommunityMapUrls string

const (
	CommunityMapTVDB CommunityMapUrls = "https://github.com/varoOP/shinkro-mapping/raw/main/tvdb-mal.yaml"
	TVDBSchema       CommunityMapUrls = "https://github.com/varoOP/shinkro-mapping/raw/main/.github/schema-tvdb.json"
	CommunityMapTMDB CommunityMapUrls = "https://github.com/varoOP/shinkro-mapping/raw/main/tmdb-mal.yaml"
	TMDBSchema       CommunityMapUrls = "https://github.com/varoOP/shinkro-mapping/raw/main/.github/schema-tmdb.json"
)

type AnimeMap struct {
	AnimeTVShows *AnimeTVShows
	AnimeMovies  *AnimeMovies
}

type AnimeMapDetails struct {
	Malid int
	Start int
}

type AnimeTVShows struct {
	Anime []AnimeTV `yaml:"AnimeMap" json:"AnimeMap"`
}

type AnimeTV struct {
	Malid        int            `yaml:"malid" json:"malid"`
	Title        string         `yaml:"title" json:"title"`
	Type         string         `yaml:"type" json:"type"`
	Tvdbid       int            `yaml:"tvdbid" json:"tvdbid"`
	TvdbSeason   int            `yaml:"tvdbseason" json:"tvdbseason"`
	Start        int            `yaml:"start" json:"start"`
	UseMapping   bool           `yaml:"useMapping" json:"useMapping"`
	AnimeMapping []AnimeMapping `yaml:"animeMapping" json:"animeMapping"`
}

type AnimeMapping struct {
	TvdbSeason int `yaml:"tvdbseason" json:"tvdbseason"`
	Start      int `yaml:"start" json:"start"`
}

type AnimeMovies struct {
	AnimeMovie []AnimeMovie `yaml:"animeMovies" json:"animeMovies"`
}

type AnimeMovie struct {
	MainTitle string `yaml:"mainTitle" json:"mainTitle"`
	TMDBID    int    `yaml:"tmdbid" json:"tmdbid"`
	MALID     int    `yaml:"malid" json:"malid"`
}

// func NewAnimeMaps(cfg *Config) (*AnimeTVDBMap, *AnimeMovies, error) {
// 	cfg.LocalMapsExist()
// 	err := loadCommunityMaps(cfg)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	err = loadLocalMaps(cfg)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return cfg.TVDBMalMap, cfg.TMDBMalMap, nil
// }

func (s *AnimeTVShows) CheckMap(tvdbid, tvdbseason, ep int) (bool, *AnimeTV) {
	candidates := s.findMatchingAnime(tvdbid, tvdbseason)
	if len(candidates) == 1 {
		return true, &candidates[0]
	} else if len(candidates) > 1 {
		anime := s.findBestMatchingAnime(ep, candidates)
		return true, &anime
	}

	return false, nil
}

func (s *AnimeTVShows) findMatchingAnime(tvdbid, tvdbseason int) []AnimeTV {
	var matchingAnime []AnimeTV
	for _, anime := range s.Anime {
		if tvdbid != anime.Tvdbid {
			continue
		}

		if !anime.UseMapping && tvdbseason == anime.TvdbSeason {
			matchingAnime = append(matchingAnime, anime)
			continue
		}

		matchingMappedAnime := s.findMatchingMappedAnime(anime, tvdbseason)
		if matchingMappedAnime != nil {
			return []AnimeTV{*matchingMappedAnime}
		}
	}

	return matchingAnime
}

func (s *AnimeTVShows) findMatchingMappedAnime(anime AnimeTV, tvdbseason int) *AnimeTV {
	if !anime.UseMapping {
		return nil
	}

	for _, animeMap := range anime.AnimeMapping {
		if tvdbseason == animeMap.TvdbSeason {
			anime.TvdbSeason = animeMap.TvdbSeason
			anime.Start = animeMap.Start
			return &anime
		}
	}

	return nil
}

func (s *AnimeTVShows) findBestMatchingAnime(ep int, candidates []AnimeTV) AnimeTV {
	var anime AnimeTV
	largestStart := 0
	for _, v := range candidates {
		if ep >= v.Start && v.Start >= largestStart {
			largestStart = v.Start
			anime = v
		}
	}

	return anime
}

func (am *AnimeMovies) CheckMap(tmdbid int) (bool, *AnimeMovie) {
	for _, animeMovie := range am.AnimeMovie {
		if animeMovie.TMDBID == tmdbid {
			return true, &animeMovie
		}
	}

	return false, nil
}

// func loadCommunityMaps(cfg *Config) error {
// 	if !cfg.CustomMapTVDB {
// 		s := &AnimeTVDBMap{}
// 		respTVDB, err := http.Get(string(CommunityMapTVDB))
// 		if err != nil {
// 			return err
// 		}

// 		err = readYamlHTTP(respTVDB, s)
// 		if err != nil {
// 			return err
// 		}

// 		cfg.TVDBMalMap = s
// 	}

// 	if !cfg.CustomMapTMDB {
// 		am := &AnimeMovies{}
// 		respTMDB, err := http.Get(string(CommunityMapTMDB))
// 		if err != nil {
// 			return err
// 		}

// 		err = readYamlHTTP(respTMDB, am)
// 		if err != nil {
// 			return err
// 		}

// 		cfg.TMDBMalMap = am
// 	}

// 	return nil
// }

// func loadLocalMaps(cfg *Config) error {
// 	if cfg.CustomMapTVDB {
// 		s := &AnimeTVDBMap{}
// 		fTVDB, err := os.Open(cfg.CustomMapTVDBPath)
// 		if err != nil {
// 			return err
// 		}

// 		err = readYamlFile(fTVDB, s)
// 		if err != nil {
// 			return err
// 		}

// 		cfg.TVDBMalMap = s
// 	}

// 	if cfg.CustomMapTMDB {
// 		am := &AnimeMovies{}
// 		fTMDB, err := os.Open(cfg.CustomMapTMDBPath)
// 		if err != nil {
// 			return err
// 		}

// 		err = readYamlFile(fTMDB, am)
// 		if err != nil {
// 			return err
// 		}

// 		cfg.TMDBMalMap = am
// 	}

// 	return nil
// }

// func ChecklocalMaps(cfg *Config) (error, bool) {
// 	loadLocalMaps(cfg)
// 	localMapLoaded := false
// 	if cfg.CustomMapTVDB {
// 		if err := validateYaml(string(TVDBSchema), cfg.TVDBMalMap); err != nil {
// 			return err, false
// 		}

// 		localMapLoaded = true
// 	}

// 	if cfg.CustomMapTMDB {
// 		if err := validateYaml(string(TMDBSchema), cfg.TMDBMalMap); err != nil {
// 			return err, false
// 		}

// 		localMapLoaded = true
// 	}

// 	return nil, localMapLoaded
// }

// func validateYaml(schema string, yaml any) error {
// 	compiler := jsonschema.NewCompiler()
// 	sch, err := compiler.Compile(schema)
// 	if err != nil {
// 		return err
// 	}

// 	var v interface{}
// 	b, err := json.Marshal(yaml)
// 	if err != nil {
// 		return err
// 	}

// 	err = json.Unmarshal(b, &v)
// 	if err != nil {
// 		return err
// 	}

// 	if err := sch.Validate(v); err != nil {
// 		return err
// 	}

// 	return nil
// }
