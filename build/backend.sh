cd $(dirname $(readlink -f "$0"))
cd ../src-go
go build -tags=with_gvisor -trimpath -o px