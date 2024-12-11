package tools

import (
	"github.com/klauspost/compress/zip"
	"github.com/metacubex/mihomo/log"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ZipDirectory 压缩指定目录及其所有子文件夹和文件
func ZipDirectory(sourceDir, outputZip string, exclude []string) error {
	zipFile, err := os.Create(outputZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {

		for _, suffix := range exclude {
			if strings.HasSuffix(path, suffix) {
				return nil
			}
		}

		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name, _ = filepath.Rel(sourceDir, path)
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				log.Errorln("os.Open path=%s, err=%v", path, err)
				return nil
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// 创建目标目录如果不存在
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		os.MkdirAll(dest, 0755)
	}

	// 解压每个文件
	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		//// 检查文件路径是否存在，避免文件泄露
		//if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
		//	return fmt.Errorf("illegal file path: %s", fpath)
		//}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), f.Mode()); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// 关闭资源
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

// CopyFile 复制文件
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

// CopyDirectory 拷贝目录及其所有子文件夹和文件
func CopyDirectory(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径并生成目标路径
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			log.Errorln("filepath.Rel path=%s,srcDir=%s, err=%v", path, srcDir, err)
			return nil
		}

		dstPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			// 如果是目录，创建目录
			err = os.MkdirAll(dstPath, info.Mode())
			if err != nil {
				log.Errorln("os.MkdirAll path=%s,dstPath=%s, err=%v", path, dstPath, err)
			}
			return nil
		} else {
			// 如果是文件，拷贝文件
			err = CopyFile(path, dstPath)
			if err != nil {
				log.Errorln("CopyFile path=%s,dstPath=%s, err=%v", path, dstPath, err)
			}

			return nil
		}
	})
}
