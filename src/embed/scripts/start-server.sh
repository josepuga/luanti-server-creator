#!/bin/bash

# Starts the Server
# By José Puga 2025.

docker run -it --rm -d \
    --net host \
    -p 30000:30000 \
    --name "%server%" \
    -v "./data:/root/.minetest:Z" \
    docker.io/josepuga/luanti-server:latest
