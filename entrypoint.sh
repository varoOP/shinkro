#!/bin/sh

# Only generate config.toml from the template if it doesn't already exist
if [ ! -f /config/config.toml ]; then

    # Transform the AnimeLibraries env var from "Library1","Library2","Library3" format to TOML array format
    if [ -n "$AnimeLibraries" ]; then
        AnimeLibraries=$(echo "$AnimeLibraries" | sed 's/,/","/g')
        export AnimeLibraries="\"$AnimeLibraries\""
    fi

    ApiKey=$(shinkro genkey)
    export ApiKey="$ApiKey"

    envsubst </app/config.toml.template >/config/config.toml
fi

# Execute the binary with arguments
exec /usr/local/bin/shinkro --config /config "$@"
