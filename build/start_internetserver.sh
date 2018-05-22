#! /bin/bash
BUILD=$(cd $(dirname $0); pwd)
sed -i "s|{redis}|$redis|g"  $BUILD/app.conf
sed -i "s|{num}|$num|g"  $BUILD/app.conf
chmod 750 $BUILD/internetserver
ps -ef|grep internetserver
cd $BUILD
./internetserver