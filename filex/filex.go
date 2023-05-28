package filex

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Read(path string) ([]byte, error) {
	fi, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

func Write(path string, b []byte, isAppend bool) error {
	flag := os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	if isAppend {
		flag = os.O_WRONLY | os.O_APPEND | os.O_CREATE
	}
	fd, err := os.OpenFile(path, flag, os.ModePerm)
	if err != nil {
		return err
	}
	defer fd.Close()
	fd.Write(b)
	return nil
}

// GetPath 获取项目路径
func GetPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		print(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// CheckDir 判断文件目录否存在
func CheckDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

// MkdirDir 创建文件夹,支持x/a/a  多层级
func MkdirDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// RemoveDir 删除文件
func RemoveDir(filePath string) error {
	return os.RemoveAll(filePath)
}
