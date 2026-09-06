package services

import (
	"errors"
	"net"
	"strings"

	"gorm.io/gorm"

	"sportswear-backend/internal/models"
)

// GeoIPService IP 地理定位服务
// 基于数据库表 ip_geo_ranges（IP 段 → 国家映射）判断访客所在国家，
// 数据由后台「IP 地理库」接口导入维护，无需引入外部 .mmdb 文件或第三方 HTTP 服务。
type GeoIPService struct {
	db *gorm.DB
}

// NewGeoIPService 创建 IP 地理定位服务
func NewGeoIPService(db *gorm.DB) *GeoIPService {
	return &GeoIPService{db: db}
}

// Enabled 判断 IP 地理库是否有数据
func (s *GeoIPService) Enabled() bool {
	var count int64
	if err := s.db.Model(&models.IPGeoRange{}).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

// LookupCountry 根据 IP 返回国家 ISO2 代码与国家名。
// 仅支持 IPv4；内网/保留地址、非法地址、库中未命中时返回空字符串。
func (s *GeoIPService) LookupCountry(ip string) (iso2, name string) {
	num, ok := ipv4ToUint32(ip)
	if !ok {
		return "", ""
	}

	var rec models.IPGeoRange
	// 取 start_ip <= num 的最大 start_ip，再校验 end_ip >= num（命中 start_ip 索引）
	if err := s.db.Where("start_ip <= ?", num).Order("start_ip DESC").First(&rec).Error; err != nil {
		return "", ""
	}
	if num > rec.EndIP {
		return "", ""
	}
	return rec.Country, rec.CountryName
}

// ListIPGeoRanges 分页查询 IP 地理库
func (s *GeoIPService) ListIPGeoRanges(page, pageSize int) ([]models.IPGeoRange, int64, error) {
	var items []models.IPGeoRange
	var total int64
	q := s.db.Model(&models.IPGeoRange{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if err := q.Order("start_ip ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// ImportIPGeoRanges 批量导入/更新 IP 段数据（按 start_ip 幂等 upsert）
func (s *GeoIPService) ImportIPGeoRanges(items []IPGeoRangeItem) (int, error) {
	if len(items) == 0 {
		return 0, errors.New("导入数据为空")
	}

	imported := 0
	for _, it := range items {
		if it.StartIP > it.EndIP {
			continue
		}
		var existing models.IPGeoRange
		err := s.db.Where("start_ip = ?", it.StartIP).First(&existing).Error
		if err == nil {
			existing.EndIP = it.EndIP
			existing.Country = it.Country
			existing.CountryName = it.CountryName
			if err := s.db.Save(&existing).Error; err != nil {
				return imported, err
			}
		} else {
			rec := models.IPGeoRange{
				StartIP:     it.StartIP,
				EndIP:       it.EndIP,
				Country:     it.Country,
				CountryName: it.CountryName,
			}
			if err := s.db.Create(&rec).Error; err != nil {
				return imported, err
			}
		}
		imported++
	}
	return imported, nil
}

// DeleteIPGeoRange 删除单条 IP 段
func (s *GeoIPService) DeleteIPGeoRange(id string) error {
	return s.db.Delete(&models.IPGeoRange{}, "id = ?", id).Error
}

// IPGeoRangeItem IP 段导入项
type IPGeoRangeItem struct {
	StartIP     int64  `json:"start_ip" binding:"required"`
	EndIP       int64  `json:"end_ip" binding:"required"`
	Country     string `json:"country" binding:"required"`
	CountryName string `json:"country_name"`
}

// ipv4ToUint32 把 IPv4 地址转为 32 位整数（用于范围查询）。
// 内网/保留地址、IPv6、非法地址返回 ok=false。
func ipv4ToUint32(ip string) (int64, bool) {
	if ip == "" || ip == "unknown" || ip == "anonymous" {
		return 0, false
	}
	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed == nil {
		return 0, false
	}
	v4 := parsed.To4()
	if v4 == nil {
		return 0, false
	}
	// 内网 / 保留地址不做查询，避免无意义数据
	if parsed.IsPrivate() || parsed.IsLoopback() || parsed.IsLinkLocalUnicast() ||
		parsed.IsLinkLocalMulticast() || parsed.IsUnspecified() || parsed.IsMulticast() {
		return 0, false
	}
	return int64(v4[0])<<24 | int64(v4[1])<<16 | int64(v4[2])<<8 | int64(v4[3]), true
}
