package skill

import (
	"github.com/xiaonanln/goworld"
	"github.com/xiaonanln/goworld/engine/common"
	"github.com/xiaonanln/goworld/examples/unity_demo/inter"
)

//
type SkillMgr struct {
	Skills map[uint64]*Skill
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

	skill = goworld.CreateEntityLocally("Skill") // 创建一个Player对象然后立刻销毁，产生一次存盘

}

func (this *SkillMgr) UpgradeSkill(skillID uint64) {
	skill, ok := this.Skills[skillID]
	if !ok {
		return
	}

	skill.Upgrade()
}
