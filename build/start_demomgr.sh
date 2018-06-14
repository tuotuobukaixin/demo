#! /bin/bash
BUILD=$(cd $(dirname $0); pwd)
sed -i "s|{{mysqluser}}|$mysqluser|g"  $BUILD/app.conf
sed -i "s|{{mysqlpwd}}|$mysqlpwd|g"  $BUILD/app.conf
sed -i "s|{{mysqlurl}}|$mysqlurl|g"  $BUILD/app.conf
sed -i "s|{{database}}|$database|g"  $BUILD/app.conf

chmod 750 $BUILD/demomgr
ps -ef|grep demomgr
cd $BUILD
if [ "$debug" == "true" ]; then
    sleep 3600
fi
./demomgr