#!/bin/sh

# Ensure /app/bin exists
mkdir -p /app/bin

# Move the executable inside /app/bin if not already there
if [ ! -f /app/bin/turplecubes ]; then
    mv /app/turplecubes /app/bin/turplecubes
    chmod +x /app/bin/turplecubes
fi

# Ensure required directories exist
mkdir -p /app/bin/turplecubes_proxy
mkdir -p /app/bin/turplecubes_volumes
mkdir -p /app/bin/turplecubes_conf

# Ensure images.json exists inside turplecubes_conf
if [ ! -f /app/bin/turplecubes_conf/images.json ]; then
    echo '{}' > /app/bin/turplecubes_conf/images.json
fi

# Start the application
exec /app/bin/turplecubes
