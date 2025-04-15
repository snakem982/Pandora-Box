package utils

import (
	"errors"
	"fmt"
	"github.com/snakem982/pandora-box/pandora/pkg/constant"
	"os"
	"os/user"
	"path/filepath"
)

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
		return os.Create(path)
	} else {
		// 其他错误
		return nil, err
	}
}

// GetUserHomeDir 获取当前用户的根目录
func GetUserHomeDir(paths ...string) (string, error) {
	// 尝试使用 os.UserHomeDir（Go 1.12+ 提供的函数）
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(append([]string{home, constant.DefaultWorkDir}, paths...)...), nil
	}

	// 如果 os.UserHomeDir 不适用，使用 os/user 包获取
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("获取当前用户失败: %v", err)
	}

	return filepath.Join(append([]string{currentUser.HomeDir, constant.DefaultWorkDir}, paths...)...), nil
}
