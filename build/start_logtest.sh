#! /bin/bash
BUILD=$(cd $(dirname $0); pwd)
cd $BUILD
for((i=1;i<=$TIMEOUT;i++));
do
	sleep 1
	time=`date +'%Y-%m-%d %H:%M:%S'`
	echo "$i test success $time"
	echo "$i test success $time" >> $BUILD/file/$JOB_NAME/result
done


