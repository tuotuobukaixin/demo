#! /bin/bash
BUILD=$(cd $(dirname $0); pwd)
cd $BUILD
wget https://raw.githubusercontent.com/tuotuobukaixin/demo/master/src/jobtest/action/action.sh
chmod 750 action.sh
bash  action.sh



