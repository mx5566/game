package excelt

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/xiaonanln/goworld/engine/gwlog"
	"os"
	"path/filepath"
	"strings"
)

var (
	ostype = os.Getenv("GOOS") // 获取系统类型
)
var listfile []string //获取文件列表

func Read(fileName string) {

	f, err := excelize.OpenFile(fileName)
	if err != nil {
		println(err.Error())
		return
	}
	// 获取工作表中指定单元格的值
	cell := f.GetCellValue("Sheet1", "B2")
	if cell == "" {
		gwlog.Error("read xlsx error" + fileName)
	}

	println(cell)
	// 获取 Sheet1 上所有单元格
	rows := f.GetRows("Sheet1")
	for _, row := range rows {
		for _, colCell := range row {
			print(colCell, "\t")
		}
		println()
	}
}

func ReadAllXlsx(path string) {
	err := GetFileList(path)
	if err != nil {
		return
	}

}

func Listfunc(path string, f os.FileInfo, err error) error {
	var strRet string
	strRet, _ = os.Getwd()

	if ostype == "windows" {
		strRet += "\\"
	} else if ostype == "linux" {
		strRet += "/"
	}

	if f == nil {
		return err
	}
	if f.IsDir() {
		return nil
	}

	strRet += path

	//用strings.HasSuffix(src, suffix)//判断src中是否包含 suffix结尾
	ok := strings.HasSuffix(strRet, ".xlsx")
	if ok {

		listfile = append(listfile, strRet) //将目录push到listfile []string中
	}

	return nil
}

func GetFileList(path string) error {
	//var strRet string
	err := filepath.Walk(path, Listfunc)

	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
		gwlog.Fatalf()
		return err
	}

	return nil
}

func ListFileFunc(p []string) {
	for index, value := range p {
		fmt.Println("Index = ", index, "Value = ", value)
	}
}
