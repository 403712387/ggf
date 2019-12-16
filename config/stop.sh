#!/bin/bash
echo "-------------Bye-----------------------`date "+%Y-%m-%d %H:%M:%S"`" >> ./logs/host.error
nohup curl -H 'User-Agent:Panda' http://127.0.0.1:6003/ggf/stop  > /dev/null 2>&1 &