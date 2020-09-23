package bev

// monster behavior
type MonsterBehavior struct {
	//npc inter.IMonster
	//bb  *b3core.Blackboard //记录行为状态

}

func NewMonsterBehavior( /*monster inter.IMonster*/ ) /*inter.IMonsterBehavior */ {
	/*&MonsterBehavior{
		//monster: monster,
		//bb: b3core.NewBlackboard(),
	}*/
}

func (ai *MonsterBehavior) Start() {

}

func (ai *MonsterBehavior) Update(dtime float32) {
	//更新行为树
	//tree := GetBevTree()
	//tree.Tick(ai.monster, ai.bb)
}
