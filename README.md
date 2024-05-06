# cwf - copy with friends

## Installation

### MacOS
```bash
$ brew install iculturebud/tap/cwf
```


## Example config
```yaml
---
 motherShipIP: "127.0.0.1"
 motherShipPort: "8787"
```


## Dependencies
- Added Zap as logging library go.uber.org/zap
- Added yaml to parse our config file gopkg.in/yaml.v3

## TODO:
- [ ] prefix paths with cwf home (preferable in a config file) - defaults to `/tmp/cwf/`
- [ ] more secure content (because base64 - wtf)
- [x] safe error handling (e.g. handle error responses in client)
- [ ] fix path for config file in resulting binary (for now it uses `pwd`, which is not a good idea, because `cwf` should lay in `<somewhere>/bin`)

## Feature list - Server
- [x] copy into cwf
- [x] stdout of cwf
- [x] clean file
- [x] check if file exists
- [x] create dir if path passed as name
- [x] list all files after date (more options for sorting?)
  - Currently only listing files in directory in modified order. NO additional sorting is supported
- [ ] adding flag to setup port additionaly reading from config file setting dir/dirDepth/port
  - Implemented with yaml (I thought yaml is toml and toml is yaml) so currently its working with a yaml file but i actually want a toml file

## Feature list - Client
- [x] copy into cwf
- [x] stdout of cwf
- [x] clean file
~~- [ ] check if file exists~~ - not needed (`cwf <filename>` returns if file exists anyway)
- [ ] hashing/enryption
- [x] list all files after date (more options for sorting?)
- [x] create dir if path passed as name

## Example commands - Executed by Client
# Sending Text to server
- [x] cat main.go | ./cwf testerfile

# Sending text range to server
~~- [ ] cat main.go | ./cwf --r 20-30 testerfile~~ -- some could specify range with other cli tools

# Pasting content from server to client
- [x] ./cwf testerfile

# Clear single file
- [x] ./cwf -d testerfile
- d -> delete

# Clear all files
- [ ] ./cwf -deleteAll

# List files in main dir and/or number of elements
- [x] ./cwf -l (list files in main cwf dir or specified subdir)
- [ ] ./cwf -lt (list files in nice file tree formatting)

# List files in specific dir
- [x] ./cwf -l dirName

# Nice to haves
- [ ] ./cwf -size testerfile

## Later ideas:
- [ ] peer to peer
- [ ] diffs over snapshots
- [ ] chown directories to specific users
- [ ] only stdout a range (e.g cwf test1 -r 25,30)

## Our Blog
https://project-folio.eu
