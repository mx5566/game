package skill

import (
	"github.com/xiaonanln/goworld/engine/common"
)

//
type SkillMgr struct {
	Skills map[uint64]Ski
}

func (this *SkillMgr) Init() {

}

func (this *SkillMgr) UseSkill(skillID uint64, targetID common.EntityID) {

}
