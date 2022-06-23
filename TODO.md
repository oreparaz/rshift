## TODO

- [ ] Print small helper URL: "you can listen now at http://...."
- [ ] Add a minimal web UI. Controls: fast forward 10 min (skip a song), or ff 60 min (skip a program).
- [ ] Add support for multiple URLs, possibly read from a configuration file.
- [ ] Improve file seeking. Add an in-memory cache and store only a few files per directory (this will make it possible to run rshift in FAT filesystems that can only handle few files per directory.)
- [ ] Save disk space: do not cache during downtime hours (night).
- [ ] Self update https://github.com/rhysd/go-github-selfupdate
- [ ] Make the stream seekable with event playlist: https://developer.apple.com/documentation/http_live_streaming/example_playlists_for_http_live_streaming/event_playlist_construction
- [ ] TLS?
