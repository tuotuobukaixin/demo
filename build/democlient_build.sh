BUILD=$(cd $(dirname $0); pwd)
server_name=democlient
source $BUILD/common.sh

rm -rf $BUILD/../bin/$server_name.tgz
rm -rf $BUILD/$server_name
mkdir -p $BUILD/$server_name
mkdir -p $BUILD/$server_name/src/

cd $BUILD/$server_name


cp -r $BUILD/../src/* $BUILD/$server_name/src/
build_server $server_name $BUILD/$server_name $BUILD/$server_name/src/$server_name


if [ -f $BUILD/$server_name/src/$server_name/$server_name ]; then
        cp $BUILD/$server_name/src/$server_name/$server_name $BUILD/../bin/
        rm -rf $BUILD/$server_name
        echo "$server_name build successfully."
    else
        rm -rf $BUILD/$server_name
        echo "$server_name build failed."
fi