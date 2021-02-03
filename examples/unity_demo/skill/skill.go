package skill

import (
	"github.com/xiaonanln/goworld/engine/common"
	"github.com/xiaonanln/goworld/engine/entity"
	"github.com/xiaonanln/goworld/engine/gwlog"
	"github.com/xiaonanln/goworld/excelt"
)

type Skill struct {
	entity.Entity
}

func (skill *Skill) DescribeEntityType(desc *entity.EntityTypeDesc) {
	desc.DefineAttr("name", "Client")
	desc.DefineAttr("lv", "Client", "persistent")
	desc.DefineAttr(common.BaseID, "Client", "persistent")
}

func (skill *Skill) OnCreated() {
	skill.Entity.OnCreated()
	skill.setDefaultAttrs()

	gwlog.DebugfE("skill OnCreated %s %d", skill.ID, skill.GetInt(common.BaseID))
}

func (skill *Skill) SetBaseID() {

}

func (skill *Skill) GetBaseID() int64 {
	return skill.GetInt(common.BaseID)
}

func (skill *Skill) setDefaultAttrs() {
	skill.Attrs.SetDefaultInt("lv", 1)
}

func (skill *Skill) Upgrade() bool {
	skillBase := excelt.GetBase(excelt.SkillTableStr, skill.GetInt(common.BaseID))
	if skillBase == nil {
		return false
	}

	// is max level
	skill.Attrs.SetInt("lv", skill.Attrs.GetInt("lv")+1)

	// save
	// skill.Save()
	return true
}

func (skill *Skill) Save() {
	skill.Entity.Save()
}
