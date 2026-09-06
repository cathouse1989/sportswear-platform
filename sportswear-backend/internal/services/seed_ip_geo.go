package services

import (
	"log"

	"gorm.io/gorm"

	"sportswear-backend/internal/models"
)

// SeedIPGeoRanges 初始化 IP 地理库（IP 段 → 国家映射）演示种子数据。
// 仅当表为空时插入（幂等）；用于让「询盘 IP 国家」链路开箱即用。
//
// 注意：以下为覆盖主要贸易国家的「演示/示例」网段（粒度 /24 或 /16，归属基于公开 RIR/ISP 分配常识），
// 并非完整、精确的全球 GeoIP 库；正式环境请通过后台「IP 地理库」接口或 SQL 导入完整数据。
func SeedIPGeoRanges(db *gorm.DB) {
	var count int64
	db.Model(&models.IPGeoRange{}).Count(&count)
	if count > 0 {
		return
	}

	type rng struct{ start, end, cc, name string }
	seeds := []rng{
		{"8.8.8.0", "8.8.8.255", "US", "美国"},          // Google Public DNS
		{"89.36.0.0", "89.36.255.255", "GB", "英国"},     //
		{"88.198.0.0", "88.198.255.255", "DE", "德国"},   // Hetzner
		{"51.75.0.0", "51.75.255.255", "FR", "法国"},     // OVH
		{"82.223.0.0", "82.223.255.255", "ES", "西班牙"}, //
		{"2.228.0.0", "2.228.255.255", "IT", "意大利"},   //
		{"5.79.0.0", "5.79.255.255", "NL", "荷兰"},       // LeaseWeb
		{"1.140.0.0", "1.140.255.255", "AU", "澳大利亚"}, //
		{"142.55.0.0", "142.55.255.255", "CA", "加拿大"}, //
		{"133.242.0.0", "133.242.255.255", "JP", "日本"}, // SAKURA Internet
		{"211.234.0.0", "211.234.255.255", "KR", "韩国"}, //
		{"165.21.0.0", "165.21.255.255", "SG", "新加坡"}, // Singtel
		{"5.195.0.0", "5.195.255.255", "AE", "阿联酋"},   // Etisalat
		{"37.224.0.0", "37.224.255.255", "SA", "沙特阿拉伯"}, // STC
		{"49.207.0.0", "49.207.255.255", "IN", "印度"},   //
		{"177.54.0.0", "177.54.255.255", "BR", "巴西"},   //
		{"189.203.0.0", "189.203.255.255", "MX", "墨西哥"}, //
		{"114.114.114.0", "114.114.114.255", "CN", "中国"}, // 114DNS
		{"1.36.0.0", "1.36.255.255", "HK", "中国香港"},   // PCCW
		{"61.216.0.0", "61.216.255.255", "TW", "中国台湾"}, // HiNet
	}

	imported := 0
	for _, s := range seeds {
		start, ok1 := ipv4ToUint32(s.start)
		end, ok2 := ipv4ToUint32(s.end)
		if !ok1 || !ok2 || start > end {
			continue
		}
		if err := db.Create(&models.IPGeoRange{
			StartIP:     start,
			EndIP:       end,
			Country:     s.cc,
			CountryName: s.name,
		}).Error; err != nil {
			log.Printf("[Seed] 插入 IP 地理库失败: %v", err)
			continue
		}
		imported++
	}
	log.Printf("[Seed] 已初始化 IP 地理库演示数据（%d 条）", imported)
}
