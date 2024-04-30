# cwf - copy with friends

## TODO:
- [ ] prefix paths with cwf home (preferable in a config file) - defaults to `/tmp/cwf/`
- [ ] more secure content (because base64 - wtf)
- [ ] chown directories to specific users

## Feature list - Server
- [x] copy into cwf
- [x] stdout of cwf
- [x] clean file
- [x] check if file exists
- [x] create dir if path passed as name
- [ ] hashing/enryption
- [ ] list all files after date (more options for sorting?)
- [ ] adding flag to setup port additionaly reading from config file setting dir/dirDepth/port

## Feature list - Client
- [ ] copy into cwf
- [ ] stdout of cwf
- [ ] clean file
- [ ] check if file exists
- [ ] hashing/enryption
- [ ] list all files after date (more options for sorting?)
- [ ] create dir if path passed as name

## Later ideas:
- [ ] peer to peer
- [ ] diffs over snapshots
- [ ] only stdout a range (e.g cwf test1 -r 25,30)
