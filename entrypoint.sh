#!/bin/sh

# Check if PGID and PUID are set, and create the group/user accordingly.
if [ -n "$PGID" ] && [ -n "$PUID" ]; then
    addgroup -g "$PGID" shinkrogrp
    adduser -D -u "$PUID" -G shinkrogrp shinkrousr
fi

# Only generate config.toml from the template if it doesn't already exist
if [ ! -f /config/config.toml ]; then

    # Transform the AnimeLibraries env var
    if [ -n "$ANIME_LIBRARIES" ]; then
        ANIME_LIBRARIES=$(echo "$ANIME_LIBRARIES" | sed 's/,/","/g')
        export ANIME_LIBRARIES="\"$ANIME_LIBRARIES\""
    fi

    SHINKRO_APIKEY=$(shinkro genkey)
    export SHINKRO_APIKEY="$SHINKRO_APIKEY"

    # Here, set the proper user for envsubst depending on if PUID/PGID are set
    if [ -n "$PUID" ] && [ -n "$PGID" ]; then
        gosu shinkrousr envsubst </app/config.toml.template >/config/config.toml
    else
        envsubst </app/config.toml.template >/config/config.toml
    fi
fi

# Change ownership of /config after generating config.toml
if [ -n "$PGID" ] && [ -n "$PUID" ]; then
    chown -R "$PUID":"$PGID" /config
fi

# Execute the binary with arguments
if [ -n "$PUID" ] && [ -n "$PGID" ]; then
    exec gosu shinkrousr /usr/local/bin/shinkro --config /config "$@"
else
    exec /usr/local/bin/shinkro --config /config "$@"
fi
