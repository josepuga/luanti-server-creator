@echo off

REM Server graceful shutdown.
REM By José Puga 2025.

docker exec -it "%server%" ^
    sh -c "kill -SIGTERM \$(pgrep minetest)"