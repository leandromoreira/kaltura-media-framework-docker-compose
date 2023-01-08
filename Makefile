run:
	docker-compose stop && docker-compose down && docker-compose build && docker-compose up

origin:
	./ffmpeg-multi-bitrate-example.sh

lowlatency_origin:
	./ffmpeg-lowlatency-multi-bitrate-example.sh

test:
	docker-compose run --rm newcontrollertest

local_test:
	CC_DECODER_URL=http://cc_decoder_url \
		       SEGMENTER_API_URL=http://segmenter_api_url \
		       SEGMENTER_KMP_URL=http://segmenter_kmp_url \
		       go test ./...

