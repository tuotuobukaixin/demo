BUILD=$(cd $(dirname $0); pwd)
server_name=logtest
source $BUILD/common.sh

rm -rf $BUILD/$server_name
mkdir -p $BUILD/$server_name
cd $BUILD/$server_name
cp $BUILD/${server_name}_dockerfile Dockerfile
cp $BUILD/start_${server_name}.sh .
docker build -t ${server_name} .
docker save -o $BUILD/../bin/${server_name}.tar ${server_name}:latest
