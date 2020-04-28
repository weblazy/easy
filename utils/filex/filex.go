package filex

import (
	"io/ioutil"
	"os"
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
