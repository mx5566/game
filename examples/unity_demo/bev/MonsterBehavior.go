package bev

import (
	b3core "github.com/magicsea/behavior3go/core"
	"github.com/xiaonanln/goworld/examples/unity_demo/npc"
)

type IMonsterBehavior interface {
	Start()
	Update(dtime float32)
}

// monster behavior
type MonsterBehavior struct {
	monster *npc.Monster
	bb      *b3core.Blackboard //记录行为状态

}

func NewMonsterBehavior(monster *npc.Monster) *MonsterBehavior {
	return &MonsterBehavior{
		monster: monster,
		bb:      b3core.NewBlackboard(),
	}
}

func (ai *MonsterBehavior) Start() {

}

func (ai *MonsterBehavior) Update(dtime float32) {
	//更新行为树
	//tree := GetBevTree()
	//tree.Tick(ai.fighter, ai.bb)
}
