# kaltura-media-framework-docker-compose

```
# running all the services
docker-compose run --rm --service-ports all_in_one_nginx

# simulating an live feed
ffmpeg -re -i a.mp4 -c copy -f flv "rtmp://localhost:1935/live/ch1_s1"

# fetching the packaged live streaming formats
curl -v localhost/clear/ch/ch1/master.m3u8
```
