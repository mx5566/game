package insent

import (
	"strings"
	"testing"
)

func TestChangeSensitiveWords(t *testing.T) {
	words := strings.Split(InvalidWords, ",")
	for _, v := range words {
		InvalidWord[v] = Null{}
	}
	Set["你妈逼的"] = Null{}
	//Set["你妈"] = Null{}
	//Set["日"] = Null{}

	var i Insensitive

	i.AddSensitiveToMap(Set)

	t.Log(sensitiveWord)

	//text := "文明用语你&* 妈,逼的你这个狗日的，怎么这么傻啊。我也是服了，我日,这些话我都说不出口"
	//fmt.Println(ChangeSensitiveWords(text, sensitiveWord))
}
