#! /bin/bash
BUILD=$(cd $(dirname $0); pwd)
ip=`ifconfig eth0 | grep "inet addr" | awk '{ print $2}' | awk -F: '{print $2}'`
sed -i "s|{{mysqluser}}|$mysqluser|g"  $BUILD/app.conf
sed -i "s|{{mysqlpwd}}|$mysqlpwd|g"  $BUILD/app.conf
sed -i "s|{{mysqlurl}}|$mysqlurl|g"  $BUILD/app.conf
sed -i "s|{{servername}}|$servername|g"  $BUILD/app.conf
sed -i "s|{{registerurl}}|$registerurl|g"  $BUILD/app.conf
sed -i "s|{{serveraddr}}|$serveraddr|g"  $BUILD/app.conf
sed -i "s|{{{podip}}|$ip|g"  $BUILD/app.conf
sed -i "s|{{database}}|$database|g"  $BUILD/app.conf
chmod 750 $BUILD/demotest
ps -ef|grep demotest
cd $BUILD
./demotest