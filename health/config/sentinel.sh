#!/bin/bash
cp /redis-config/sentinel.conf /sentinel.conf
while ! ping -c 1 redis-premium-master-0.redissvc.default.svc.cluster.local; do
    echo 'Waiting for server'
    sleep 1
done

redis-sentinel /sentinel.conf
