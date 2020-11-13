package map_file

import (
	"encoding/json"
	"github.com/xiaonanln/goworld/engine/common"
	"github.com/xiaonanln/goworld/engine/entity"
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
	StartPos   Position3   `json:"start_pos"`
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

type Position3 struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
	Z int32 `json:"z"`
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

func (this *MapInfo) IsHasBlockPosition(x, y float64) bool {
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

func (this *MapInfo) IsValidPosition(x, y int32) bool {
	if x < 0 || y < 0 {
		return false
	}

	if x >= this.Width || y >= this.Height {
		return false
	}

	return true
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
	t.Pos.X = x
	t.Pos.Y = y
	t.W = m
}

// 格子坐标转为实际坐标
func GridToPos(pos Position, mapInfo *MapInfo) (vec entity.Vector3) {
	vec.X = entity.Coord(pos.X*mapInfo.TileWidth + mapInfo.StartPos.X)
	vec.Y = 0
	vec.Z = entity.Coord(-pos.Y*mapInfo.TileHeight + mapInfo.StartPos.Z)
	return
}

/*
   def navigate(self, start, end):
       start = HeapAstar.Node(None, int(round((start[0] + 30.0))), int(round((30.0-start[2]))))
       end = HeapAstar.Node(None, int(round((end[0] + 30.0))), int(round((30.0-end[2]))))
       self.path_list = HeapAstar.transform_path_list(HeapAstar.find_path(start, end))
*/
/*
# 输出坐标的转换
def transform_path_list(path_list):
    if path_list:
        print "crude path:", path_list
        return [(p[0]-30.0, 0, 30.0-p[1]) for p in path_list]
    else:
        return []
*/
// 解释
// 首先我们的阻挡的是二维的格子坐标
// 通过三位的坐标西(x,y,z)按照规则生成的 其中二维的格子是横向是X 纵向是Y
// 对于实际的坐标 我们x坐标转换为格子坐标的就是 (实际的x坐标 - + 开始的坐标X) / 单个格子的宽度设置
// 实际坐标转为格子坐标
// 0,1,1,1,1,1
// 1,1,1,1,1,1
// 1,1,1,1,1,1
// 1,1,1,1,1,1
// 1,1,1,1,1,1
// 1,1,1,1,1,1
// 上面左上角0对应的地图转换的开始坐标，这个开始坐标对应的地图的世界坐标不一定是(0,0,0)对应的可能是(-35,0,35)
// 要转的坐标是 (-33.2, 0, 35)
// X = -33.2 - (-35)(开始坐标) / 单个格子的宽(2) = 1.8 / 2 = 0
func PosToGrid(pos entity.Vector3, mapInfo *MapInfo) (xStartTile, yStartTile int32) {
	xStartTile, yStartTile = int32(math.Floor(float64(pos.X-entity.Coord(mapInfo.StartPos.X)))/float64(mapInfo.TileWidth)), int32(math.Round(float64(entity.Coord(mapInfo.StartPos.Z)-pos.Z))/float64(mapInfo.TileHeight))
	return
}

// -35~45(0,79)
