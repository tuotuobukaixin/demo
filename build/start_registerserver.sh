#! /bin/bash
BUILD=$(cd $(dirname $0); pwd)
sed -i "s|{{mysqluser}}|$mysqluser|g"  $BUILD/app.conf
sed -i "s|{{mysqlpwd}}|$mysqlpwd|g"  $BUILD/app.conf
sed -i "s|{{mysqlurl}}|$mysqlurl|g"  $BUILD/app.conf

chmod 750 $BUILD/registerserver
ps -ef|grep registerserver
cd $BUILD
./registerserver