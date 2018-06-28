#! /bin/bash
chmod 750 redis-cli
chmod 750 jobtest
nohup ./jobtest &
sleep $timeout
./redis-cli -h $redis_url -p  $redis_port rpush joblist $jobname