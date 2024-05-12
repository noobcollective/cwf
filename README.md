# cwf - copy with friends
`cwf` is a small tool to quickly share content (e.g. logfiles, commands, ...) with your friends and colleagues. It can be run local or on a server, so nobody needs to set things up like port forwarding on their local machine.

:warning: Warning:\
As of now we do **not recommend** sending sensitive data! Although the transport can be handled via HTTPS, the files on the server are just `base64` encoded and can be accessed via the given endpoints from any IP.

## Demo
![cwf_demo](https://github.com/noobcollective/cwf/assets/99724871/d2a2c286-516a-4db5-b3b4-d8ff34e8867c)


## Installation
We provide a homebrew package to install `cwf` as a binary. \
It's totally possible to run the server with the binary. \
However if you are planning to run `cwf` on a server (e.g. VPS) we recommend using our [docker image](#docker-image).

### MacOS & Linux
```bash
brew install noobcollective/tap/cwf
```

## Example config
```yaml
---
 motherShipIP: 127.0.0.1
 motherShipPort: 8787
 motherShipSSL: true  <- false to disable HTTPS
```

## Client Usage
### Send content to server and save it in a single file.
```
echo "Hello Clipboard!" | cwf testfile
```
`Saved to: testfile`

### Get content of a saved file from server.
```
cwf testfile
```
`Hello Clipboard!`

### Send content to server and save file in a subdirectory.
```
echo "Hello Clipboard from subdirectory!" | cwf testdir/testfile
```
`Saved to: testdir/testfile`

### Get content of a saved file in a subdirectory from server.
```
cwf testdir/testerfile
```
`Hello Clipboard from subdirectory!`

### Delete single file
```
cwf -d testerfile
```
`Deleted file: testerfile`

### Delete single file in subdirectory
```
cwf -d testdir/testfile
```
`Deleted file: testfile`

### Delete all files in subdirectory
```
cwf -d testdir/
```
`Deleted directory: testdir`

### List files and directories in configured base directory.
```
cwf -l
```
```
Type    Name         Modified
Dir    testdir      2006-01-02 15:04:05
File   testfile.cwf 2006-01-02 15:04:05
```

### List files and directories in a subdirectory.
```
cwf -l testdir
```
```
Type    Name         Modified
File   testfile.cwf 2006-01-02 15:04:05
```

## Docker Image
- Get the image with:
    ```bash
    docker pull noobcollective/cwf
    ```
- Start `cwf` with docker in HTTP Mode:
    ```bash
    docker run \
        --name cwf-server \
        -p <host_port>:<container_port> \
        -v <volume_name>:<container_filesDir> \
        -d noobcollective/cwf -serve
    ```
- Start `cwf` with docker in HTTPS Mode:
    ```bash
    docker run \
        --name cwf-server \
        -p <host_port>:<container_port> \
        -v <volume_name>:<container_filesDir> \
        -v <container_ssl_certs_dir>:<host_ssl_certs_dir> \
        -d noobcollective/cwf -serve -https -certfile <path_to_ssl_cert> -keyfile <path_to_ssl_key>
    ```
- `<container_port>`: Configurable via `-port` argument in cwf -> defaults to 8787
- `<volume_name>`: Name of the docker volume for persisting data.
- `<container_filesDir>`: Configurable via `-filesDir` argument in cwf -> defaults to `/tmp/cwf/`

When using the `-https` flag, you need to mount the directory where the SSL certificates are located on your host machine \
and provide the names of your ssl certificate and key to `cwf` with the `-certfile` and `-keyfile` arguments.
- `<host_ssl_certs_dir>`: Directory where the SSL certificates are located on your host machine.
- `<container_ssl_certs_dir>`: Configurable via `-certsdir` argument in cwf -> defaults to `/etc/crypts`

## Roadmap:
- [ ] get size of files
- [ ] diffs over snapshots
- [ ] chown directories to specific users
- [ ] only stdout a range (e.g cwf test1 -r 25,30)


## Dependencies
- Added Zap as logging library go.uber.org/zap
- Added yaml to parse our config file gopkg.in/yaml.v3


## TODOs:
- [x] prefix paths with cwf home (preferable in a config file) - defaults to `/tmp/cwf/`
- [ ] more secure content (because base64 - wtf)
- [x] safe error handling (e.g. handle error responses in client)
- [x] fix path for config file in resulting binary (for now it uses `pwd`, which is not a good idea, because `cwf` should lay in `<somewhere>/bin`)


## Feature list - Server
- [x] copy into cwf
- [x] stdout of cwf
- [x] clean file
- [x] check if file exists
- [x] create dir if path passed as name
- [x] list all files after date (more options for sorting?)
  - Currently only listing files in directory in modified order. NO additional sorting is supported
- [x] adding flag to setup port additionaly reading from config file setting dir/dirDepth/port
  - Implemented with yaml (I thought yaml is toml and toml is yaml) so currently its working with a yaml file but i actually want a toml file


## Feature list - Client
- [x] copy into cwf
- [x] stdout of cwf
- [x] clean file
~~- [ ] check if file exists~~ - not needed (`cwf <filename>` returns if file exists anyway)
- [ ] hashing/enryption
- [x] list all files after date (more options for sorting?)
- [x] create dir if path passed as name


## Our Blog
https://project-folio.eu
