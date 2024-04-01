#!/bin/sh

# Function to check if a group with the given GID exists and return its name
get_group_name_by_gid() {
    getent group "$1" | cut -d: -f1
}

# Function to check if a user with the given UID exists and return its name
get_user_name_by_uid() {
    getent passwd "$1" | cut -d: -f1
}

# Attempt to create the group if PGID is provided and does not already exist
if [ -n "$PGID" ]; then
    GROUP_NAME=$(get_group_name_by_gid "$PGID")
    if [ -z "$GROUP_NAME" ]; then
        GROUP_NAME="shinkrogrp"
        addgroup -g "$PGID" "$GROUP_NAME"
    fi
else
    GROUP_NAME="shinkrogrp"
    # Optionally create a default group if PGID is not set
fi

# Attempt to create the user if PUID is provided and does not already exist
if [ -n "$PUID" ]; then
    USER_NAME=$(get_user_name_by_uid "$PUID")
    if [ -z "$USER_NAME" ]; then
        USER_NAME="shinkrousr"
        adduser -D -u "$PUID" -G "$GROUP_NAME" "$USER_NAME"
    fi
else
    USER_NAME="shinkrousr"
    # Optionally create a default user if PUID is not set
fi

# Transform and export the ANIME_LIBRARIES environment variable
if [ -n "$ANIME_LIBRARIES" ]; then
    ANIME_LIBRARIES=$(echo "$ANIME_LIBRARIES" | sed 's/,/","/g')
    export ANIME_LIBRARIES="\"$ANIME_LIBRARIES\""
fi

# Generate and export the SHINKRO_APIKEY
SHINKRO_APIKEY=$(shinkro genkey)
export SHINKRO_APIKEY="$SHINKRO_APIKEY"

# Only generate config.toml from the template if it doesn't already exist
if [ ! -f /config/config.toml ]; then
    if [ -n "$PUID" ] && [ -n "$PGID" ]; then
        gosu "$USER_NAME" envsubst </app/config.toml.template >/config/config.toml
    else
        envsubst </app/config.toml.template >/config/config.toml
    fi
fi

# Change ownership of /config, considering existing or newly created user/group
chown -R "$PUID":"$PGID" /config

# Execute the main process
if [ -n "$PUID" ] && [ -n "$PGID" ]; then
    exec gosu "$USER_NAME" /usr/local/bin/shinkro --config /config "$@"
else
    exec /usr/local/bin/shinkro --config /config "$@"
fi
