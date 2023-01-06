# kaltura-media-framework-docker-compose

## Running all the services

```
docker-compose stop && docker-compose down && docker-compose build && docker-compose up
```

## Simulating an live feed

```
./ffmpeg-multi-bitrate-example.sh
```

## Playing the stream
```
vlc http://localhost:9090/clear/ch/ch2/master.m3u8
# or you can open the url in your hls/dash player
```
