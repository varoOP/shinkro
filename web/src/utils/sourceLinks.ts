/**
 * Generates external links for different source databases
 */

export type SourceDB = "TVDB" | "TMDB" | "AniDB" | "MAL";

export const getSourceLink = (sourceDB: string, sourceID: number): string | null => {
    if (!sourceDB || !sourceID) {
        return null;
    }
    
    const lowerSourceDB = sourceDB.toLowerCase().trim();
    
    switch (lowerSourceDB) {
        case "tvdb":
            return `https://thetvdb.com/index.php?tab=series&id=${sourceID}`;
        case "tmdb":
            // For anime, this is typically a TV show
            return `https://www.themoviedb.org/tv/${sourceID}`;
        case "anidb":
            return `https://anidb.net/anime/${sourceID}`;
        case "myanimelist":
        case "mal":
            return `https://myanimelist.net/anime/${sourceID}`;
        default:
            console.warn(`Unknown sourceDB: "${sourceDB}" (normalized: "${lowerSourceDB}")`);
            return null;
    }
};

export const getMALLink = (malId: number): string => {
    return `https://myanimelist.net/anime/${malId}`;
};

