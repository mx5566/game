package player

import (
	"fmt"
	pproto "github.com/golang/protobuf/proto"
	"github.com/xiaonanln/goworld"
	"github.com/xiaonanln/goworld/engine/common"
	"github.com/xiaonanln/goworld/engine/consts"
	"github.com/xiaonanln/goworld/engine/entity"
	"github.com/xiaonanln/goworld/engine/gwlog"
	"github.com/xiaonanln/goworld/examples/unity_demo/skill"
	"github.com/xiaonanln/goworld/proto"
	"strconv"
)

// Player 对象代表一名玩家
type Player struct {
	entity.Entity

	mgr *skill.SkillMgr
}

func (a *Player) DescribeEntityType(desc *entity.EntityTypeDesc) {
	desc.SetPersistent(true).SetUseAOI(true, 100)
	desc.DefineAttr("name", "AllClients", "Persistent")
	desc.DefineAttr("lv", "AllClients", "Persistent")
	desc.DefineAttr("hp", "AllClients")
	desc.DefineAttr("hpmax", "AllClients")
	desc.DefineAttr("action", "AllClients")
	desc.DefineAttr("spaceKind", "Persistent")
	desc.DefineAttr("skillId", "Persistent")
}

// OnCreated 在Player对象创建后被调用
func (a *Player) OnCreated() {
	a.Entity.OnCreated()
	a.setDefaultAttrs()
	a.mgr = new(skill.SkillMgr)
	a.mgr.Owner = a

	// 技能怎么load

	gwlog.DebugfE("Player OnCreated type[%s]", a.mgr.Owner.(*Player).TypeName)
}

// setDefaultAttrs 设置玩家的一些默认属性
func (a *Player) setDefaultAttrs() {
	a.Attrs.SetDefaultInt("spaceKind", 1)
	a.Attrs.SetDefaultStr("name", "noname")
	a.Attrs.SetDefaultInt("lv", 1)
	a.Attrs.SetDefaultInt("hp", 100)
	a.Attrs.SetDefaultInt("hpmax", 100)
	a.Attrs.SetDefaultStr("action", "idle")

	skills := new(entity.ListAttr)
	skills.AppendStr("100001")
	a.Attrs.SetDefaultListAttr("skillId", skills)

	a.SetClientSyncing(true)
}

// GetSpaceID 获得玩家的场景ID并发给调用者
func (a *Player) GetSpaceID(callerID common.EntityID) {
	a.Call(callerID, "OnGetPlayerSpaceID", a.ID, a.Space.ID)
}

func (p *Player) enterSpace(spaceKind int) {
	if p.Space.Kind == spaceKind {
		return
	}
	if consts.DEBUG_SPACES {
		gwlog.Infof("%s enter space from %d => %d", p, p.Space.Kind, spaceKind)
	}
	goworld.CallServiceShardKey("SpaceService", strconv.Itoa(spaceKind), "EnterSpace", p.ID, spaceKind)
}

// OnClientConnected is called when client is connected
func (a *Player) OnClientConnected() {
	gwlog.Infof("%s client connected", a)

	gwlog.TraceErrorEx("Player OnClientConnected kind[%d]", int(a.GetInt("spaceKind")))

	a.enterSpace(int(a.GetInt("spaceKind")))
}

// OnClientDisconnected is called when client is lost
func (a *Player) OnClientDisconnected() {
	gwlog.Infof("%s client disconnected", a)
	a.Destroy()
}

// EnterSpace_Client is enter space RPC for client
func (a *Player) EnterSpace_Client(kind int) {
	a.enterSpace(kind)
}

// DoEnterSpace is called by SpaceService to notify avatar entering specified space
func (a *Player) DoEnterSpace(kind int, spaceID common.EntityID) {
	// let the avatar enter space with spaceID
	a.EnterSpace(spaceID, entity.Vector3{})
}

//func (a *Player) randomPosition() entity.Vector3 {
//	minCoord, maxCoord := -400, 400
//	return entity.Vector3{
//		X: entity.Coord(minCoord + rand.Intn(maxCoord-minCoord)),
//		Y: 0,
//		Z: entity.Coord(minCoord + rand.Intn(maxCoord-minCoord)),
//	}
//}

// OnEnterSpace is called when avatar enters a space
func (a *Player) OnEnterSpace() {
	gwlog.Infof("%s ENTER SPACE %s", a, a.Space)
	a.SetClientSyncing(true)
}

func (a *Player) SetAction_Client(action string) {
	if a.GetInt("hp") <= 0 { // dead already
		return
	}

	a.Attrs.SetStr("action", action)
}

// test protobuf
func (a *Player) Test_Client() {
	p := &proto.Person{
		Name: "Jack",
		Age:  10,
		From: "China",
	}
	fmt.Println("原始数据:", p)

	// 序列化
	dataMarshal, err := pproto.Marshal(p)
	if err != nil {
		fmt.Println("proto.Unmarshal.Err: ", err)
		return
	}
	fmt.Println("编码数据:", dataMarshal)
	// 反序列化
	person := proto.Person{}
	err = pproto.Unmarshal(dataMarshal, &person)
	if err != nil {
		fmt.Println("proto.Unmarshal.Err: ", err)
		return
	}

	fmt.Printf("解码数据: 姓名：%s 年龄：%d 国籍：%s ", person.GetName(), person.GetAge(), person.GetFrom())
}

func (a *Player) ShootMiss_Client() {
	a.CallAllClients("Shoot")
}

func (a *Player) ShootHit_Client(victimID common.EntityID) {
	a.CallAllClients("Shoot")
	victim := a.Space.GetEntity(victimID)
	if victim == nil {
		gwlog.Warnf("Shoot %s, but monster not found", victimID)
		return
	}

	if victim.Attrs.GetInt("hp") <= 0 {
		return
	}

	gwlog.Infof("Shoot %s, monster hp %d", victimID, victim.Attrs.GetInt("hp"))
	monster := victim.I.(*Monster)
	monster.TakeDamage(50)
}

func (player *Player) TakeDamage(damage int64) {
	hp := player.GetInt("hp")
	if hp <= 0 {
		return
	}

	hp = hp - damage
	if hp < 0 {
		hp = 0
	}

	player.Attrs.SetInt("hp", hp)

	if hp <= 0 {
		// now player dead ...
		player.Attrs.SetStr("action", "death")
		player.SetClientSyncing(false)
		// triggle server logic
		// onDead
	}
}

func (player *Player) UserSkill_Client() {
	// skillID := 10001
	// targetID := 1001
	// targetID不存在就是直接空放

}
