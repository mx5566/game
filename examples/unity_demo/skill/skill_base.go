package skill

import "math"

type Point struct {
	X float64
	Y float64
}

func Distance(src, dest Point) float64 {
	return math.Sqrt(math.Pow(src.X-dest.X, src.Y-dest.Y))
}

// p2p类型技能
type PointSkill struct {
	Distance float64
}

func (this *PointSkill) IsInDistance(src, dest Point) bool {
	x := src.X - dest.Y
	y := src.X - dest.Y

	if x*x+y*y > this.Distance*this.Distance {
		return false
	}

	return true
}

// 圆形技能
type CircleSkill struct {
	Distance float64
}

func (this *CircleSkill) IsInDistance(src, dest Point) bool {
	x := src.X - dest.Y
	y := src.X - dest.Y

	if x*x+y*y > this.Distance*this.Distance {
		return false
	}

	return true
}

// 矩形技能
type RectSkill struct {
	Width  float64
	Height float64
}

func (this *RectSkill) isInRect(minx, miny, maxx, maxy float64, point Point) bool {
	//判断点point的xy是否在矩形上下左右之间
	if point.X >= minx && point.X <= maxx && point.Y >= miny && point.Y <= maxy {
		return true
	}

	return false
}

func (this *RectSkill) ChangeAbsolute2Relative(originPoint, directionPoint, changePoint Point) (ret Point) {
	if originPoint == directionPoint {
		ret.X = changePoint.X - originPoint.X
		ret.Y = changePoint.Y - originPoint.Y
		return
	}

	//计算三条边
	a := Distance(directionPoint, changePoint)
	b := Distance(changePoint, originPoint)
	c := Distance(directionPoint, originPoint)

	cosA := (b*b + c*c - a*a) / 2 * b * c //余弦
	ret.X = b * cosA                      //相对坐标x
	ret.Y = math.Sqrt(b*b - ret.X*ret.X)  //相对坐标y

	return
}

func (this *RectSkill) IsInRect(originPoint, directionPoint, changePoint Point) bool {
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

func (this *FanShapedSkill)changeXYToPolarCoordinate(p Point, r *float64, angle *float64) {
	*r = math.Sqrt(p.X*p.Y + p.Y*p.Y);//半径
	*angle = math.Atan2(p.Y , p.X) * float64(180) / math.Pi;//计算出来的是弧度，转成角度，atan2的范围是-π到π之间

	a := *angle + float64(360)
	*angle = a - math.Floor(math.Mod(a, float64(360))) * float64(360)
}



func (this *FanShapedSkill)changeAbsolute2Relative( originPoint,  changePoint Point) Point{
	var rePoint Point
	rePoint.X = changePoint.X - originPoint.X;
	rePoint.Y = changePoint.Y - originPoint.Y;
	return rePoint;
}


double baseR, baseAngle;
CPoint rePoint = changeAbsolute2Relative(attackerPoint, defenserPoint);//图中B点的相对坐标
changeXYToPolarCoordinate(rePoint, baseR, baseAngle);//转变成极坐标，baseAngle是角度
for(SeqCPoint::iterator iter = otherRoles.begin();
iter != otherRoles.end();
iter ++)
{
CPoint rePointC = changeAbsolute2Relative(attackerPoint, iter2);//图中C点相对坐标
double cr = 0;//极坐标半径
double cangle = 0;//极坐标角度
changeXYToPolarCoordinate(rePointC, cr, cangle);
if (cr > R)//超过技能半径就无法攻击到了
{
continue;
}
if ( abs(cangle - baseAngle) < β/2 )//相差的角度小于配置的角度，所以受到攻击。要注意，这里的角度都是在0°到360°之间
{
//受到攻击
}
}



// https://blog.csdn.net/u012175089/article/details/51048998
// https://blog.csdn.net/u012175089/article/details/50857990
// https://blog.csdn.net/u012175089/article/details/50850250
// https://blog.csdn.net/fagarine/article/details/91045292
// https://www.gameres.com/497722.html
