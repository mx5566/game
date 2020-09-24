package player

import (
	"github.com/xiaonanln/goworld/engine/common"
	"github.com/xiaonanln/goworld/examples/unity_demo/bev"
	mycommon "github.com/xiaonanln/goworld/examples/unity_demo/common"
	"github.com/xiaonanln/goworld/examples/unity_demo/inter"
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

	attackCD       time.Duration
	lastAttackTime time.Time

	ai inter.IMonsterBehavior
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
	monster.Entity.OnCreated()
	/*monster.ai = */ bev.NewMonsterBehavior( /*monster*/ )
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

	monster.attackCD = time.Second
	monster.lastAttackTime = time.Now()
}

func (monster *Monster) AI() {
	// 用behaviors3go来实现一个基本的ai模块判断
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
	return 2
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
	ent := monster.Space.GetEntity(id)
	if ent == nil {
		return false
	}

	if ent.TypeName != "Player" {
		return false
	}

	player := ent.I.(*Player)

	monster.CallAllClients("DisplayAttack", ent.ID)

	if ent.GetInt("hp") <= 0 {
		return false
	}

	player.TakeDamage(monster.GetDamage())
	monster.Attrs.SetStr("action", "attack")
	return true
}

func (monster *Monster) Idle() {
	monster.Attrs.SetStr("action", "idle")
	return
}

func (monster *Monster) Move(id common.EntityID) bool {
	ent := monster.Space.GetEntity(id)
	if ent == nil {
		return false
	}

	if monster.IsInterestedIn(ent) {
		return false
	}

	myPos := monster.GetPosition()
	direction := ent.GetPosition().Sub(myPos)
	direction.Y = 0

	t := direction.Normalized().Mul(monster.GetSpeed() * 30 / 1000.0)
	monster.SetPosition(myPos.Add(t))
	monster.FaceTo(ent)

	monster.Attrs.SetStr("action", "move")
	return true
}
