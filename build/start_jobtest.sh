#! /bin/bash
BUILD=$(cd $(dirname $0); pwd)
wget  https://blog.csdn.net/pkueecser/article/details/50433460
chmod 750 $BUILD/jobtest
ps -ef|grep jobtest
cd $BUILD
nohup ./jobtest &
sleep $TIMEOUT
echo "test success $JOB_NAME" >> $BUILD/file/result

