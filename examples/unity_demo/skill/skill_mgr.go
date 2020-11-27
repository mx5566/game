package skill

import (
	"github.com/xiaonanln/goworld"
	"github.com/xiaonanln/goworld/engine/common"
	"github.com/xiaonanln/goworld/examples/unity_demo/inter"
)

//
type SkillMgr struct {
	Skills map[uint64]*goworld.Entity
	Owner  inter.IPlayer
}

func (this *SkillMgr) Init() {

}

func (this *SkillMgr) UseSkill(skillID uint64, targetID common.EntityID) {

}

func (this *SkillMgr) LearnSkill(skillID uint64) {
	skill, ok := this.Skills[skillID]
	if ok {
		return
	}

	skill = goworld.CreateEntityLocallyByExternal("Skill", map[string]interface{}{common.BaseID: 1})
	this.Skills[skillID] = skill

}

func (this *SkillMgr) UpgradeSkill(skillID uint64) {
	skill, ok := this.Skills[skillID]
	if !ok {
		return
	}

	entityID := skill.ID

	skill.Upgrade()
}
