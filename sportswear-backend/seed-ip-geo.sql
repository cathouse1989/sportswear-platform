-- ============================================================
--  IP 地理库（ip_geo_ranges）建表 + 演示种子数据
--  表结构对齐 GORM 模型 models.IPGeoRange（AutoMigrate 也会生成同样结构）
--  用途：询盘按 IP 解析国家（写入 leads.ip_country）
--  说明：以下为覆盖主要贸易国家的「演示/示例」网段，正式环境请导入完整 GeoIP 库
-- ============================================================

CREATE TABLE IF NOT EXISTS ip_geo_ranges (
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now(),
    deleted_at   timestamptz,
    start_ip     bigint NOT NULL,
    end_ip       bigint NOT NULL,
    country      varchar(10) NOT NULL,
    country_name varchar(100)
);

CREATE INDEX IF NOT EXISTS idx_ip_geo_start              ON ip_geo_ranges (start_ip);
CREATE INDEX IF NOT EXISTS idx_ip_geo_ranges_country     ON ip_geo_ranges (country);
CREATE INDEX IF NOT EXISTS idx_ip_geo_ranges_deleted_at  ON ip_geo_ranges (deleted_at);

-- 仅当表为空时插入（幂等，避免重复执行时叠加数据）
INSERT INTO ip_geo_ranges (start_ip, end_ip, country, country_name)
SELECT v.*
FROM (VALUES
    ('8.8.8.0'::inet - '0.0.0.0'::inet,       '8.8.8.255'::inet - '0.0.0.0'::inet,      'US', '美国'),
    ('89.36.0.0'::inet - '0.0.0.0'::inet,     '89.36.255.255'::inet - '0.0.0.0'::inet,   'GB', '英国'),
    ('88.198.0.0'::inet - '0.0.0.0'::inet,    '88.198.255.255'::inet - '0.0.0.0'::inet,  'DE', '德国'),
    ('51.75.0.0'::inet - '0.0.0.0'::inet,     '51.75.255.255'::inet - '0.0.0.0'::inet,   'FR', '法国'),
    ('82.223.0.0'::inet - '0.0.0.0'::inet,    '82.223.255.255'::inet - '0.0.0.0'::inet,  'ES', '西班牙'),
    ('2.228.0.0'::inet - '0.0.0.0'::inet,     '2.228.255.255'::inet - '0.0.0.0'::inet,   'IT', '意大利'),
    ('5.79.0.0'::inet - '0.0.0.0'::inet,      '5.79.255.255'::inet - '0.0.0.0'::inet,    'NL', '荷兰'),
    ('1.140.0.0'::inet - '0.0.0.0'::inet,     '1.140.255.255'::inet - '0.0.0.0'::inet,   'AU', '澳大利亚'),
    ('142.55.0.0'::inet - '0.0.0.0'::inet,    '142.55.255.255'::inet - '0.0.0.0'::inet,  'CA', '加拿大'),
    ('133.242.0.0'::inet - '0.0.0.0'::inet,   '133.242.255.255'::inet - '0.0.0.0'::inet, 'JP', '日本'),
    ('211.234.0.0'::inet - '0.0.0.0'::inet,   '211.234.255.255'::inet - '0.0.0.0'::inet, 'KR', '韩国'),
    ('165.21.0.0'::inet - '0.0.0.0'::inet,    '165.21.255.255'::inet - '0.0.0.0'::inet,  'SG', '新加坡'),
    ('5.195.0.0'::inet - '0.0.0.0'::inet,     '5.195.255.255'::inet - '0.0.0.0'::inet,   'AE', '阿联酋'),
    ('37.224.0.0'::inet - '0.0.0.0'::inet,    '37.224.255.255'::inet - '0.0.0.0'::inet,  'SA', '沙特阿拉伯'),
    ('49.207.0.0'::inet - '0.0.0.0'::inet,    '49.207.255.255'::inet - '0.0.0.0'::inet,  'IN', '印度'),
    ('177.54.0.0'::inet - '0.0.0.0'::inet,    '177.54.255.255'::inet - '0.0.0.0'::inet,  'BR', '巴西'),
    ('189.203.0.0'::inet - '0.0.0.0'::inet,   '189.203.255.255'::inet - '0.0.0.0'::inet, 'MX', '墨西哥'),
    ('114.114.114.0'::inet - '0.0.0.0'::inet, '114.114.114.255'::inet - '0.0.0.0'::inet, 'CN', '中国'),
    ('1.36.0.0'::inet - '0.0.0.0'::inet,      '1.36.255.255'::inet - '0.0.0.0'::inet,    'HK', '中国香港'),
    ('61.216.0.0'::inet - '0.0.0.0'::inet,    '61.216.255.255'::inet - '0.0.0.0'::inet,  'TW', '中国台湾')
) AS v(start_ip, end_ip, country, country_name)
WHERE NOT EXISTS (SELECT 1 FROM ip_geo_ranges LIMIT 1);
