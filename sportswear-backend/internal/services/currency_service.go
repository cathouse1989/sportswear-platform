package services

import (
	"fmt"
	"strings"

	"gorm.io/gorm"

	"sportswear-backend/internal/models"
)

// CurrencyService 多货币服务
// 面向全球站点，提供货币列表、汇率转换、价格本地化格式化
type CurrencyService struct {
	db *gorm.DB
}

// NewCurrencyService 创建货币服务
func NewCurrencyService(db *gorm.DB) *CurrencyService {
	return &CurrencyService{db: db}
}

// ListCurrencies 货币列表（按排序）
func (s *CurrencyService) ListCurrencies(onlyActive bool) ([]models.Currency, error) {
	var currencies []models.Currency
	query := s.db.Model(&models.Currency{})
	if onlyActive {
		query = query.Where("is_active = ?", true)
	}
	err := query.Order("sort_order ASC").Find(&currencies).Error
	return currencies, err
}

// GetCurrency 获取指定货币
func (s *CurrencyService) GetCurrency(code string) (*models.Currency, error) {
	var currency models.Currency
	err := s.db.Where("code = ? AND is_active = ?", code, true).First(&currency).Error
	if err != nil {
		return nil, fmt.Errorf("货币不存在: %s", code)
	}
	return &currency, nil
}

// GetDefaultCurrency 获取默认货币
func (s *CurrencyService) GetDefaultCurrency() (*models.Currency, error) {
	var currency models.Currency
	err := s.db.Where("is_default = ? AND is_active = ?", true, true).First(&currency).Error
	if err != nil {
		// 回退到基准货币
		err = s.db.Where("is_base = ? AND is_active = ?", true, true).First(&currency).Error
	}
	if err != nil {
		return nil, fmt.Errorf("无可用货币")
	}
	return &currency, nil
}

// Convert 金额从基准货币（USD）转换为目标货币
// amount 是基准货币（USD）金额，返回目标货币金额
func (s *CurrencyService) Convert(amountUSD float64, targetCode string) (float64, error) {
	target, err := s.GetCurrency(targetCode)
	if err != nil {
		return amountUSD, err
	}
	return amountUSD * target.Rate, nil
}

// FormatPrice 本地化价格格式化
// 示例：FormatPrice(1234.5, "USD") -> "$1,234.50"
//
//	FormatPrice(1234.5, "EUR") -> "€1.234,50"
func (s *CurrencyService) FormatPrice(amount float64, currencyCode string) (string, error) {
	currency, err := s.GetCurrency(currencyCode)
	if err != nil {
		return "", err
	}
	return FormatPriceWithCurrency(amount, currency), nil
}

// FormatPriceWithCurrency 根据货币配置格式化价格
func FormatPriceWithCurrency(amount float64, c *models.Currency) string {
	// 处理小数位数
	format := fmt.Sprintf("%%.%df", c.DecimalPlaces)
	numStr := fmt.Sprintf(format, amount)

	// 分离整数与小数部分
	parts := strings.Split(numStr, ".")
	intPart := parts[0]
	decPart := ""
	if len(parts) > 1 {
		decPart = parts[1]
	}

	// 添加千分位分隔符
	if c.ThousandsSeparator != "" && len(intPart) > 3 {
		var builder strings.Builder
		n := len(intPart)
		for i := 0; i < n; i++ {
			if i > 0 && (n-i)%3 == 0 {
				builder.WriteString(c.ThousandsSeparator)
			}
			builder.WriteByte(intPart[i])
		}
		intPart = builder.String()
	}

	// 组装数字
	formatted := intPart
	if decPart != "" {
		formatted += c.DecimalSeparator + decPart
	}

	// 添加货币符号
	if c.Symbol == "" {
		formatted += " " + c.Code
	} else if c.SymbolPosition == "suffix" {
		formatted += c.Symbol
	} else {
		formatted = c.Symbol + formatted
	}

	return formatted
}

// InitDefaults 初始化默认货币
func (s *CurrencyService) InitDefaults() error {
	var count int64
	s.db.Model(&models.Currency{}).Count(&count)
	if count > 0 {
		return nil
	}
	return s.db.Create(&models.DefaultCurrencies).Error
}
