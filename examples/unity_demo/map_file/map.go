package map_file

import (
	"encoding/json"
	"github.com/xiaonanln/goworld/engine/common"
	"github.com/xiaonanln/goworld/engine/gwioutil"
	"github.com/xiaonanln/goworld/engine/gwlog"
	"io/ioutil"
	"math"
	"os"
	"runtime"
)

// Kind* constants refer to tile kinds for input and output.
const (
	// KindPlain (.) is a plain tile with a movement cost of 1.
	KindPlain = iota
)

// KindCosts map tile kinds to movement costs.
var KindCosts = map[int]float64{
	KindPlain: 1.0,
}

// 名字是key
var MapBaseInfo = make(map[string]MapInfo)

// Uid是key
var MapBaseInfoID = make(map[int32]MapInfo)

// map struct
type MapInfo struct {
	Name       string      `json:"name"`
	Uid        string      `json:"uid"`
	ID         int32       `json:"id"`
	TileWidth  int32       `json:"tile_width"`
	TileHeight int32       `json:"tile_height"`
	Width      int32       `json:"width"`
	Height     int32       `json:"height"`
	MapBlock   [][]int32   `json:"blocks"`
	MapObjects []MapObject `json:"objects"`
}

func init() {
	LoadAllMapRes()
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

type Position struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

// map object
type MapObject struct {
	ID  int64    `json:"id"`
	Pos Position `json:"pos"`
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
func LoadAllMapRes() {
	ostype := runtime.GOOS // 获取系统类型

	var splicing = ""
	if ostype == "windows" {
		splicing = "\\examples\\unity_demo\\map_file\\"
	} else if ostype == "linux" {
		splicing = "/examples/unity_demo/map_file/"
	}

	var listpath = "." + splicing

	var filter gwioutil.FileFilter
	_ = filter.GetFileList(listpath, common.Json)

	gwlog.DebugfE("map file load cout[%d]", len(filter.ListFile))

	var m MapInfo
	for _, path := range filter.ListFile {
		err := m.Init(path)
		if err != nil {
			gwlog.PanicfE("load file err file[%s] [%s]", path, err.Error())
		}

		MapBaseInfo[m.Name] = m
		MapBaseInfoID[m.ID] = m
	}
}

type Map struct {
	MapInfo
	// all grid
	Grids map[int32]map[int32]*Grid
}

// 初始化地图配置信息 A*Star
func (m *Map) Init(mapInfo MapInfo) {
	m.MapInfo = mapInfo
	m.Grids = make(map[int32]map[int32]*Grid)
	// 填充所有的格子 用于寻路
	for y, row := range m.MapInfo.MapBlock {
		for x, _ := range row {
			m.SetTile(&Grid{}, int32(x), int32(y))
		}
	}
}

// Tile gets the tile at the given coordinates in the world.
func (m *Map) Tile(x, y int32) *Grid {
	if m.Grids[x] == nil {
		return nil
	}
	return m.Grids[x][y]
}

// SetTile sets a tile at the given coordinates in the world.
func (m *Map) SetTile(t *Grid, x, y int32) {
	if m.Grids[x] == nil {
		m.Grids[x] = map[int32]*Grid{}
	}
	m.Grids[x][y] = t
	t.pos.X = x
	t.pos.Y = y
	t.W = m
}
