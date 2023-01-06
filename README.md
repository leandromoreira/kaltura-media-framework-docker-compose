# kaltura-media-framework-docker-compose

A docker image (and a docker-compose) for the passthrough for [Kaltura Live Media Framework](https://github.com/kaltura/media-framework) just for fun/learning.
<img width="1440" alt="image" src="https://user-images.githubusercontent.com/55913/211084171-52b607bd-4030-40e6-a41d-be9743ea926c.png">

## Running all the services

```
docker-compose stop && docker-compose down && docker-compose build && docker-compose up
```

## Simulating a multi-bitrate live stream feed

```
./ffmpeg-multi-bitrate-example.sh
```

## Playing the live stream
```
vlc http://localhost:9090/clear/ch/ch2/master.m3u8
# or you can open the url in your hls/dash player
```
