package skill

import (
	"github.com/xiaonanln/goworld/engine/entity"
	"math"
)

func Distance(src, dest entity.Vector3) float64 {
	return math.Sqrt(math.Pow(float64(src.X-dest.X), float64(src.Y-dest.Y)))
}

// p2p类型技能
type PointSkill struct {
	Distance float64
}

func (this *PointSkill) IsInDistance(src, dest entity.Vector3) bool {
	x := src.X - dest.Y
	y := src.X - dest.Y

	if float64(x*x+y*y) > this.Distance*this.Distance {
		return false
	}

	return true
}

// 圆形技能
type CircleSkill struct {
	Distance float64
}

func (this *CircleSkill) IsInDistance(src, dest entity.Vector3) bool {
	x := src.X - dest.Y
	y := src.X - dest.Y

	if float64(x*x+y*y) > this.Distance*this.Distance {
		return false
	}

	return true
}

// 矩形技能
type RectSkill struct {
	Width  float64
	Height float64
}

func (this *RectSkill) isInRect(minx, miny, maxx, maxy float64, point entity.Vector3) bool {
	//判断点point的xy是否在矩形上下左右之间
	if float64(point.X) >= minx && float64(point.X) <= maxx && float64(point.Y) >= miny && float64(point.Y) <= maxy {
		return true
	}

	return false
}

func (this *RectSkill) ChangeAbsolute2Relative(originPoint, directionPoint, changePoint entity.Vector3) (ret entity.Vector3) {
	if originPoint == directionPoint {
		ret.X = changePoint.X - originPoint.X
		ret.Y = changePoint.Y - originPoint.Y
		return
	}

	//计算三条边
	a := Distance(directionPoint, changePoint)
	b := Distance(changePoint, originPoint)
	c := Distance(directionPoint, originPoint)

	cosA := (b*b + c*c - a*a) / 2 * b * c                       //余弦
	ret.X = entity.Coord(b * cosA)                              //相对坐标x
	ret.Y = entity.Coord(math.Sqrt(b*b - float64(ret.X*ret.X))) //相对坐标y

	return
}

func (this *RectSkill) IsInRect(originPoint, directionPoint, changePoint entity.Vector3) bool {
	//检测每一个角色是否在矩形内。
	ret := this.ChangeAbsolute2Relative(originPoint, directionPoint, changePoint) //相对坐标
	//skillWidth为图中宽度，skillLong为图中长度
	//skillWidth := 50.0 //矩形攻击区域的宽度
	//skillLong := 50.0  //矩形攻击区域的高度

	//宽度是被AB平分的，从A点开始延伸长度
	return this.isInRect(0, -this.Width/2, this.Height, this.Width/2, ret) //相对坐标下攻击范围
}

// 扇形
type FanShapedSkill struct {
	Radius float64
	Angle  float64
}

func (this *FanShapedSkill) ChangeXYToPolarCoordinate(p entity.Vector3) (r float64, angle float64) {
	r = math.Sqrt(float64(p.X*p.Y + p.Y*p.Y))                               //半径
	angle = math.Atan2(float64(p.Y), float64(p.X)) * float64(180) / math.Pi //计算出来的是弧度，转成角度，atan2的范围是-π到π之间

	a := angle + float64(360)
	angle = a - math.Floor(math.Mod(a, float64(360)))*float64(360)
	return
}

//|
//|
//|朝向 direction
//|位置 position
//|-----------
func (this *FanShapedSkill) ChangeAbsolute2Relative(originPoint, directionPoint entity.Vector3) entity.Vector3 {
	var rePoint entity.Vector3
	rePoint.X = directionPoint.X - originPoint.X
	rePoint.Y = directionPoint.Y - originPoint.Y
	return rePoint
}

// a-> 根据朝向获取目标列表
func (this *FanShapedSkill) IsInShape(attackerPoint, targetPoint, directionPoint entity.Vector3) bool {

	// 根据攻击点和朝向得到范围
	rePoint := this.ChangeAbsolute2Relative(attackerPoint, directionPoint)
	rePoint.Normalize()

	rePoint.Mul(entity.Coord(this.Radius))
	_, angle := this.ChangeXYToPolarCoordinate(rePoint)

	rePoint = this.ChangeAbsolute2Relative(attackerPoint, targetPoint)

	r1, angle1 := this.ChangeXYToPolarCoordinate(rePoint)

	if r1 > this.Radius {
		return false
	}

	if math.Abs(angle1-angle) > this.Angle/2 {
		return false
	}

	return true
}

// 叉积
func (this *FanShapedSkill) IsInShapeByMul(attackerPoint, targetPoint, directionPoint entity.Vector3) bool {
	// 先算半径
	dis := float64(attackerPoint.DistanceTo(targetPoint))
	if dis > this.Radius {
		return false
	}

	// 根据攻击点和朝向得到范围
	rePoint := this.ChangeAbsolute2Relative(attackerPoint, directionPoint)
	rePoint.Normalize()

	rePointTarget := this.ChangeAbsolute2Relative(attackerPoint, targetPoint)
	rePointTarget.Normalize()

	jiaodu := math.Acos(float64(rePoint.X*rePointTarget.X+rePoint.Y*rePointTarget.Y)) * 180 / math.Pi // 夹角大小

	if jiaodu > this.Angle/2 {
		return false
	}

	return true
}

// https://blog.csdn.net/u012175089/article/details/51048998
// https://blog.csdn.net/u012175089/article/details/50857990
// https://blog.csdn.net/u012175089/article/details/50850250
// https://blog.csdn.net/fagarine/article/details/91045292
// https://www.gameres.com/497722.html
