#!/bin/bash
set -e
PWD=$(cd $(dirname $0); pwd)
BASE_DIR=${PWD}/..
BUILD_LOG_DIR=${PWD}/logs
mkdir -p $BUILD_LOG_DIR

date -R > $BUILD_LOG_DIR/time.log

ALL_SUBSYSTEM=(paas_controller k8s_adaptor k8s_runtime healthcheck runtime_supervisor backup_recovery_agent backup_recovery_server access_mgr ops bluprint_generator cloudify_runtime cde_api_server orch composer)
for system in ${ALL_SUBSYSTEM[@]}
do
{
    echo "=================begin to build ${system} at `date -R`=================" >> $BUILD_LOG_DIR/time.log
    bash ${PWD}/${system}_build.sh $BUILDTYPE > $BUILD_LOG_DIR/${system}.log 2>&1
    echo "=================finish build ${system} at `date -R`=================" >> $BUILD_LOG_DIR/time.log
}&
done
wait


begin_number=`cat $BUILD_LOG_DIR/time.log | grep begin | wc -l`
finish_number=`cat $BUILD_LOG_DIR/time.log | grep  finish| wc -l`

if [ $begin_number != $finish_number ];then
  echo Miss packages, build is fail
  exit 1
fi

date -R >> $BUILD_LOG_DIR/time.log