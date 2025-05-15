cd $(dirname $(readlink -f "$0"))
cd ../src-go
go build -tags=with_gvisor -trimpath -ldflags "-X github.com/snakem982/pandora-box/api.Version=v1.2.3" -o px