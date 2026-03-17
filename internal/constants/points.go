package constants

type PointsSourceType int8
type PointsChangeType int8

const (
	PointsSourceConsume  PointsSourceType = 10 // 消费
	PointsSourceRecharge PointsSourceType = 20 // 充值
	PointsSourceActivity PointsSourceType = 30 // 活动
	PointsSourceAdmin    PointsSourceType = 40 // 管理员修改

	PointsChangeTypeAdd    PointsChangeType = 1 // 增加积分
	PointsChangeTypeReduce PointsChangeType = 2 // 减少积分
)

// String 返回类型标题
func (t PointsSourceType) String() string {
	switch t {
	case PointsSourceConsume:
		return "消费"
	case PointsSourceRecharge:
		return "充值"
	case PointsSourceActivity:
		return "活动"
	case PointsSourceAdmin:
		return "管理员修改"
	default:
		return "未知"
	}
}

// IsValid 检查类型是否有效
func (t PointsSourceType) IsValid() bool {
	switch t {
	case PointsSourceConsume, PointsSourceRecharge,
		PointsSourceActivity, PointsSourceAdmin:
		return true
	default:
		return false
	}
}

// GetPointsSourceTypes 获取所有类型映射关系
func GetPointsSourceTypes() map[int8]string {
	return map[int8]string{
		int8(PointsSourceConsume):  PointsSourceConsume.String(),
		int8(PointsSourceRecharge): PointsSourceRecharge.String(),
		int8(PointsSourceActivity): PointsSourceActivity.String(),
		int8(PointsSourceAdmin):    PointsSourceAdmin.String(),
	}
}
