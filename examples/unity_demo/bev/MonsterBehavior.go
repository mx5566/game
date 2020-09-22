package bev

import (
	b3core "github.com/magicsea/behavior3go/core"
	"github.com/xiaonanln/goworld/examples/unity_demo/inter"
)

// monster behavior
type MonsterBehavior struct {
	monster inter.IMonster
	bb      *b3core.Blackboard //记录行为状态

}

func NewMonsterBehavior(monster inter.IMonster) inter.IMonsterBehavior {
	return &MonsterBehavior{
		monster: monster,
		bb:      b3core.NewBlackboard(),
	}
}

func (ai *MonsterBehavior) Start() {

}

func (ai *MonsterBehavior) Update(dtime float32) {
	//更新行为树
	tree := GetBevTree()
	tree.Tick(ai.monster, ai.bb)
}
