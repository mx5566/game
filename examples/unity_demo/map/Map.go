package _map

import "github.com/xiaonanln/goworld/engine/entity"

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

// block struct
type MapBlock struct {
}

// map object
type MapObject struct {
	ID  int64          `json:"id"`
	Pos entity.Vector3 `json:"pos"`
}
