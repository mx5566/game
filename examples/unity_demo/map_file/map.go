package map_file

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
)

// map struct
type MapInfo struct {
	Name       string      `json:"name"`
	Uid        string      `json:"uid"`
	TileWidth  int32       `json:"tile_width"`
	TileHeight int32       `json:"tile_height"`
	Width      int32       `json:"width"`
	Height     int32       `json:"height"`
	MapBlock   [][]int     `json:"blocks"`
	MapObjects []MapObject `json:"objects"`
}

func (this *MapInfo) Init(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, this)
	if err != nil {
		return err
	}

	return nil
}

func (this *MapInfo) Load(name string) error {
	return nil
}

type Position struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

// map object
type MapObject struct {
	ID  int64    `json:"id"`
	Pos Position `json:"pos"`
}

func (this *MapInfo) LoadBlock() {

}

func (this *MapInfo) IsHasBlockGrid(x, y int32) bool {
	if x < 0 || y < 0 {
		return true
	}

	if x >= this.Width || y >= this.Height {
		return true
	}
	if this.MapBlock[x][y] == 1 {
		return true
	}

	return false
}

func (this *MapInfo) IsHasBlockPostion(x, y float64) bool {
	if x < 0 || y < 0 {
		return true
	}

	// 80    70 75 79 80 81
	width := int32(math.Floor(x / float64(this.TileWidth))) //float32(this.TileWidth * this.Width)
	height := int32(math.Floor(x / float64(this.TileHeight)))
	if width >= this.Width || height >= this.Height {
		return true
	}

	if this.MapBlock[width][height] == 1 {
		return true
	}

	return false
}

// load all map jsonfile
func LoadAllMaps() {

}

type Map struct {
	MapInfo
	// all grid
	Grids map[int]map[int]*Grid
}
