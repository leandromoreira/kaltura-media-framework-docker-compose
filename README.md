# kaltura-media-framework-docker-compose

```
# running all the services
docker-compose up

# in case you know what's happening individualy
# docker-compose logs -f controller
# docker-compose logs -f allinonenginx

# simulating an live feed
ffmpeg -re -i a.mp4 -c copy -f flv "rtmp://localhost:1935/live/ch1_s1"

# fetching live streaming from the packager
curl http://localhost:9090/clear/ch/ch1/master.m3u8
```
