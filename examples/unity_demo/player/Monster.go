package player

import (
	"github.com/xiaonanln/goworld/engine/common"
	"github.com/xiaonanln/goworld/examples/unity_demo/bev"
	mycommon "github.com/xiaonanln/goworld/examples/unity_demo/common"
	"github.com/xiaonanln/goworld/examples/unity_demo/inter"
	"github.com/xiaonanln/goworld/examples/unity_demo/map_file"
	"github.com/xiaonanln/goworld/excelt"
	"strings"
	"time"

	"github.com/xiaonanln/goworld/engine/entity"
	"github.com/xiaonanln/goworld/engine/gwlog"
)

// Monster type
type Monster struct {
	entity.Entity   // Entity type should always inherit entity.Entity
	movingToTarget  *entity.Entity
	attackingTarget *entity.Entity
	lastTickTime    time.Time

	lastAttackTime time.Time
	lastMoveTime   time.Time

	// new
	ai inter.IMonsterBehavior

	// last move target
	moveTarget *TargetEntity
}

// last move to target
type TargetEntity struct {
	Target       *entity.Entity   // move to target
	LastMoveTime time.Time        // last move time
	Paths        []*map_file.Grid // last find path record
}

func (monster *Monster) DescribeEntityType(desc *entity.EntityTypeDesc) {
	desc.SetUseAOI(true, 100)
	desc.DefineAttr("name", "AllClients")
	desc.DefineAttr("lv", "AllClients")
	desc.DefineAttr("hp", "AllClients")
	desc.DefineAttr("hpmax", "AllClients")
	desc.DefineAttr("action", "AllClients")
}

func (monster *Monster) OnCreated() {
	// createEntity 里面基本数据已经预先加载完毕

	monster.Entity.OnCreated()

	monster.ai = bev.NewMonsterBehavior(monster)
	gwlog.DebugfE("monster OnCreated %s %d", monster.ID, monster.GetInt(common.BaseID))
}

func (monster *Monster) OnEnterSpace() {
	monster.setDefaultAttrs()

	monster.AddTimer(time.Millisecond*100, "AI")
	monster.lastTickTime = time.Now()
	monster.AddTimer(time.Millisecond*30, "Tick")
}

func (monster *Monster) setDefaultAttrs() {
	monster.Attrs.SetDefaultStr("name", "minion")
	monster.Attrs.SetDefaultInt("lv", 1)
	monster.Attrs.SetDefaultInt("hpmax", 100)
	monster.Attrs.SetDefaultInt("hp", 100)
	monster.Attrs.SetDefaultStr("action", "idle")

	monster.lastAttackTime = time.Now()
	monster.lastMoveTime = time.Now()
}

func (monster *Monster) SetBaseID(id int64) {

}

func (monster *Monster) AI() {
	// 用behaviors3go来实现一个基本的ai模块判断
	//fmt.Println("monster Ai start")
	dtime := float32(mycommon.FRAME_TIME) / float32(1000)
	monster.ai.Update(dtime)

	/*var nearestPlayer *entity.Entity
	for entity := range monster.InterestedIn {

		if entity.TypeName != "Player" {
			continue
		}

		if entity.GetInt("hp") <= 0 {
			// dead
			continue
		}

		if nearestPlayer == nil || nearestPlayer.DistanceTo(&monster.Entity) > entity.DistanceTo(&monster.Entity) {
			nearestPlayer = entity
		}
	}

	if nearestPlayer == nil {
		monster.Idling()
		return
	}

	if nearestPlayer.DistanceTo(&monster.Entity) > monster.GetAttackRange() {
		monster.MovingTo(nearestPlayer)
	} else {
		monster.Attacking(nearestPlayer)
	}*/
}

func (monster *Monster) Tick() {
	return
	/*
		if monster.attackingTarget != nil && monster.IsInterestedIn(monster.attackingTarget) {
			now := time.Now()
			if !now.Before(monster.lastAttackTime.Add(monster.attackCD)) {
				monster.FaceTo(monster.attackingTarget)
				monster.attack(monster.attackingTarget.I.(*Player))
				monster.lastAttackTime = now
			}
			return
		}

		if monster.movingToTarget != nil && monster.IsInterestedIn(monster.movingToTarget) {
			mypos := monster.GetPosition()
			direction := monster.movingToTarget.GetPosition().Sub(mypos)
			direction.Y = 0

			t := direction.Normalized().Mul(monster.GetSpeed() * 30 / 1000.0)
			monster.SetPosition(mypos.Add(t))
			monster.FaceTo(monster.movingToTarget)
			return
		}*/

}

func (monster *Monster) GetSpeed() entity.Coord {
	return 10
}

func (monster *Monster) GetAttackRange() entity.Coord {
	return 3
}

func (monster *Monster) Idling() {
	if monster.movingToTarget == nil && monster.attackingTarget == nil {
		return
	}

	monster.movingToTarget = nil
	monster.attackingTarget = nil
	monster.Attrs.SetStr("action", "idle")
}

func (monster *Monster) MovingTo(player *entity.Entity) {
	if monster.movingToTarget == player {
		// moving target not changed
		return
	}

	monster.movingToTarget = player
	monster.attackingTarget = nil
	monster.Attrs.SetStr("action", "move")
}

func (monster *Monster) Attacking(player *entity.Entity) {
	if monster.attackingTarget == player {
		return
	}

	monster.movingToTarget = nil
	monster.attackingTarget = player
	monster.Attrs.SetStr("action", "attack")
}

func (monster *Monster) attack(player *Player) {
	monster.CallAllClients("DisplayAttack", player.ID)

	if player.GetInt("hp") <= 0 {
		return
	}

	player.TakeDamage(monster.GetDamage())
}

func (monster *Monster) GetDamage() int64 {
	return 10
}

func (monster *Monster) TakeDamage(damage int64) {
	hp := monster.GetInt("hp")
	hp = hp - damage
	if hp < 0 {
		hp = 0
	}

	monster.Attrs.SetInt("hp", hp)
	gwlog.Infof("%s TakeDamage %d => hp=%d", monster, damage, hp)
	if hp <= 0 {
		monster.Attrs.SetStr("action", "death")
		monster.Destroy()
	}
}

////////////////////////////////new add//////////////////////////

func (monster *Monster) Shot() {
	gwlog.Debugf("Shot Test")
}

func (monster *Monster) Hp() int64 {
	return monster.GetInt("hp")
}

func (monster *Monster) HpMax() int64 {
	return monster.GetInt("hpmax")
}

func (monster *Monster) GetNearestTarget(typeName string) *entity.Entity {
	var nearestPlayer *entity.Entity
	for ent := range monster.InterestedIn {

		// fast than > < != ==
		if strings.Compare(ent.TypeName, typeName) != 0 {
			continue
		}

		if ent.GetInt("hp") <= 0 {
			// dead
			continue
		}

		if nearestPlayer == nil || nearestPlayer.DistanceTo(&monster.Entity) > ent.DistanceTo(&monster.Entity) {
			nearestPlayer = ent
		}
	}

	return nearestPlayer
}

func (monster *Monster) Attack(id common.EntityID) bool {
	//gwlog.DebugfE("monster Attack ", id)
	ent := monster.Space.GetEntity(id)
	if ent == nil {
		return false
	}

	if ent.TypeName != "Player" {
		return false
	}

	if ent.GetInt("hp") <= 0 {
		return false
	}

	// 查表去获取表数据
	npcBase := excelt.GetBase(excelt.NpcTableStr, monster.GetInt(common.BaseID))
	// 没有找到对应的怪物ID
	if npcBase == nil {
		return false
	}

	now := time.Now()
	// 攻击间隔判断 毫秒
	delta := time.Duration(npcBase.(*excelt.NpcBase).AttackInter)
	if now.Before(monster.lastAttackTime.Add(delta * time.Millisecond)) {
		return false
	}

	// reset attacktime
	monster.lastAttackTime = now

	player := ent.I.(*Player)

	monster.CallAllClients("DisplayAttack", ent.ID)

	player.TakeDamage(monster.GetDamage())
	monster.Attrs.SetStr("action", "attack")
	return true
}

func (monster *Monster) Idle() {
	//gwlog.DebugfE("monster Idle ")

	monster.Attrs.SetStr("action", "idle")
	return
}

func (monster *Monster) Move(id common.EntityID) bool {
	//gwlog.DebugfE("monster Move1 %v", id)

	ent := monster.Space.GetEntity(id)
	if ent == nil {
		return false
	}

	if !monster.IsInterestedIn(ent) {
		return false
	}

	if ent.GetInt("hp") <= 0 {
		return false
	}

	if monster.Hp() <= 0 {
		return false
	}

	myPos := monster.GetPosition()
	dest := ent.GetPosition()

	// 二者的距离小于某个值
	if myPos.DistanceTo(dest) < 0.9 {
		return true
	}

	// astar 寻路
	paths := monster.Space.I.(*MySpace).FindPathA(myPos, dest)
	if paths == nil || len(paths) == 0 {
		return false
	}

	// *map_file.Grid
	grid := paths[0]
	gwlog.DebugfE("move grid %v %v", grid, map_file.GridToPos(grid.Pos, &monster.Space.I.(*MySpace).Map.MapInfo))

	// 覆盖
	dest = map_file.GridToPos(grid.Pos, &monster.Space.I.(*MySpace).Map.MapInfo)

	if myPos.DistanceTo(dest) < 0.9 {
		paths = paths[1:]

		if len(paths) > 0 {
			grid = paths[0]
			dest = map_file.GridToPos(grid.Pos, &monster.Space.I.(*MySpace).Map.MapInfo)
		} else {
			return false
		}
	}

	direction := dest.Sub(myPos)
	direction.Y = 0

	t := direction.Normalized().Mul(monster.GetSpeed() * 30 / 1000.0)

	monster.SetPosition(myPos.Add(t))
	monster.FaceTo(ent)

	monster.Attrs.SetStr("action", "move")

	var target TargetEntity
	target.LastMoveTime = time.Now()
	target.Paths = paths
	target.Target = ent

	monster.moveTarget = &target

	return true
}

// 3 异常错误
// 2 对象距离大于设定的距离
// 1 距离小于等于设定的距离
func (monster *Monster) CheckDistance(id common.EntityID) int {
	ent := monster.Space.GetEntity(id)
	if ent == nil {
		return 3
	}

	myPos := monster.GetPosition()
	dest := ent.GetPosition()

	// 查表去获取表数据
	npcBase := excelt.GetBase(excelt.NpcTableStr, monster.GetInt(common.BaseID))
	// 没有找到对应的怪物ID
	if npcBase == nil {
		return 3
	}

	dis := float32(myPos.DistanceTo(dest))
	// 二者的距离大于配置的距离
	if dis > npcBase.(*excelt.NpcBase).AttackDistance {
		return 2
	}

	// <=
	return 1
}

func (monster *Monster) CheckHpgte() bool {

	hp := monster.Hp()

	npcBase := excelt.GetBase(excelt.NpcTableStr, monster.GetInt(common.BaseID))
	// 没有找到对应的怪物ID
	if npcBase == nil {
		return false
	}

	base := int64(100) // npcBase.(*excelt.NpcBase).AttackDistance

	return hp >= base
}

func (monster *Monster) CheckHplt() bool {
	hp := monster.Hp()

	npcBase := excelt.GetBase(excelt.NpcTableStr, monster.GetInt(common.BaseID))
	// 没有找到对应的怪物ID
	if npcBase == nil {
		return false
	}

	base := int64(100) // npcBase.(*excelt.NpcBase).AttackDistance

	return hp < base
}

func (monster *Monster) Flee() bool {

	return true
}
