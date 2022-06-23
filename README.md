# Time-shifted HLS streaming


https://github.com/oreparaz/rshift


Those living in California know that listening to the BBC Radio 4 live stream isn't the same thing ---
when you're waking up, people in London are already eating supper.
Fine disc jockeys won't play same music at 8am or 8pm; time zones messes this up.

This simple project allows you to listen to internet radio streams (HLS) with a configurable time delay.
You end up listening to your favorite show at the proper
(local) time: you listen at 8pm California time what was broadcasted at 8pm in London.

It works by serving a time shifted M3U file and caching the fragmented stream.
Time shifting into the future isn't supported yet, unfortunately.

## Deploy

You need a server 24/7 to run this. You don't need a very powerful machine. A home server
with a large disk will work, or any small cloud VM instance will work fine too.

Example usage:
```
curl -LO https://github.com/oreparaz/rshift/releases/latest/download/rshift
./rshift -download-m3u8-url https://radio/playlist.m3u8 -output-path=/path/to 
```

Then you can visit http://host:8080/timeshift/1234.m3u8 if you want to listen to the stream with a time delay of 1234 seconds.
(You'll probably need NAT to access from clients outside of your LAN).

## FAQ

__Will this eventually eat all my disk space?__
`rshift` keeps only 5 days of stream data and deletes any files older than that. So no, it won't end up filling all your disk.

__Which client should I use?__
Works great with Safari on iOS/MacOS. On iOS the stream plays even when Safari is minimized or the screen is locked.

__Which cloud provider should I use?__
If you don't want to keep your own server at home, I found Google Compute Platform is a cool environment to run rshift.
You can run this on a micro instance that qualifies for the free forever tier.
You pay for _egress_ traffic (_ingress_ is for free), so you just pay for what you listen.

__No authentication?__

You can set the following environment variables to turn basic HTTP auth on:
```
RSHIFT_USERNAME=foo RSHIFT_PASSWORD=bar ./rshift ...
```

## Status

This is "code complete" for my (very) specific use case.
`rshift` has been running solidly for more than one year.
Feature-wise it's pretty barebones.
There is a tons of items in the TODO if you want to tackle on those.

