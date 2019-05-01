#!/bin/bash
if [[ ${HOSTNAME} == 'redis-premium-master-0' ]]; then
  redis-server /redis-config/master.conf
else
  redis-server /redis-config/slave.conf
fi
