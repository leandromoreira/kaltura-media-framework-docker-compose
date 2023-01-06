# kaltura-media-framework-docker-compose

```
# running all the services
docker-compose stop && docker-compose down && docker-compose build && docker-compose up

# in case you know what's happening individualy
# docker-compose logs -f controller
# docker-compose logs -f allinonenginx

# simulating an live feed
# download a.mp4 video if you don't have it already
# wget http://cdnapi.kaltura.com/p/2035982/playManifest/entryId/0_w4l3m87h/flavorId/0_vsu1xutk/format/download/a.mp4
ffmpeg -stream_loop -1 -re -i a.mp4 -c copy -f flv "rtmp://localhost:1935/live/ch1_s1"

# fetching live streaming from the packager or open it in any hls capable player
curl http://localhost:9090/clear/ch/ch1/master.m3u8
```
