package skill

import "github.com/xiaonanln/goworld/engine/entity"

type Skill struct {
	entity.Entity
}

func (skill *Skill) DescribeEntityType(desc *entity.EntityTypeDesc) {
	desc.SetUseAOI(true, 100)
	desc.DefineAttr("name", "Client")
	desc.DefineAttr("lv", "Client", "persistent")
}
