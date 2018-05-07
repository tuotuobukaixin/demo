BUILD=$(cd $(dirname $0); pwd)




function build_server(){
    export GOPATH=$2:$BUILD/../third_party
    cd $3
    go build -gcflags=-trimpath=$2 -o $1
}