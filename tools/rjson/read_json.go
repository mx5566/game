package rjson

import (
	"encoding/json"
	"io/ioutil"
)

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

func (jst *JsonStruct) Load(filename string, v interface{}) (err error) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err1 := ioutil.ReadFile(filename)
	if err1 != nil {
		err = err1
		return
	}

	//fmt.Println(data)
	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}

	//fmt.Print(v)
	return
}
