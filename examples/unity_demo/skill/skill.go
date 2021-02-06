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
	desc.SetPersistent(true)
	desc.DefineAttr("name", "Client")
	desc.DefineAttr("lv", "Client", "persistent")
	desc.DefineAttr(common.BaseID, "Client", "persistent")
	desc.DefineAttr("cooldown", "Client", "persistent")
}

func (skill *Skill) OnCreated() {
	skill.Entity.OnCreated()
	skill.setDefaultAttrs()

	gwlog.DebugfE("skill OnCreated %s %d", skill.ID, skill.GetInt(common.BaseID))
}

func (skill *Skill) SetBaseID() {

}

func (skill *Skill) GetBaseID() uint32 {
	return uint32(skill.GetInt(common.BaseID))
}

func (skill *Skill) setDefaultAttrs() {
	skill.Attrs.SetDefaultInt("lv", 1)
	skill.Attrs.SetDefaultInt("cooldown", 0)
}

func (skill *Skill) Upgrade() bool {
	skillBase := excelt.GetBase(excelt.SkillTableStr, skill.GetInt(common.BaseID))
	if skillBase == nil {
		return false
	}

	// is max level
	skill.Attrs.SetInt("lv", skill.Attrs.GetInt("lv")+1)

	// 技能会不会影响属性

	// save
	// skill.Save()
	return true
}

func (skill *Skill) Save() {
	skill.Entity.Save()
}

func (skill *Skill) PrintSkill() {
	gwlog.DebugfE("skillID[%d] skillLv[%d] cd[%d]", skill.Attrs.GetInt(common.BaseID), skill.Attrs.GetInt("lv"),
		skill.Attrs.GetInt("colldown"))
}
