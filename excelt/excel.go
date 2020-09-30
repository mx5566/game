package excelt

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/xiaonanln/goworld/engine/gwlog"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var (
	ostype = runtime.GOOS // 获取系统类型
)
var listfile []string //获取文件列表
var MapItem map[string]*ItemBase

type TableRes interface {
	Get(key string)
	GetInt(key int64)
}

// 加载table
func init() {
	Load()
}

func Load() {
	var listpath = "."
	_ = GetFileList(listpath)
	pathHead := "excelt"

	if ostype == "windows" {
		pathHead += "\\"
	} else if ostype == "linux" {
		pathHead += "/"
	}
	for _, path := range listfile {
		switch path {
		case pathHead + "equip.xlsx":
			LoadEquip(path)
		case pathHead + "item.xlsx":
			LoadItem(path)
		default:
			fmt.Println("error path " + path)
		}
	}
}

func LoadItem(path string) {
	items := Read(path, "ID")
	fmt.Println("load table item !!!")
	fmt.Println(items)
}

func LoadEquip(path string) {
	equip := Read(path, "ID")
	fmt.Println("load table equip !!!")
	fmt.Println(equip)

}

func compressStr(str string) string {
	if str == "" {
		return ""
	}
	//匹配一个或多个空白符的正则表达式
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(str, "")
}

func ReadAllXlsx(path string) {
	err := GetFileList(path)
	if err != nil {
		return
	}
}

func CombineKeys(keys ...interface{}) string {
	//sort.Strings(keys)
	com := []string{}
	for _, key := range keys {
		switch key.(type) {
		case int, int32, int64, int8, int16, uint, uint32, uint64, uint16, uint8:
			com = append(com, reflect.ValueOf(key).String())
		case string:
			com = append(com, key.(string))
		default:
			fmt.Println("unkonw type " + reflect.TypeOf(key).String())
		}
	}
	return strings.Join(com, "_")
}

func Read(fileName string, keys ...string) map[interface{}]map[string]interface{} {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		println(err.Error())
		return nil
	}

	// 找出key
	mapKeys := make(map[string]interface{})
	for _, value := range keys {
		fmt.Println(value)
		mapKeys[value] = value
	}

	var mapFields = make(map[interface{}]map[string]interface{})
	var mapFieldNames = make(map[string]string)
	var sliceFieldNames = []string{}
	var sliceFieldTypes = []string{}
	// 获取 Sheet1 上所有单元格
	rows := f.GetRows("Sheet1")
	for index, row := range rows {
		// 第一行算是一种注释
		if index == 0 {
			for _, colCell := range row {
				if colCell == "" {
					log.Panic("fileName " + fileName + " has field empty 0!!!")
				}
				//fmt.Print(colCell)
				mapFieldNames[colCell] = colCell
			}
			continue
		}

		// 第二行是字段名字
		if index == 1 {
			for _, colCell := range row {
				if colCell == "" {
					log.Panic("fileName " + fileName + " has field empty 1!!!")
				}
				//fmt.Print(colCell)
				sliceFieldNames = append(sliceFieldNames, colCell)
			}
			continue
		}

		// 第三行是数据类型
		if index == 2 {
			for _, colCell := range row {
				if colCell == "" {
					log.Panic("fileName " + fileName + " has field empty 2!!!")
				}

				colCell = compressStr(colCell)
				//fmt.Print(colCell)
				sliceFieldTypes = append(sliceFieldTypes, colCell)
			}
			continue
		}

		oneMapFields := make(map[string]interface{})
		comKeys := []string{}
		for index1, colCell := range row {
			// 实际的值判断
			fieldName := sliceFieldNames[index1]
			if _, ok := mapKeys[fieldName]; ok {
				comKeys = append(comKeys, colCell)
			}

			switch sliceFieldTypes[index1] {
			case "int64", "int32", "int":
				ret, _ := strconv.Atoi(colCell)
				oneMapFields[fieldName] = ret
			case "float32":
				//ret, _ := strconv.Atoi(colCell)
				//strconv.FormatFloat(float64, 'E', -1, 32)
				ret, _ := strconv.ParseFloat(colCell, 32)
				oneMapFields[fieldName] = float32(ret)
			case "float64":
				ret, _ := strconv.ParseFloat(colCell, 64)
				oneMapFields[fieldName] = ret
			case "string":
				oneMapFields[fieldName] = colCell
			case "[]int":
				sli := strings.Split(colCell, ",")
				sliTemp := []int{}
				for _, value := range sli {
					ret, _ := strconv.Atoi(value)
					sliTemp = append(sliTemp, ret)
				}
				// 设置数组
				oneMapFields[fieldName] = sliTemp
			case "[]string":
				sli := strings.Split(colCell, "|")
				// 设置数组
				oneMapFields[fieldName] = sli
			case "map[string]string": // key1,value1|key2,value2

			}
		}
		//sort.Strings(comKeys)
		mapFields[strings.Join(comKeys, "_")] = oneMapFields
	}

	//fmt.Println(mapFields)

	return mapFields
}

type BaseI interface {
}

type EquipBase struct {
	ItemBase
}

type ItemBase struct {
	ID       int64
	Name     string
	Type     uint16
	Quality  uint8
	Ratio1   float32
	Ratio2   float64
	BufferID []int32
	Names    []string
}

func Listfunc(path string, f os.FileInfo, err error) error {
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
		gwlog.Fatalf("filepath.Walk() returned %v\n", err)
		return err
	}

	return nil
}

func ListFileFunc(p []string) {
	for index, value := range p {
		fmt.Println("Index = ", index, " Value = ", value)
		if index == 0 {
			Read(value, "ID")
		}
	}
}
