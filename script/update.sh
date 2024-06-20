cd $(dirname $(readlink -f "$0"))
cd ../
export http_proxy=http://127.0.0.1:10000
export https_proxy=http://127.0.0.1:10000
go get -u