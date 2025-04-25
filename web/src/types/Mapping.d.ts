export interface Mapping {
    tvdb_enabled: boolean;
    tmdb_enabled: boolean;
    tvdb_path: string;
    tmdb_path: string;
}

export interface ValidateMap {
    yamlPath: string;
    isTVDB: boolean;
}