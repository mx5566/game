package skill

import (
	"github.com/xiaonanln/goworld"
	"github.com/xiaonanln/goworld/engine/common"
	"github.com/xiaonanln/goworld/examples/unity_demo/inter"
)

//
type SkillMgr struct {
	Skills map[uint32]*goworld.Entity
	Owner  inter.IPlayer
}

func (this *SkillMgr) Load() {
	// load skill entity
	// playerID$1001
}

func (this *SkillMgr) DefaultSkill(playerID common.EntityID, skillBaseID uint32) {
	skill, ok := this.Skills[skillBaseID]
	if ok {
		skill.Call(playerID, "LearnSkillResult", skillBaseID, -1)
		return
	}

	skill = goworld.CreateEntityLocallyByExternal("Skill", map[string]interface{}{common.BaseID: skillBaseID})
	if skill == nil {
		skill.Call(playerID, "LearnSkillResult", skillBaseID, -1)
		return
	}

	this.Skills[skillBaseID] = skill

	skill.Call(playerID, "LearnSkillResult", skillBaseID, 0)
}

func (this *SkillMgr) UseSkill(skillID uint64, targetID common.EntityID) {

}

func (this *SkillMgr) LearnSkill(playerID common.EntityID, skillBaseID uint32) {
	skill, ok := this.Skills[skillBaseID]
	if ok {
		skill.Call(playerID, "LearnSkillResult", skillBaseID, -1)
		return
	}

	skill = goworld.CreateEntityLocallyByExternal("Skill", map[string]interface{}{common.BaseID: skillBaseID})

	this.Skills[skillBaseID] = skill

	skill.Call(playerID, "LearnSkillResult", skillBaseID, 0)
}

func (this *SkillMgr) UpgradeSkill(skillID uint32) {
	skill, ok := this.Skills[skillID]
	if !ok {
		return
	}

	s := skill.I.(*Skill)

	s.Upgrade()
}

func (this *SkillMgr) PrintSkills() {
	for _, v2 := range this.Skills {
		v2.I.(*Skill).PrintSkill()
	}
}
