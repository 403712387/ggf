#!/bin/bash
ulimit -n 10240
ulimit -c unlimited
ulimit -v unlimited
ulimit -m unlimited
GOMAXPROCS=2
mkdir logs >> /dev/null 2>&1
echo "-------------Welcome to ggf service------------`date "+%Y-%m-%d %H:%M:%S"`" >> ./logs/ggf.error
nohup ./ggf  >> ./logs/ggf.error 2>&1 &