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
	desc.SetUseAOI(true, 100)
	desc.DefineAttr("name", "Client")
	desc.DefineAttr("lv", "Client", "persistent")
	desc.DefineAttr("id", "Client", "persistent")
}

func (skill *Skill) OnCreated() {
	skill.Entity.OnCreated()

	gwlog.DebugfE("skill OnCreated %s %d", skill.ID, skill.GetInt(common.BaseID))
}

func (skill *Skill) setDefaultAttrs() {
	skill.Attrs.SetDefaultInt("lv", 1)
}

func (skill *Skill) Upgrade() bool {
	skillBase := excelt.GetBase(excelt.SkillTableStr, skill.GetInt(common.BaseID))
	if skillBase == nil {
		return false
	}

	// save
	skill.Save()
	return true
}

func (skill *Skill) Save() {

}
