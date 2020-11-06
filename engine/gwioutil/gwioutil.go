package gwioutil

import (
	"github.com/xiaonanln/goworld/engine/gwlog"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type timeoutError interface {
	Timeout() bool // Is it a timeout error
}

// IsTimeoutError checks if the error is a timeout error
func IsTimeoutError(err error) bool {
	if err == nil {
		return false
	}

	err = errors.Cause(err)
	ne, ok := err.(timeoutError)
	return ok && ne.Timeout()
}

// WriteAll write all bytes of data to the writer
func WriteAll(conn io.Writer, data []byte) error {
	left := len(data)
	for left > 0 {
		n, err := conn.Write(data)
		if n == left && err == nil { // handle most common case first
			return nil
		}

		if n > 0 {
			data = data[n:]
			left -= n
		}

		if err != nil && !IsTimeoutError(err) {
			return err
		}
	}
	return nil
}

// ReadAll reads from the reader until all bytes in data is filled
func ReadAll(conn io.Reader, data []byte) error {
	left := len(data)
	for left > 0 {
		n, err := conn.Read(data)
		if n == left && err == nil { // handle most common case first
			return nil
		}

		if n > 0 {
			data = data[n:]
			left -= n
		}

		if err != nil && !IsTimeoutError(err) {
			return err
		}
	}
	return nil
}

// 查找指定的文件列表
type FileFilter struct {
	// file list in directory
	ListFile []string
	// 后缀 eg:.go .xlsx .txt ...
	Suffix string
}

func (this *FileFilter) Listfunc(path string, f os.FileInfo, err error) error {
	var strRet string
	/*strRet, _ = os.Getwd()

	if ostype == "windows" {
		strRet += "\\"
	} else if ostype == "linux" {
		strRet += "/"
	}*/

	if f == nil {
		return err
	}
	if f.IsDir() {
		return nil
	}

	strRet += path

	//用strings.HasSuffix(src, suffix)//判断src中是否包含 suffix结尾
	ok := strings.HasSuffix(strRet, this.Suffix)
	if ok {
		this.ListFile = append(this.ListFile, strRet) //将目录push到listfile []string中
	}

	return nil
}

func (this *FileFilter) GetFileList(path, suffix string) error {
	this.Suffix = suffix
	//var strRet string
	err := filepath.Walk(path, this.Listfunc)

	if err != nil {
		gwlog.FatalfE("filepath.Walk() returned %v\n", err)
		return err
	}

	return nil
}
