# Luanti Server Creator

Para crear un servidor de Luanti en Docker con 2 clicks de ratón.

## Cómo usar
En el equipo que hará de servidor, tendréis que tener instalado Docker.

Usuarios de Windows: Para ese Sistema Operativo, la aplicación está en fase **totalmente experimental**, yo no uso Windows y aunque he procurado programarlo multiplataforma y lo he compilado para Windows, no he podido probarlo personalmente, así que no puedo asegurar que funcione.


### Descargar
Si no quieres compilarlo por ti mismo, descarga los ejecutables que están en el directorio `bin/` del repositorio. `servercreator` para Linux y `servercreator.exe` para Windows.

### Desde el cliente del juego
1. (Opcional) Crea un nuevo mundo del tipo de juego que se quiera (minetest, mineclone, voxelibre...). `Importante`, aunque sea opcional este paso, es recomendable lanzar el servidor con un mundo nuevo.
2. (Opcional) Activarle los mods que se quieran usar.

### Lanzamos la aplicación
Es recomendable hacerlo desde la linea de comandos para ver los posibles mensajes de error.

*NOTA*: La primera vez que se lanza se creará un fichero `config.ini` donde puedes indicar en qué directorio se crearán los servidores. Por defecto se crean dentro `<app path>/servers`. Si quieres cambiarlo, ciérrala, modifica el valor y ábrela de nuevo.

1. En el desplegable, seleccionamos nuestro mundo.
 
   ![](./screenshot-01.png "Word selection")

2. Puedes seleccionar la versión del servidor a usar o dejar la última disponible (latest).

3. Edita las lineas del fichero de configuración del servidor adaptadas a tus necesidades, no necesitan mucha explicación.

   ![](./screenshot-02.png "Word selection")

4. Pulsamos `Create server`. Se generarán todos los ficheros necesarios dentro del directorio servers con el nombre del mundo. Por ejemplo, si el mundo se llama Deimos, se creará dentro de `<app path>/servers/Deimos/`.

   ![](./screenshot-03.png "")

### Usar el contenedor del servidor
Si seguimos con el ejemplo anterior vemos los ficheros creados
```bash
$ cd servers/Deimos && ls
start-server.sh
start-server.bat
stop-server.sh
stop-server.bat
data/
```
Los archivos con extensión `.sh` son de Linux y los `.bat` de Windows.

1. Lanza el contenedor con `start-server.sh` y conéctate, no hay más que explicar. La primera vez se tendrá que descargar y generar la imagen y tardará un poco más. Esto sólo sucedera una vez.
2. Si quieres parar el contenedor de forma segura, puedes hacer un comando `/shutdown` desde el cliente, o bien ejecutar `stop-server.sh`.

## Caveats and gotchas

**Error Luanti data directory does not exists**. Edita `config.ini` y pon el valor de la clave `data_path` en el directorio donde Luanti guarda los datos.

***
Usuarios de Windows: Es posible que tengáis que darle permisos al directorio para compartirlo. Docker Desktop > Settings > Resources > File Sharing.
***
Los nombres de los mundos serán `sanitizados`: Convertidos a minúsculas y reemplazando caracteres especiales a `_`. **Esto sólo afecta internamente al nombre de los contenedores**. Sin embargo, nombre parecidos como `King Realm@` y `king realm!` serán convertidos al mismo nombre `king_realm_`, lo que hará que el último acabe sobrescribiendo el anterior.

***
Si el servidor no se lanza o se comporta como quieres, elimina la opción `-d` de `start-server.sh` para entrar en el modo interactivo y ver los mensajes.
-`CTRL-C` cerrar servidor 
-`CTRL-P` + `CTRL-Q` dejarlo en segundo plano.
