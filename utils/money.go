// Package utils 提供金额精度计算工具。
// 基于 github.com/shopspring/decimal，避免浮点精度问题。
// Money2 系列保留 2 位小数（人民币等货币场景）
// Money6 系列保留 6 位小数（汇率、单价等高精度场景）
package utils

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// ------------------- 2位小数 -------------------

// Money2Add 加法，保留2位小数，四舍五入
func Money2Add(a, b string) (string, error) {
	return moneyOp(a, b, 2, func(x, y decimal.Decimal) decimal.Decimal {
		return x.Add(y)
	})
}

// Money2Sub 减法，保留2位小数，四舍五入
func Money2Sub(a, b string) (string, error) {
	return moneyOp(a, b, 2, func(x, y decimal.Decimal) decimal.Decimal {
		return x.Sub(y)
	})
}

// Money2Mul 乘法，保留2位小数，四舍五入
func Money2Mul(a, b string) (string, error) {
	return moneyOp(a, b, 2, func(x, y decimal.Decimal) decimal.Decimal {
		return x.Mul(y)
	})
}

// Money2Div 除法，保留2位小数，四舍五入；除数为0时返回错误
func Money2Div(a, b string) (string, error) {
	return moneyDiv(a, b, 2)
}

// ------------------- 6位小数 -------------------

// Money6Add 加法，保留6位小数，四舍五入
func Money6Add(a, b string) (string, error) {
	return moneyOp(a, b, 6, func(x, y decimal.Decimal) decimal.Decimal {
		return x.Add(y)
	})
}

// Money6Sub 减法，保留6位小数，四舍五入
func Money6Sub(a, b string) (string, error) {
	return moneyOp(a, b, 6, func(x, y decimal.Decimal) decimal.Decimal {
		return x.Sub(y)
	})
}

// Money6Mul 乘法，保留6位小数，四舍五入
func Money6Mul(a, b string) (string, error) {
	return moneyOp(a, b, 6, func(x, y decimal.Decimal) decimal.Decimal {
		return x.Mul(y)
	})
}

// Money6Div 除法，保留6位小数，四舍五入；除数为0时返回错误
func Money6Div(a, b string) (string, error) {
	return moneyDiv(a, b, 6)
}

// ------------------- 比较 -------------------

// MoneyCmp 比较两个金额字符串，返回 -1/0/1
func MoneyCmp(a, b string) (int, error) {
	x, err := decimal.NewFromString(a)
	if err != nil {
		return 0, err
	}
	y, err := decimal.NewFromString(b)
	if err != nil {
		return 0, err
	}
	return x.Cmp(y), nil
}

// MoneyIsZero 判断金额是否为0
func MoneyIsZero(a string) (bool, error) {
	x, err := decimal.NewFromString(a)
	if err != nil {
		return false, err
	}
	return x.IsZero(), nil
}

// MoneyIsNegative 判断金额是否为负数
func MoneyIsNegative(a string) (bool, error) {
	x, err := decimal.NewFromString(a)
	if err != nil {
		return false, err
	}
	return x.IsNegative(), nil
}

// ------------------- 格式化 -------------------

// MoneyFormat2 格式化为2位小数字符串（不做运算，仅格式化）
func MoneyFormat2(a string) (string, error) {
	x, err := decimal.NewFromString(a)
	if err != nil {
		return "", err
	}
	return x.StringFixed(2), nil
}

// MoneyFormat6 格式化为6位小数字符串
func MoneyFormat6(a string) (string, error) {
	x, err := decimal.NewFromString(a)
	if err != nil {
		return "", err
	}
	return x.StringFixed(6), nil
}

// ------------------- 内部工具 -------------------

func moneyOp(a, b string, places int32, op func(decimal.Decimal, decimal.Decimal) decimal.Decimal) (string, error) {
	x, err := decimal.NewFromString(a)
	if err != nil {
		return "", err
	}
	y, err := decimal.NewFromString(b)
	if err != nil {
		return "", err
	}
	return op(x, y).StringFixed(places), nil
}

func moneyDiv(a, b string, places int32) (string, error) {
	x, err := decimal.NewFromString(a)
	if err != nil {
		return "", err
	}
	y, err := decimal.NewFromString(b)
	if err != nil {
		return "", err
	}
	if y.IsZero() {
		return "", fmt.Errorf("division by zero")
	}
	return x.Div(y).StringFixed(places), nil
}
