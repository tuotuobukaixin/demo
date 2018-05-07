#! /bin/bash
BUILD=$(cd $(dirname $0); pwd)
sed -i "s|{{mysqluser}}|$mysqluser|g"  $BUILD/app.conf
sed -i "s|{{mysqlpwd}}|$mysqlpwd|g"  $BUILD/app.conf
sed -i "s|{{mysqlurl}}|$mysqlurl|g"  $BUILD/app.conf
sed -i "s|{{servername}}|$servername|g"  $BUILD/app.conf
sed -i "s|{{registerurl}}|$registerurl|g"  $BUILD/app.conf
sed -i "s|{{serveraddr}}|$serveraddr|g"  $BUILD/app.conf
chmod 750 $BUILD/gameserver
ps -ef|grep gameserver
cd $BUILD
./gameserver