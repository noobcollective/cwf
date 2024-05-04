# cwf - copy with friends

# Dependencies
- Added Zap as logging library go.uber.org/zap 
- Added yaml to parse our config file gopkg.in/yaml.v3 

## TODO:
- [ ] prefix paths with cwf home (preferable in a config file) - defaults to `/tmp/cwf/`
- [ ] more secure content (because base64 - wtf)

- [ ] chown directories to specific users
- [x] safe error handling (e.g. handle error responses in client)

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
- [ ] check if file exists
- [ ] hashing/enryption
- [x] list all files after date (more options for sorting?)
- [x] create dir if path passed as name

## Example commands - Executed by Client
# Sending Text to server
- [x] cat main.go | ./cwf testerfile

# Sending text range to server
- [ ] cat main.go | ./cwf --r 20-30 testerfile

# Pasting content from server to client
- [x] ./cwf testerfile

# Clear single file
- [x] ./cwf -d testerfile
- d -> delete

# Clear all files
- [ ] ./cwf -deleteAll

# List files in main dir and/or number of elements
- [x] ./cwf -l 
- [ ] ./cwf -lt
- l -> list
- lt -> show file tree

# List files in specific dir
- [ ] ./cwf -l dirName - needs client implementation

# Nice to haves
- [ ] ./cwf -size testerfile

## Later ideas:
- [ ] peer to peer
- [ ] diffs over snapshots
- [ ] only stdout a range (e.g cwf test1 -r 25,30)
