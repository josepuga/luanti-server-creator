@echo off

REM Starts the Server
REM By José Puga 2025. GPL3

docker run -it --rm -d ^
    --net host ^
    -p 30000:30000 ^
    --name "%server%" ^
    -v "%cd%\data:/root/.minetest" ^
    docker.io/josepuga/luanti-server:%tag%
