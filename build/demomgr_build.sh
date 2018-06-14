BUILD=$(cd $(dirname $0); pwd)
server_name=demomgr
source $BUILD/common.sh

rm -rf $BUILD/../bin/$server_name.tgz
rm -rf $BUILD/$server_name
mkdir -p $BUILD/$server_name
mkdir -p $BUILD/$server_name/src/

cd $BUILD/$server_name


cp -r $BUILD/../src/* $BUILD/$server_name/src/
build_server $server_name $BUILD/$server_name $BUILD/$server_name/src/$server_name


if [ -f $BUILD/$server_name/src/$server_name/$server_name ]; then
        cd $BUILD/$server_name
        cp $BUILD/${server_name}_dockerfile Dockerfile
        cp $BUILD/start_${server_name}.sh .
        docker build -t swr.cn-north-1.myhuaweicloud.com/cce-demo/${server_name}:latest .
        cp -r $BUILD/$server_name/src/$server_name/conf .
        cp $BUILD/$server_name/src/$server_name/$server_name .
        tar -zcvf $BUILD/../bin/$server_name.tgz conf $server_name
        rm -rf $BUILD/$server_name
        echo "$server_name build successfully."
    else
        rm -rf $BUILD/$server_name
        echo "$server_name build failed."
fi