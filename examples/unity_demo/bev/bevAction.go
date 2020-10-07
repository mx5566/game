//ai节点动作
package bev

import (
	"fmt"
	b3 "github.com/magicsea/behavior3go"
	b3config "github.com/magicsea/behavior3go/config"
	b3core "github.com/magicsea/behavior3go/core"
	"github.com/xiaonanln/goworld/engine/common"
	"github.com/xiaonanln/goworld/engine/gwlog"
	mycommon "github.com/xiaonanln/goworld/examples/unity_demo/common"
	"github.com/xiaonanln/goworld/examples/unity_demo/inter"
	"math/rand"
	"time"
)

//---------------------------------------condition------------------------------------------------
//HaveTarget
type HaveTarget struct {
	b3core.Condition
	index string
}

func (this *HaveTarget) Initialize(setting *b3config.BTNodeCfg) {
	this.Condition.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
}

func (this *HaveTarget) OnTick(tick *b3core.Tick) b3.Status {
	id := tick.Blackboard.GetInt32(this.index, "", "")
	if id < 1 {
		return b3.FAILURE
	}

	/*f := tick.GetTarget().(inter.IMonster)
	_, b := f.gl.entitys[id]
	if !b {
		tick.Blackboard.Set(this.index, int32(0), "", "")
		return b3.FAILURE
	}*/
	return b3.SUCCESS
}

//CheckBool
type CheckBool struct {
	b3core.Condition
	keyname string
}

func (this *CheckBool) Initialize(setting *b3config.BTNodeCfg) {
	this.Condition.Initialize(setting)
	this.keyname = setting.GetPropertyAsString("keyname")
}

func (this *CheckBool) OnTick(tick *b3core.Tick) b3.Status {
	var b = tick.Blackboard.GetBool(this.keyname, "", "")
	if b {
		return b3.SUCCESS
	}
	return b3.FAILURE
}

//---------------------------------------actions------------------------------------------------

type RandWait struct {
	b3core.Action
	minTime int64
	maxTime int64
}

func (this *RandWait) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.minTime = setting.GetPropertyAsInt64("timemini")
	this.maxTime = setting.GetPropertyAsInt64("timemax")
}

func (this *RandWait) OnOpen(tick *b3core.Tick) {
	var startTime int64 = time.Now().UnixNano() / 1000000
	tick.Blackboard.Set("startTime", startTime, tick.GetTree().GetID(), this.GetID())
	end := this.minTime + rand.Int63n(this.maxTime-this.minTime)
	tick.Blackboard.Set("endTime", startTime+end, tick.GetTree().GetID(), this.GetID())
}

func (this *RandWait) OnTick(tick *b3core.Tick) b3.Status {
	var currTime = time.Now().UnixNano() / 1000000
	var endTime = tick.Blackboard.GetInt64("endTime", tick.GetTree().GetID(), this.GetID())

	if currTime > endTime {
		return b3.SUCCESS
	}

	return b3.RUNNING
}

//RandAction
type RandAction struct {
	b3core.Action
	index string
	min   float64
	max   float64
}

func (this *RandAction) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
	this.min = setting.GetProperty("min")
	this.max = setting.GetProperty("max")
}

func (this *RandAction) OnTick(tick *b3core.Tick) b3.Status {
	val := this.min + rand.Float64()*(this.max-this.min)
	tick.Blackboard.Set(this.index, val, "", "")
	return b3.SUCCESS
}

//RandMove
type RandMove struct {
	b3core.Action
}

func (this *RandMove) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
}

func (this *RandMove) OnTick(tick *b3core.Tick) b3.Status {
	f := tick.GetTarget().(inter.IMonster)
	f.Move("")
	return b3.SUCCESS
}

//Shoot
type Shoot struct {
	b3core.Action
}

func (this *Shoot) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
}

func (this *Shoot) OnTick(tick *b3core.Tick) b3.Status {
	f := tick.GetTarget().(inter.IMonster)
	f.Shot()
	return b3.SUCCESS
}

//TurnTarget
type TurnTarget struct {
	b3core.Action
	index string
}

func (this *TurnTarget) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
}

func (this *TurnTarget) OnTick(tick *b3core.Tick) b3.Status {
	id := tick.Blackboard.GetInt32(this.index, "", "")
	if id < 1 {
		return b3.FAILURE
	}
	/*f := tick.GetTarget().(inter.IMonster)
	tball, b := f.gl.entitys[id]
	if !b {
		tick.Blackboard.Set(this.index, uint32(0), "", "")
		return b3.FAILURE
	}

	v := tball.GetPos().Sub(*f.pos)
	a := v.AngleY() * 180 / math.Pi

	gwlog.Infof("%v TurnTarget angle=%v  v=%v,%v", f.id, a, v.X, v.Y)
	f.Move(a)*/

	return b3.SUCCESS
}

//FindItem
type FindItem struct {
	b3core.Action
	index string
	etype mycommon.EntityType
	dis   float32
}

func (this *FindItem) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
	this.etype = mycommon.EntityType(setting.GetPropertyAsInt("etype"))
	this.dis = float32(setting.GetProperty("range"))
}

func (this *FindItem) OnTick(tick *b3core.Tick) b3.Status {
	_ = tick.GetTarget().(inter.IMonster)
	tick.Blackboard.Set(this.index, int32(0), "", "")

	/*ball := f.FindNearItem(this.dis, this.etype)
	if nil == ball {
		return b3.FAILURE
	}

	id := ball.GetID()
	tick.Blackboard.Set(this.index, id, "", "")
	log.Info("FindItem %v dis:%v", id, this.dis)
	// var currTime int64 = time.Now().UnixNano() / 1000000
	// tick.Blackboard.Set("targetTime", currTime, "", "")
	*/

	return b3.SUCCESS
}

//SubTree
type SubTreeNode struct {
	b3core.Action
	sTree    *b3core.BehaviorTree
	treeName string
}

func (this *SubTreeNode) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.treeName = setting.GetPropertyAsString("treeName")
	this.sTree = CreateBevTree(this.treeName)
	if nil == this.sTree {
		gwlog.Errorf("SubTreeNode Get SubTree Failed, treeName: ", this.treeName)
	}
	gwlog.Infof("SubTreeNode::Initialize ", this, " treeName ", this.treeName)
}

func (this *SubTreeNode) OnTick(tick *b3core.Tick) b3.Status {
	if nil == this.sTree {
		return b3.ERROR
	}
	if tick.GetTarget() == nil {
		panic("unknow error!")
	}
	tar := tick.GetTarget()
	return this.sTree.Tick(tar, tick.Blackboard)
}

//随机
type RandomComposite struct {
	b3core.Composite
}

func (this *RandomComposite) OnOpen(tick *b3core.Tick) {
	tick.Blackboard.Set("runningChild", -1, tick.GetTree().GetID(), this.GetID())
}

func (this *RandomComposite) OnTick(tick *b3core.Tick) b3.Status {
	var child = tick.Blackboard.GetInt("runningChild", tick.GetTree().GetID(), this.GetID())
	if -1 == child {
		child = int(rand.Uint32()) % this.GetChildCount()
	}

	var status = this.GetChild(child).Execute(tick)
	if status == b3.RUNNING {
		tick.Blackboard.Set("runningChild", child, tick.GetTree().GetID(), this.GetID())
	} else {
		tick.Blackboard.Set("runningChild", -1, tick.GetTree().GetID(), this.GetID())
	}
	return status
}

//HpLess
type HpLess struct {
	b3core.Condition
	rate float32
}

func (this *HpLess) Initialize(setting *b3config.BTNodeCfg) {
	this.Condition.Initialize(setting)
	this.rate = float32(setting.GetProperty("rate"))
}

func (this *HpLess) OnTick(tick *b3core.Tick) b3.Status {
	f := tick.GetTarget().(inter.IMonster)
	rate := float32(f.Hp()) / float32(f.HpMax())
	if rate < this.rate {
		return b3.SUCCESS
	}
	return b3.FAILURE
}

//Parallel
type ParallelComposite struct {
	b3core.Composite
	failCond int //1有一个失败就失败 0全失败才失败
	succCond int //1有一个成功就成功 0全成功才成功
	//如果不能确定状态 那就有running返回running，不然失败
}

func (this *ParallelComposite) Initialize(setting *b3config.BTNodeCfg) {
	this.Composite.Initialize(setting)
	this.failCond = setting.GetPropertyAsInt("fail_cond")
	this.succCond = setting.GetPropertyAsInt("succ_cond")
}

func (this *ParallelComposite) OnTick(tick *b3core.Tick) b3.Status {
	var failCount int
	var succCount int
	var hasRunning bool
	for i := 0; i < this.GetChildCount(); i++ {
		var status = this.GetChild(i).Execute(tick)
		if status == b3.FAILURE {
			failCount++
		} else if status == b3.SUCCESS {
			succCount++
		} else {
			hasRunning = true
		}
	}
	if (this.failCond == 0 && failCount == this.GetChildCount()) || (this.failCond == 1 && failCount > 0) {
		return b3.FAILURE
	}
	if (this.succCond == 0 && succCount == this.GetChildCount()) || (this.succCond == 1 && succCount > 0) {
		return b3.FAILURE
	}
	if hasRunning {
		return b3.RUNNING
	}
	return b3.FAILURE
}

/////////////////////////////////myself////////////////////////////////
// idle

// find nearest target  if not find monster idling
// if find target and target distance than
type FindTarget struct {
	b3core.Action
	index    string // 值存储在blackboard >> index是key,对应的值是目标的ID
	typeName string // 目标类型 怪物、玩家、宠物等等

}

func (this *FindTarget) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
	this.typeName = setting.GetPropertyAsString("typeName")

	fmt.Println("Ai bev FindTarget Initialize -->", this.index, " ", this.typeName)

}

func (this *FindTarget) OnTick(tick *b3core.Tick) b3.Status {
	object := tick.GetTarget().(inter.IMonster)
	tick.Blackboard.Set(this.index, "", "", "")

	fmt.Println("Ai bev FindTarget OnTick -->", this.index, " ", this.typeName)


	nearestTarget := object.GetNearestTarget(this.typeName)
	if nearestTarget == nil {
		return b3.FAILURE
	}

	// set nearest target
	tick.Blackboard.Set(this.index, nearestTarget.ID, "", "")
	return b3.SUCCESS
}

// attack node
type AttackTarget struct {
	b3core.Action
	index string
}

func (this *AttackTarget) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")

	fmt.Println("Ai bev AttackTarget Initialize -->", this.index)

}

func (this *AttackTarget) OnTick(tick *b3core.Tick) b3.Status {
	object := tick.GetTarget().(inter.IMonster)
	id := tick.Blackboard.Get(this.index, "", "").(string)
	if id == "" {
		return b3.FAILURE
	}

	fmt.Println("Ai bev AttackTarget OnTick -->", this.index)

	ret := object.Attack(common.EntityID(id))

	if ret {
		return b3.SUCCESS
	}

	return b3.FAILURE
}

// move target
type MoveToTarget struct {
	b3core.Action
	index string
}

func (this *MoveToTarget) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")

	fmt.Println("Ai bev MoveToTarget Initialize -->", this.index)

}

func (this *MoveToTarget) OnTick(tick *b3core.Tick) b3.Status {
	object := tick.GetTarget().(inter.IMonster)
	id := tick.Blackboard.Get(this.index, "", "").(common.EntityID)
	if id.IsNil() {
		return b3.FAILURE
	}

	fmt.Println("Ai bev MoveToTarget OnTick id -->", id)


	ret := object.Move(id)

	if ret {
		return b3.SUCCESS
	}

	return b3.FAILURE
}

// idle node
type Idle struct {
	b3core.Action
}

func (this *Idle) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
}

func (this *Idle) OnTick(tick *b3core.Tick) b3.Status {
	object := tick.GetTarget().(inter.IMonster)
	object.Idle()
	return b3.SUCCESS
}
