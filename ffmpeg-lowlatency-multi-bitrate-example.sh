ffmpeg -hide_banner \
-re -f lavfi -i "testsrc2=size=1280x720:rate=30,format=yuv420p" \
-f lavfi -i "sine=frequency=1000:sample_rate=4800" \
-c:v libx264 -preset ultrafast -tune zerolatency -profile:v high \
-b:v 1400k -bufsize 2800k -x264opts keyint=30:min-keyint=30:scenecut=-1 \
-c:a aac -b:a 128k -f flv rtmp://localhost:1935/live/ll_ch3_720 \
-c:v libx264 -preset ultrafast -tune zerolatency -profile:v high \
-b:v 750k -bufsize 1500k -s 640x360 -x264opts keyint=30:min-keyint=30:scenecut=-1 \
-c:a aac -b:a 128k -f flv rtmp://localhost:1935/live/ll_ch3_360
