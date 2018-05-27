#! /bin/bash
BUILD=$(cd $(dirname $0); pwd)
sed -i "s|{{mysqluser}}|$mysqluser|g"  $BUILD/app.conf
sed -i "s|{{mysqlpwd}}|$mysqlpwd|g"  $BUILD/app.conf
sed -i "s|{{mysqlurl}}|$mysqlurl|g"  $BUILD/app.conf
sed -i "s|{{timeout}}|$timeout|g"  $BUILD/app.conf
sed -i "s|{{jobname}}|$jobname|g"  $BUILD/app.conf
chmod 750 $BUILD/jobtest
cd $BUILD
./jobtest


