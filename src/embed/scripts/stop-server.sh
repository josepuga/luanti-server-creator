#!/bin/bash

# Server graceful shutdown.
# By José Puga 2025. GPL3

docker exec -it "%server%" \
    bash -c "kill -SIGTERM \$(pgrep luanti)"