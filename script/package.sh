cd $(dirname $(readlink -f "$0"))
cd ../
if [ "$1" = "windows" ]; then
  wails build -tags with_gvisor -skipbindings -m -s -trimpath -nosyncgomod -platform windows
elif [ "$1" = "darwin" ]; then
  wails build -tags with_gvisor -skipbindings -m -s -trimpath -nosyncgomod -platform darwin
elif [ "$1" = "linux" ]; then
  wails build -tags with_gvisor -skipbindings -m -s -trimpath -nosyncgomod -platform linux
elif [ "$1" = "all" ]; then
  wails build -tags with_gvisor -skipbindings -m -s -trimpath -nosyncgomod -platform "darwin/amd64,darwin/arm64,windows/amd64,windows/arm64"
else
  wails build -tags with_gvisor -skipbindings -m -s -trimpath -nosyncgomod
fi
