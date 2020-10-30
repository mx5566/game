package insent

import (
	"fmt"
	"strings"
)

type Null struct{}

var sensitiveWord = make(map[string]interface{})
var Set = make(map[string]Null)

const InvalidWords = " ,~,!,@,#,$,%,^,&,*,(,),_,-,+,=,?,<,>,.,—,，,。,/,\\,|,《,》,？,;,:,：,',‘,；,“,"

var InvalidWord = make(map[string]Null) //无效词汇，不参与敏感词汇判断直接忽略

type Insensitive struct {
}

func (this *Insensitive) Init() {
	words := strings.Split(InvalidWords, ",")
	for _, v := range words {
		InvalidWord[v] = Null{}
	}
}

// 日 本 人
// 日 本 狗 东 西
//

//生成违禁词集合
func (this *Insensitive) AddSensitiveToMap(set map[string]Null) {
	for key := range set {
		str := []rune(key)
		nowMap := sensitiveWord

		for i := 0; i < len(str); i++ {
			if _, ok := nowMap[string(str[i])]; !ok { //如果该key不存在，

				thisMap := make(map[string]interface{})
				thisMap["isEnd"] = false
				nowMap[string(str[i])] = thisMap
				nowMap = thisMap

				//fmt.Println("NO key2 ", nowMap, " sensitiveWord->", sensitiveWord)
				fmt.Printf("3 sensitiveWord[%p]  &sensitiveWord[%p] nowMap[%p] &nowMap[%p]\n", sensitiveWord, &sensitiveWord, nowMap, &nowMap)

			} else {
				nowMap = nowMap[string(str[i])].(map[string]interface{})
			}

			if i == len(str)-1 {
				nowMap["isEnd"] = true
			}
		}
	}
}

/*
func CheckSensitiveWord(txt string,  beginIndex int,  matchType int) int{
	//敏感词结束标识位：用于敏感词只有1位的情况
 	flag := false;
	//匹配标识数默认为0
 	matchFlag  := 0;
 	var word byte;
 	nowMap := sensitiveWord;
	for i := beginIndex; i < len(txt); i++ {
	word = txt.charAt(i);
	//获取指定key
	nowMap = (Map) nowMap.get(word);
	if (nowMap != null) {//存在，则判断是否为最后一个
	//找到相应key，匹配标识+1
	matchFlag++;
	//如果为最后一个匹配规则,结束循环，返回匹配标识数
	if ("1".equals(nowMap.get("isEnd"))) {
	//结束标志位为true
	flag = true;
	//最小规则，直接返回,最大规则还需继续查找
	if (MinMatchTYpe == matchType) {
	break;
	}
	}
	} else {//不存在，直接返回
	break;
	}
	}
	if (matchFlag < 2 || !flag) {//长度必须大于等于1，为词
	matchFlag = 0;
	}
	return matchFlag;
}*/

//敏感词汇转换为*
func ChangeSensitiveWords(txt string, sensitive map[string]interface{}) (word string) {
	str := []rune(txt)
	nowMap := sensitive
	start := -1
	tag := -1
	for i := 0; i < len(str); i++ {
		if _, ok := InvalidWord[(string(str[i]))]; ok || string(str[i]) == "," {
			continue
		}
		if thisMap, ok := nowMap[string(str[i])].(map[string]interface{}); ok {
			tag++
			if tag == 0 {
				start = i

			}
			isEnd, _ := thisMap["isEnd"].(bool)
			if isEnd {
				for y := start; y < i+1; y++ {
					str[y] = 42
				}
				nowMap = sensitive
				start = -1
				tag = -1

			} else {
				nowMap = nowMap[string(str[i])].(map[string]interface{})
			}

		} else {
			if start != -1 {
				i = start + 1
			}
			nowMap = sensitive
			start = -1
			tag = -1
		}
	}

	return string(str)
}

/*
func main() {
	words := strings.Split(InvalidWords,",")
	for _, v := range words {
		InvalidWord[v] = Null{}
	}
	Set["你妈逼的"] = Null{}
	Set["你妈"] = Null{}
	Set["日"] = Null{}
	AddSensitiveToMap(Set)
	text := "文明用语你&* 妈,逼的你这个狗日的，怎么这么傻啊。我也是服了，我日,这些话我都说不出口"
	fmt.Println(ChangeSensitiveWords(text,sensitiveWord))

}*/
