package model

// Region 定义模型
type Region struct {
	ID    uint32 `json:"id" gorm:"primaryKey"`
	Name  string `json:"name" gorm:"size:255;not null;default:'';comment:区划名称"`
	Pid   uint32 `json:"pid" gorm:"not null;default:0;comment:父级ID"`
	Code  string `json:"code" gorm:"size:255;not null;default:'';comment:区划编码"`
	Level int8   `json:"level" gorm:"type:tinyint;not null;default:1;comment:层级 ( 1省级 2市级 3区/县级 ) "`

	Children []*Region `json:"children,omitempty" gorm:"-"`
}

// TableName 指定数据表名称
func (m *Region) TableName() string {
	return TableNamePrefix + "region"
}

// RegionList 地区列表
type RegionList []Region

// Tree 将地区列表转换为树形结构
func (list RegionList) Tree() []*Region {
	regionMap := make(map[uint32]*Region)
	for i := range list {
		regionMap[list[i].ID] = &list[i]
	}

	var rootRegions []*Region
	for i := range list {
		if list[i].Pid == 0 {
			rootRegions = append(rootRegions, &list[i])
			continue
		}

		if parent, ok := regionMap[list[i].Pid]; ok {
			parent.Children = append(parent.Children, &list[i])
		}
	}

	return rootRegions
}
