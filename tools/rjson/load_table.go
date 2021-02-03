package rjson

import (
	"fmt"
	"github.com/xiaonanln/goworld/engine/common"
	"github.com/xiaonanln/goworld/engine/gwioutil"
	"github.com/xiaonanln/goworld/engine/gwlog"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

//TODO 查找所有的TODO位置
var (
	ostype = runtime.GOOS // 获取系统类型
)

// 所有表的预定义字符串
var (
	ItemTableStr  = "item"
	EquipTableStr = "equip"
	NpcTableStr   = "npc"
	SkillTableStr = "skill"
)

// 加载table
func init() {
	//Load()
}

func Load() {
	var listpath = "."
	if ostype == "windows" {
		listpath += "\\tools\\rjson\\"
	} else if ostype == "linux" {
		listpath += "/tools/rjson/"
	}

	var filter gwioutil.FileFilter
	_ = filter.GetFileList(listpath, common.Json)

	list := filter.ListFile

	gwlog.DebugfE("json file load count[%v]", len(list))

	pathHead := "tools"

	if ostype == "windows" {
		pathHead += "\\rjson\\"
	} else if ostype == "linux" {
		pathHead += "/rjson/"
	}
	// TODO 按照示例添加表
	for _, path := range list {
		switch path {
		case pathHead + "equip.json":
			LoadEquipTable(path)
		case pathHead + "item.json":
			LoadItemTable(path)
		case pathHead + "npc.json":
			LoadNpcTable(path)
		default:
			fmt.Println("error path " + path)
		}
	}
}

func compressStr(str string) string {
	if str == "" {
		return ""
	}
	//匹配一个或多个空白符的正则表达式
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(str, "")
}

// 把key转换位字符串
func CombineKeys(keys ...interface{}) string {
	//sort.Strings(keys)
	fmt.Println(keys...)
	com := []string{}
	for _, key := range keys {
		switch key.(type) {
		case int, int32, int64, int8, int16:
			com = append(com, strconv.FormatInt(reflect.ValueOf(key).Int(), 10))
		case uint, uint32, uint64, uint16, uint8:
			com = append(com, strconv.FormatUint(reflect.ValueOf(key).Uint(), 10))
		case string:
			com = append(com, key.(string))
		default:
			fmt.Println("unkonw type "+reflect.TypeOf(key).String(), " ", key)
		}
	}
	return strings.Join(com, "_")
}

// 把key转换位字符串
func CombineKeysEx(keys []interface{}) string {
	//sort.Strings(keys)
	com := []string{}
	for _, key := range keys {
		switch key.(type) {
		case int, int32, int64, int8, int16:
			com = append(com, strconv.FormatInt(reflect.ValueOf(key).Int(), 10))
		case uint, uint32, uint64, uint16, uint8:
			com = append(com, strconv.FormatUint(reflect.ValueOf(key).Uint(), 10))
		case string:
			com = append(com, key.(string))
		default:
			gwlog.ErrorfE("unkonw type %s  key %v", reflect.TypeOf(key).String(), key)
		}
	}
	return strings.Join(com, "_")
}

////////////////////////////////////////////////////////////////////////////////////////
// 基本的表数据结构

var MapItemsBase map[interface{}]ItemBase
var MapNpcBase map[interface{}]NpcBase
var MapEquipsBase map[interface{}]EquipBase

type TableBase interface {
}

type ItemBase struct {
	ID       int64    `json:"ID"`
	Name     string   `json:"Name"`
	Type     uint16   `json:"Type"`
	Quality  uint8    `json:"Quality"`
	Ratio1   float32  `json:"Ratio1"`
	Ratio2   float64  `json:"Ratio2"`
	BufferID []int32  `json:"Ids"`
	Names    []string `json:"Names"`
}

type NpcBase struct {
	ID             int64   `json:"ID"`
	Name           string  `json:"Name"`
	Type           uint16  `json:"Type"`
	Level          uint16  `json:"Level"`
	Hp             int64   `json:"Hp"`
	AttackInter    int32   `json:"AttackInter"`
	AttackDistance float32 `json:"AttackDistance"`
}

type EquipBase struct {
	ItemBase
	// external attr
}

////////////////////////////////////////////////////////////////////////////
func LoadItemTable(path string) {
	fmt.Println("load table item !!!")
	j := NewJsonStruct()
	var ii []ItemBase

	err := j.Load(path, &ii)
	if err != nil {
		gwlog.Error(err)
		os.Exit(0)
	}

	MapItemsBase = make(map[interface{}]ItemBase)
	for _, value := range ii {
		primary := strconv.FormatInt(value.ID, 10)
		MapItemsBase[primary] = value
	}

	fmt.Println(MapItemsBase)
}

func LoadEquipTable(path string) {
	fmt.Println("load table equip !!!")

	j := NewJsonStruct()
	var ii []EquipBase

	err := j.Load(path, &ii)
	if err != nil {
		gwlog.Error(err)
		os.Exit(0)
	}

	MapEquipsBase = make(map[interface{}]EquipBase)
	for _, value := range ii {
		// 主键处理的逻辑
		// 如果有多个key 格式是 key1_key2
		primary := strconv.FormatInt(value.ID, 10)
		MapEquipsBase[primary] = value
	}

	fmt.Println(MapEquipsBase)

}

func LoadNpcTable(path string) {
	fmt.Println("load table npc !!!")

	j := NewJsonStruct()
	var ii []NpcBase

	err := j.Load(path, &ii)
	if err != nil {
		gwlog.Error(err)
		os.Exit(0)
	}

	MapNpcBase = make(map[interface{}]NpcBase)
	for _, value := range ii {
		// 主键处理的逻辑
		// 如果有多个key 格式是 key1_key2
		primary := strconv.FormatInt(value.ID, 10)
		MapNpcBase[primary] = value
	}

	fmt.Println(MapNpcBase)
}

// TODO: 需要加入对用的表返回
func GetBase(name string, keys ...interface{}) interface{} {
	if len(name) == 0 {
		return nil
	}

	gwlog.DebugfE("GetBase name[%s] key[%v]", name, keys)
	keyCom := CombineKeysEx(keys)
	switch name {
	case ItemTableStr:
		if base, ok := MapItemsBase[keyCom]; ok {
			return &base
		}
	case EquipTableStr:
		if base, ok := MapEquipsBase[keyCom]; ok {
			return &base
		}
	case NpcTableStr:
		if base, ok := MapNpcBase[keyCom]; ok {
			return &base
		}
	default:
		gwlog.ErrorfE("GetBase Error name %s", name)
	}

	return nil
}
