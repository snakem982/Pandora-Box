cd $(dirname $(readlink -f "$0"))
cd ../
wails dev -tags with_gvisor -skipbindings -m -s -nosyncgomod -appargs "-dev"