#!/usr/bin/env sh
set -eu

envsubst '${CONTROLLER_API_URL}' < /usr/local/openresty/nginx/conf/nginx.conf.template > /usr/local/openresty/nginx/conf/nginx.conf

/usr/local/openresty/bin/openresty -g  'daemon off;'
