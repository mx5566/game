package skill

import (
	"github.com/xiaonanln/goworld"
	"github.com/xiaonanln/goworld/engine/common"
	"github.com/xiaonanln/goworld/examples/unity_demo/inter"
	"github.com/xiaonanln/goworld/examples/unity_demo/player"
)

//
type SkillMgr struct {
	Skills map[uint64]*goworld.Entity
	Owner  inter.IPlayer
}

func (this *SkillMgr) Load() {
	// load skill entity
	// playerID$1001
}

func (this *SkillMgr) UseSkill(skillID uint64, targetID common.EntityID) {

}

func (this *SkillMgr) LearnSkill(skillBaseID uint64) {
	skill, ok := this.Skills[skillBaseID]
	if ok {
		this.Owner.(*player.Player).CallClient("LearnSkill", 1)
		return
	}

	skill = goworld.CreateEntityLocallyByExternal("Skill", map[string]interface{}{common.BaseID: skillBaseID})
	this.Skills[skillBaseID] = skill

	this.Owner.(*player.Player).SetSkillID(skill.ID, 0)
	this.Owner.(*player.Player).CallClient("LearnSkill", 0)
}

func (this *SkillMgr) UpgradeSkill(skillID uint64) {
	skill, ok := this.Skills[skillID]
	if !ok {
		return
	}

	s := skill.I.(*Skill)

	s.Upgrade()
}
