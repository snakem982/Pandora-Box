//go:build darwin

package update

import (
	"archive/zip"
	"bytes"
	"fmt"
	C "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/log"
	"io"
	"os"
	"pandora-box/backend/constant"
	"pandora-box/backend/system/proxy"
	"pandora-box/backend/tools"
	"path/filepath"
	"runtime"
)

func IsNeedUpdate() (bool, string) {
	all, _ := tools.ConcurrentHttpGet(constant.PandoraVersionUrl)
	if all != nil && len(all) > 0 {
		ver := string(all)
		return ver != constant.PandoraVersion, ver
	}

	return false, constant.PandoraVersion
}

func Replace(version string) bool {
	// 获取架构
	arch := runtime.GOARCH
	// 获取下载地址
	url := fmt.Sprintf(constant.PandoraDownloadUrl, version, "macos", arch)
	// 获取压缩文件
	all, _ := tools.ConcurrentHttpGet(url)
	if all != nil && len(all) > 0 {
		log.Infoln("new version file download success")
		// 解压
		err := UnzipFromBytes(C.Path.HomeDir(), all)
		if err != nil {
			log.Errorln("new version file unzip error:", err)
			return false
		}
		log.Infoln("new version file unzip success")
		exePath, err := os.Executable()
		if err != nil {
			log.Errorln("get exe path error：", err)
			return false
		}
		// 进行文件替换
		appFileName := fmt.Sprintf("pandora-box-%s.app", arch)
		src := filepath.Join(C.Path.HomeDir(), appFileName, "/Contents/MacOS/Pandora-Box")
		_, err = proxy.Command("mv", src, exePath)
		if err != nil {
			log.Errorln("new version file replace error:", err)
			return false
		}
		log.Infoln("new version file replace success")
		// 删除版本文件
		_, err = proxy.Command("rm", "-rf", filepath.Join(C.Path.HomeDir(), appFileName))

		return true
	}

	return false
}

// UnzipFromBytes 解压压缩字节流
// @params dst string 解压后目标路径
// @params src []byte 压缩字节流
func UnzipFromBytes(dst string, src []byte) error {
	// 通过字节流创建zip的Reader对象
	zr, err := zip.NewReader(bytes.NewReader(src), int64(len(src)))
	if err != nil {
		return err
	}

	// 解压
	return Unzip(dst, zr)
}

// Unzip 解压压缩文件
// @params dst string      解压后的目标路径
// @params src *zip.Reader 压缩文件可读流
func Unzip(dst string, src *zip.Reader) error {
	// 强制转换一遍目录
	dst = filepath.Clean(dst)

	// 遍历压缩文件
	for _, file := range src.File {
		// 在闭包中完成以下操作可以及时释放文件句柄
		err := func() error {
			// 跳过文件夹
			if file.Mode().IsDir() {
				return nil
			}

			// 配置输出目标路径
			filename := filepath.Join(dst, file.Name)
			// 创建目标路径所在文件夹
			e := os.MkdirAll(filepath.Dir(filename), 0750)
			if e != nil {
				return e
			}

			// 打开这个压缩文件
			zfr, e := file.Open()
			if e != nil {
				return e
			}
			defer func(zfr io.ReadCloser) {
				err := zfr.Close()
				if err != nil {

				}
			}(zfr)

			// 创建目标文件
			fw, e := os.Create(filename)
			if e != nil {
				return e
			}
			defer func(fw *os.File) {
				err := fw.Close()
				if err != nil {

				}
			}(fw)

			// 执行拷贝
			_, e = io.Copy(fw, zfr)
			if e != nil {
				return e
			}

			// 拷贝成功
			return nil
		}()

		// 是否发生异常
		if err != nil {
			return err
		}
	}

	// 解压完成
	return nil
}
