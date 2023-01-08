run:
	docker-compose stop && docker-compose down && docker-compose build && docker-compose up

origin:
	./ffmpeg-multi-bitrate-example.sh

lowlatency_origin:
	./ffmpeg-lowlatency-multi-bitrate-example.sh

test:
	docker-compose run --rm newcontrollertest

local_test:
	SEGMENTER_KMP_URL=segmenter_kmp_url go test ./...
