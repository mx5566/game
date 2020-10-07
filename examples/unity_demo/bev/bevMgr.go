//ai全局管理
package bev

import (
	b3 "github.com/magicsea/behavior3go"
	b3config "github.com/magicsea/behavior3go/config"
	b3core "github.com/magicsea/behavior3go/core"
	b3loader "github.com/magicsea/behavior3go/loader"
	"github.com/xiaonanln/goworld/engine/gwlog"
	"os"
)

func InitBev() {
	str,_ := os.Getwd()

	gwlog.Debugf("before----------------- %s", "start bev mgr" + str)
	mapTrees = make(map[string]*b3core.BehaviorTree)


	bevMainTree = CreateBevTree(str + "/examples/unity_demo/bev/monster.json")
	if bevMainTree == nil {
		gwlog.DebugfE("CreateBevTree error-----------------------")
	}

	gwlog.Debugf("after ----------------- %s", bevMainTree.GetTitile())
}

//主树
var bevMainTree *b3core.BehaviorTree

func GetBevTree() *b3core.BehaviorTree {
	return bevMainTree
}

//创建一个行为树
var mapTrees map[string]*b3core.BehaviorTree

func CreateBevTree(name string) *b3core.BehaviorTree {
	gwlog.Debugf("CreateBevTree %s", name)
	b, ok := mapTrees[name]
	if ok {
		return b
	}

	gwlog.Infof("create tree:%v", name)

	config, ok := b3config.LoadTreeCfg(name)
	if !ok {
		gwlog.Fatal("LoadTreeCfg fail:" + name)
	}
	extMaps := createExtStructMaps()
	tree := b3loader.CreateBevTreeFromConfig(config, extMaps)
	tree.Print()
	mapTrees[name] = tree
	return tree
}

//自定义的节点
func createExtStructMaps() *b3.RegisterStructMaps {
	st := b3.NewRegisterStructMaps()
	//actions
	st.Register("Rand", &RandAction{})
	st.Register("RandWait", &RandWait{})
	st.Register("TurnTarget", &TurnTarget{})
	st.Register("RandMove", &RandMove{})
	st.Register("Shoot", &Shoot{})
	st.Register("FindItem", &FindItem{})
	st.Register("SubTree", &SubTreeNode{})

	st.Register("FindTarget", &FindTarget{})
	st.Register("AttackTarget", &AttackTarget{})
	st.Register("MoveToTarget", &MoveToTarget{})
	st.Register("Idle", &Idle{})

	//conditions
	st.Register("HaveTarget", &HaveTarget{})
	st.Register("CheckBool", &CheckBool{})
	st.Register("HpLess", &HpLess{})

	//composite
	st.Register("Random", &RandomComposite{})
	st.Register("Parallel", &ParallelComposite{})
	return st
}
