#!/bin/sh

# nohup /apps/prometheus/pushgateway &

# /apps/prometheus/metrics-forwarder version
# /apps/prometheus/metrics-forwarder -h

nohup /apps/prometheus/metrics-forwarder start \
    --log ./log.conf \
    --config ./config.conf \
    >> /apps/logs/forwarder/start.log 2>&1 &
