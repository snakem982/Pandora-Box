cd $(dirname $(readlink -f "$0"))
cd ../build/bin
zip -r macos-amd64.zip ./pandora-box-amd64.app && rm -rf pandora-box-amd64.app
zip -r macos-arm64.zip ./pandora-box-arm64.app && rm -rf pandora-box-arm64.app
zip -r windows-amd64.zip pandora-box-amd64.exe && rm -rf pandora-box-amd64.exe
zip -r windows-arm64.zip pandora-box-arm64.exe && rm -rf pandora-box-arm64.exe