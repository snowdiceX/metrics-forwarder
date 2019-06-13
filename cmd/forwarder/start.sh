#!/bin/sh

# nohup /apps/prometheus/pushgateway &

nohup /apps/prometheus/metrics-forwarder start \
    --pull http://127.0.0.1:26660/metrics \
    --push http://127.0.0.1:9091 \
    >> /apps/logs/forwarder/start.log 2>&1 &
