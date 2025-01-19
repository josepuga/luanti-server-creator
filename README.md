# Luanti Server Creator

To create a Luanti server in Docker with just 2 clicks.

## How to Use
On the machine that will act as the server, you need to have Docker installed.

Windows users: For this Operating System, the application is in a **completely experimental phase**. I don't use Windows, and although I've tried to program it as multiplatform and compiled it for Windows, I haven't been able to test it personally, so I can't guarantee it will work.

### Download
If you don’t want to compile it yourself, download the executables found in the `bin/` directory of the repository. Use `servercreator` for Linux and `servercreator.exe` for Windows.

### From the Game Client
1. (Optional) Create a new world of the desired game type (minetest, mineclone, voxelibre...). **Important**: Although this step is optional, it is recommended to launch the server with a new world.
2. (Optional) Enable the mods you want to use.

### Launch the Application
It’s recommended to do this from the command line to see any error messages.

*NOTE*: The first time it’s launched, a `config.ini` file will be created where you can specify the directory where servers will be created. By default, they are created inside `<app path>/servers`. If you want to change it, close the application, modify the value, and reopen it.

1. In the dropdown menu, select your world.
   
   ![](./screenshot-01.png "World selection")

2. You can select the server version to use or leave it as the latest available (latest).

3. Edit the lines in the server configuration file to suit your needs; they don't need much explanation.

   ![](./screenshot-02.png "World selection")

4. Press `Create server`. All necessary files will be generated inside the servers directory with the world's name. For example, if the world is called Deimos, it will be created inside `<app path>/servers/Deimos/`.

   ![](./screenshot-03.png "")

### Using the Server Container
Continuing with the above example, you can see the created files:
```bash
$ cd servers/Deimos && ls
start-server.sh
start-server.bat
stop-server.sh
stop-server.bat
data/
```
Files with the `.sh` extension are for Linux, and `.bat` files are for Windows.

1. Launch the container with `start-server.sh` and connect; there’s nothing more to explain. The first time, it will need to download and generate the image, so it will take a bit longer. This only happens once.
2. To safely stop the container, you can use the `/shutdown` command from the client or execute `stop-server.sh`.

## Caveats and Gotchas

**Error: Luanti data directory does not exist**. Edit `config.ini` and set the `data_path` key to the directory where Luanti stores its data.

***
Windows users: You may need to grant permissions to the directory to share it. Go to Docker Desktop > Settings > Resources > File Sharing.
***
World names will be sanitized: Converted to lowercase and special characters replaced with `_`. **This only affects internal container names.** However, similar names like `King Realm@` and `king realm!` will be converted to the same name, `king_realm_`, causing the latter to overwrite the former.

***
If the server does not launch or behaves unexpectedly, remove the `-d` option from `start-server.sh` to enter interactive mode and see the messages.
- `CTRL-C`: Close server.
- `CTRL-P` + `CTRL-Q`: Detach it to run in the background.
