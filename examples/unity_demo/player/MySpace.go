package player

import (
	"github.com/xiaonanln/goworld/engine/astar"
	"github.com/xiaonanln/goworld/engine/common"
	mycommon "github.com/xiaonanln/goworld/examples/unity_demo/common"
	"github.com/xiaonanln/goworld/examples/unity_demo/map_file"
	"strconv"
	"time"

	"github.com/xiaonanln/goTimer"
	"github.com/xiaonanln/goworld"
	"github.com/xiaonanln/goworld/engine/entity"
	"github.com/xiaonanln/goworld/engine/gwlog"
)

const (
	_SPACE_DESTROY_CHECK_INTERVAL = time.Minute * 5
)

// MySpace is the custom space type
type MySpace struct {
	goworld.Space // Space type should always inherit from entity.Space

	destroyCheckTimer entity.EntityTimerID

	// 地图信息 用来寻路阻挡判断
	Map *map_file.Map
}

// OnSpaceCreated is called when the space is created
func (space *MySpace) OnSpaceCreated() {
	// notify the SpaceService that it's ok
	space.EnableAOI(100)

	goworld.CallServiceShardKey("SpaceService", strconv.Itoa(space.Kind), "NotifySpaceLoaded", space.Kind, space.ID)
	space.AddTimer(time.Second*5, "DumpEntityStatus")
	space.AddTimer(time.Second*5, "SummonMonsters")

	// 根据kind 找到对应的地图的配置静态的信息 阻挡 地图基本信息
	// init map
	space.Map = new(map_file.Map)
	space.Map.Init(map_file.MapBaseInfoID[int32(space.Kind)])

}

func (space *MySpace) DumpEntityStatus() {
	space.ForEachEntity(func(e *entity.Entity) {
		gwlog.Debugf(">>> %s @ position %s, neighbors=%d", e, e.GetPosition(), len(e.InterestedIn))
	})
}

func (space *MySpace) SummonMonsters() {
	if space.CountEntities("Monster") < space.CountEntities("Player")*2 {
		space.CreateEntityByExternal("Monster", entity.Vector3{10.0, 0.0, 10.0}, map[string]interface{}{common.BaseID: 1})
	}
}

// OnEntityEnterSpace is called when entity enters space
func (space *MySpace) OnEntityEnterSpace(entity *entity.Entity) {
	if entity.TypeName == "Player" {
		space.onPlayerEnterSpace(entity)
	}
}

func (space *MySpace) onPlayerEnterSpace(entity *entity.Entity) {
	gwlog.Debugf("Player %s enter space %s, total avatar count %d", entity, space, space.CountEntities("Player"))
	space.clearDestroyCheckTimer()
}

// OnEntityLeaveSpace is called when entity leaves space
func (space *MySpace) OnEntityLeaveSpace(entity *entity.Entity) {
	if entity.TypeName == "Player" {
		space.onPlayerLeaveSpace(entity)
	}
}

func (space *MySpace) onPlayerLeaveSpace(entity *entity.Entity) {
	gwlog.Infof("Player %s leave space %s, left avatar count %d", entity, space, space.CountEntities("Player"))
	if space.CountEntities("Player") == 0 {
		// no avatar left, start destroying space
		space.setDestroyCheckTimer()
	}
}

func (space *MySpace) setDestroyCheckTimer() {
	if space.destroyCheckTimer != 0 {
		return
	}

	space.destroyCheckTimer = space.AddTimer(_SPACE_DESTROY_CHECK_INTERVAL, "CheckForDestroy")
}

// CheckForDestroy checks if the space should be destroyed
func (space *MySpace) CheckForDestroy() {
	avatarCount := space.CountEntities("Player")
	if avatarCount != 0 {
		gwlog.Panicf("Player count should be 0, but is %d", avatarCount)
	}

	goworld.CallServiceShardKey("SpaceService", strconv.Itoa(space.Kind), "RequestDestroy", space.Kind, space.ID)
}

func (space *MySpace) clearDestroyCheckTimer() {
	if space.destroyCheckTimer == 0 {
		return
	}

	space.CancelTimer(space.destroyCheckTimer)
	space.destroyCheckTimer = 0
}

// ConfirmRequestDestroy is called by SpaceService to confirm that the space
func (space *MySpace) ConfirmRequestDestroy(ok bool) {
	if ok {
		if space.CountEntities("Player") != 0 {
			gwlog.Panicf("%s ConfirmRequestDestroy: avatar count is %d", space, space.CountEntities("Player"))
		}
		space.Destroy()
	}
}

func (space *MySpace) FindPathA(start, dest entity.Vector3) []*map_file.Grid {
	xStartTile, yStartTile := map_file.PosToGrid(start, &space.Map.MapInfo)
	xDestTile, yDestTile := map_file.PosToGrid(dest, &space.Map.MapInfo)

	// 防止越界
	if !space.Map.IsValidPosition(xStartTile, yStartTile) || !space.Map.IsValidPosition(xDestTile, yDestTile) {
		return nil
	}

	s := space.Map.Grids[xStartTile][yStartTile]
	e := space.Map.Grids[xDestTile][yDestTile]

	// dist 就是图数据结构从a->b的边的权重累计
	p, _, found := astar.Path(s, e)
	if !found {
		return nil
	}

	ret := []*map_file.Grid{}
	for _, v := range p {
		ret = append(ret, v.(*map_file.Grid))
	}

	return ret
}

// OnGameReady is called when the game server is ready
func (space *MySpace) OnGameReady() {
	timer.AddCallback(time.Millisecond*1000, checkServerStarted)
}

func checkServerStarted() {
	ok := isAllServicesReady()
	gwlog.Infof("checkServerStarted: %v", ok)
	if ok {
		onAllServicesReady()
	} else {
		timer.AddCallback(time.Millisecond*1000, checkServerStarted)
	}
}

func isAllServicesReady() bool {
	for _, serviceName := range mycommon.SERVICE_NAMES {
		if !goworld.CheckServiceEntitiesReady(serviceName) {
			gwlog.Infof("%s entities are not ready ...", serviceName)
			return false
		}
	}
	return true
}

func onAllServicesReady() {
	gwlog.Infof("All services are ready!")
}
