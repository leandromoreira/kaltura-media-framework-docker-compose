# kaltura-media-framework-docker-compose

A docker image (and a docker-compose) for the [Kaltura Live Media Framework](https://github.com/kaltura/media-framework) passthrough just for fun/learning.
<img width="1440" alt="image" src="https://user-images.githubusercontent.com/55913/211084171-52b607bd-4030-40e6-a41d-be9743ea926c.png">

## Running all the services

```
docker-compose stop && docker-compose down && docker-compose build && docker-compose up
```

## Simulating a multi-bitrate live stream feed

```
./ffmpeg-multi-bitrate-example.sh
# to test the low latency stream use ./ffmpeg-lowlatency-multi-bitrate-example.sh
```

## Playing the live stream

```
vlc http://localhost:9090/clear/ch/ch2/master.m3u8
# or you can open the url in your hls/dash player

# HLS AES-128 http://localhost:9090/aes128/ch/ch2/master.m3u8
# HLS SAMPLE-AES http://localhost:9090/cbcs/ch/ch2/master.m3u8
# DASH http://localhost:9090/clear/ch/ch2/manifest.mpd

# to test the low latency stream use vlc http://localhost:9090/clear/ch/ch3/master.m3u8
```

## TODO

* port the `transcoding` part https://github.com/kaltura/media-framework/tree/master/conf to `main.go`.
  * explore other options such as DRM.
* refactor the golang code. :bomb:
  * use struct instead of `map[string]interface{}`, once the API format is known one can use the [json to struct](https://json2struct.mervine.net/) to aid this task.
* template nginx.conf to set hosts by enviroment variable
* try to do the `ffmpeg` origin simulator over docker instead of locally.
* add more workflow examples such as: `srt input`, `mpegts`, and, etc.
* add the infamous big buck bunny file in loop? `ffmpeg -stream_loop -1 -re -i bbb.mp4`.
* maybe create a `docker-compose` in which each server has its own process.
