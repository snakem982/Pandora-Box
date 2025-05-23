package utils

import (
	"errors"
	"fmt"
	"github.com/snakem982/pandora-box/pkg/constant"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"sync"
)

var HomeDir string
var once sync.Once

func InitHomeDir(homeDir string) {
	once.Do(func() {
		HomeDir = homeDir
	})
}

// SetPermissions 设置权限允许所有用户读写
func SetPermissions(filePath string) error {
	if runtime.GOOS == "windows" {
		// Windows: 使用 icacls 赋予所有用户读写权限
		cmd := exec.Command("icacls", filePath, "/grant", "*S-1-1-0:(R,W)")
		return cmd.Run()
	} else {
		// macOS & Linux: 赋予 0666 权限（所有用户可读写）
		return os.Chmod(filePath, 0666)
	}
}

// FileExists 检查文件是否存在
func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false // 文件不存在
	}
	return !info.IsDir() // 如果路径存在且不是目录，则文件存在
}

// SaveFile 文件保存
func SaveFile(savePath string, content []byte) (bool, error) {
	// 检查路径是否合法
	if savePath == "" {
		return false, errors.New("保存路径不能为空")
	}

	// 检查文件是否存在
	if _, err := os.Stat(savePath); err == nil {
		// 如果文件存在，先删除
		err = os.Remove(savePath)
		if err != nil {
			return false, fmt.Errorf("删除文件失败: %v", err)
		}
	} else if !os.IsNotExist(err) {
		return false, fmt.Errorf("检查文件状态失败: %v", err)
	}

	// 创建保存路径的所有必要目录
	err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm)
	if err != nil {
		return false, fmt.Errorf("创建目录失败: %v", err)
	}

	// 创建并保存文件
	err = os.WriteFile(savePath, content, os.ModePerm)
	if err != nil {
		return false, fmt.Errorf("保存文件失败: %v", err)
	}
	_ = SetPermissions(savePath)

	return true, nil
}

// DeletePath 删除指定路径
func DeletePath(path string) error {
	// 检查路径是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("路径 %s 不存在", path)
	}

	// 使用 os.RemoveAll 删除路径
	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("删除路径失败: %v", err)
	}

	return nil
}

// CreateFile 根据路径创建文件，如果文件存在直接返回
func CreateFile(path string) (*os.File, error) {
	// 确保目录存在
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// 检查文件是否存在
	if _, err := os.Stat(path); err == nil {
		// 文件已存在，直接打开
		return os.OpenFile(path, os.O_RDWR, os.ModePerm)
	} else if os.IsNotExist(err) {
		// 文件不存在，创建新文件
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		_ = SetPermissions(path)
		return file, nil
	} else {
		// 其他错误
		return nil, err
	}
}

// CreateFileForAppend 以追加模式打开或创建文件，保证目录存在
func CreateFileForAppend(path string) (*os.File, error) {
	// 确保目录存在
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// 打开文件，使用追加模式
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	_ = SetPermissions(path)
	return file, nil
}

// ReadFile 根据传入的文件路径获取文件中的内容
func ReadFile(filePath string) (string, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("文件不存在: %s", filePath)
	}

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件时出错: %v", err)
	}

	// 返回文件内容作为字符串
	return string(content), nil
}

// GetUserHomeDir 获取当前用户的根目录
func GetUserHomeDir(paths ...string) string {
	if HomeDir != "" {
		return filepath.Join(append([]string{HomeDir}, paths...)...)
	}

	// 尝试使用 os.UserHomeDir（Go 1.12+ 提供的函数）
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(append([]string{home, constant.DefaultWorkDir}, paths...)...)
	}

	// 如果 os.UserHomeDir 不适用，使用 os/user 包获取
	currentUser, _ := user.Current()
	return filepath.Join(append([]string{currentUser.HomeDir, constant.DefaultWorkDir}, paths...)...)
}
