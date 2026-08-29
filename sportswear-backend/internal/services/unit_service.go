package services

import (
	"fmt"
	"strings"
)

// UnitService 单位本地化服务
// 面向全球站点：尺寸对照（US/EU/UK/CM）与重量单位转换（kg/lb）
type UnitService struct{}

// NewUnitService 创建单位服务
func NewUnitService() *UnitService {
	return &UnitService{}
}

// SizeStandard 服装尺寸对照（运动服饰常见尺码）
var sizeStandardMap = map[string][]string{
	"US": {"XS", "S", "M", "L", "XL", "XXL", "3XL"},
	"EU": {"42", "44", "46", "48", "50", "52", "54"},
	"UK": {"34", "36", "38", "40", "42", "44", "46"},
	"CM": {"160/84A", "165/88A", "170/92A", "175/96A", "180/100A", "185/104A", "190/108A"},
}

// GetSizeChart 获取尺寸对照表
func (s *UnitService) GetSizeChart() map[string][]string {
	return sizeStandardMap
}

// GetSizesFor 获取指定标准的尺码列表
func (s *UnitService) GetSizesFor(standard string) []string {
	if sizes, ok := sizeStandardMap[standard]; ok {
		return sizes
	}
	return sizeStandardMap["US"]
}

// ConvertSize 尺码转换（从一个标准到另一个标准）
func (s *UnitService) ConvertSize(size string, from, to string) (string, error) {
	fromSizes, ok1 := sizeStandardMap[from]
	toSizes, ok2 := sizeStandardMap[to]
	if !ok1 || !ok2 {
		return "", fmt.Errorf("不支持的尺码标准: %s -> %s", from, to)
	}

	idx := -1
	for i, s := range fromSizes {
		if strings.EqualFold(s, size) {
			idx = i
			break
		}
	}
	if idx < 0 {
		return "", fmt.Errorf("尺码不存在: %s", size)
	}
	if idx >= len(toSizes) {
		idx = len(toSizes) - 1
	}
	return toSizes[idx], nil
}

// WeightUnit 重量单位
type WeightUnit string

const (
	WeightUnitKG WeightUnit = "kg"
	WeightUnitLB WeightUnit = "lb"
)

// ConvertWeight 重量单位转换
func (s *UnitService) ConvertWeight(value float64, from, to WeightUnit) (float64, error) {
	if from == to {
		return value, nil
	}
	switch {
	case from == WeightUnitKG && to == WeightUnitLB:
		return value * 2.20462, nil
	case from == WeightUnitLB && to == WeightUnitKG:
		return value / 2.20462, nil
	default:
		return 0, fmt.Errorf("不支持的重量单位转换: %s -> %s", from, to)
	}
}

// FormatWeight 格式化重量
func (s *UnitService) FormatWeight(value float64, unit WeightUnit) string {
	return fmt.Sprintf("%.1f %s", value, unit)
}

// LengthUnit 长度单位
type LengthUnit string

const (
	LengthUnitCM LengthUnit = "cm"
	LengthUnitIN LengthUnit = "in"
)

// ConvertLength 长度单位转换
func (s *UnitService) ConvertLength(value float64, from, to LengthUnit) (float64, error) {
	if from == to {
		return value, nil
	}
	switch {
	case from == LengthUnitCM && to == LengthUnitIN:
		return value / 2.54, nil
	case from == LengthUnitIN && to == LengthUnitCM:
		return value * 2.54, nil
	default:
		return 0, fmt.Errorf("不支持的长度单位转换: %s -> %s", from, to)
	}
}
