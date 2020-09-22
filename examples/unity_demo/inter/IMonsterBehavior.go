package inter

/*
@Time : 2020/9/22 21:31
@Author : Administrator
@File : IMonsterBehavior
@Software: GoLand
@Version: 1.0.0
@Description:
*/
type IMonsterBehavior interface {
	Start()
	Update(dtime float32)
}