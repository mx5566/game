package _map

// map struct
type Map struct {
	Name       string `json:"name"`
	Uid        string `json:"uid"`
	TileWidth  int32  `json:"tile_width"`
	TileHeight int32  `json:"tile_height"`
	Width      int32  `json:"width"`
	Height     int32  `json:"height"`
	MapObjects []MapObject
}

func (this *Map) Init() {

}

// block struct
type MapBlock struct {
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
