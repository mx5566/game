package map_file

import (
	"fmt"
	"testing"
)

func TestMap_Init(t *testing.T) {
	var m MapInfo
	err := m.Init("./block.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(m)
}
