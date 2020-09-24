package inter

import (
	"github.com/xiaonanln/goworld/engine/common"
	"github.com/xiaonanln/goworld/engine/entity"
)

/*
@Time : 2020/9/22 21:42
@Author : Administrator
@File : IMonster
@Software: GoLand
@Version: 1.0.0
@Description:
*/

type IMonster interface {
	entity.IEntity

	Move(id common.EntityID) bool
	Shot()
	Hp() int
	HpMax() int
	// 获取最近的目标实体
	GetNearestTarget(typeName string) *entity.Entity
	Attack(id common.EntityID) bool
	Idle()
}
