package map_file

import (
	"fmt"
	"testing"
)

func TestMap_Init(t *testing.T) {
	var a = 1.1
	var b = 1.5
	var c = 1.6
	fmt.Println(int32(float32(a)))
	fmt.Println(int32(float32(b)))
	fmt.Println(int32(float32(c)))

	var m Map
	err := m.Init("./block.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(m)
}
