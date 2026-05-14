package xjson

import (
	"fmt"
	"strconv"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
)

var StrictMode = false // 🔥 严格模式：解析失败直接报错

func init() {
	registerInt("int", func(v int64) any { return int(v) })
	registerInt("int8", func(v int64) any { return int8(v) })
	registerInt("int16", func(v int64) any { return int16(v) })
	registerInt("int32", func(v int64) any { return int32(v) })
	registerInt("int64", func(v int64) any { return v })

	registerFloat("float32", func(v float64) any { return float32(v) })
	registerFloat("float64", func(v float64) any { return v })

	registerBool()
	registerDecimal()
	registerTime()
}

//
// ======================== INT ========================
//

func registerInt(typeName string, cast func(int64) any) {
	jsoniter.RegisterTypeDecoderFunc(typeName, func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		switch iter.WhatIsNext() {

		case jsoniter.NumberValue:
			set(ptr, cast(iter.ReadInt64()))

		case jsoniter.StringValue:
			str := iter.ReadString()
			val, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				handleError(iter, typeName, "invalid int: "+str)
				return
			}
			set(ptr, cast(val))

		case jsoniter.BoolValue:
			if iter.ReadBool() {
				set(ptr, cast(1))
			} else {
				set(ptr, cast(0))
			}

		default:
			iter.Skip()
		}
	})
}

//
// ======================== FLOAT ========================
//

func registerFloat(typeName string, cast func(float64) any) {
	jsoniter.RegisterTypeDecoderFunc(typeName, func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		switch iter.WhatIsNext() {

		case jsoniter.NumberValue:
			set(ptr, cast(iter.ReadFloat64()))

		case jsoniter.StringValue:
			str := iter.ReadString()
			val, err := strconv.ParseFloat(str, 64)
			if err != nil {
				handleError(iter, typeName, "invalid float: "+str)
				return
			}
			set(ptr, cast(val))

		case jsoniter.BoolValue:
			if iter.ReadBool() {
				set(ptr, cast(1))
			} else {
				set(ptr, cast(0))
			}

		default:
			iter.Skip()
		}
	})
}

//
// ======================== BOOL ========================
//

func registerBool() {
	jsoniter.RegisterTypeDecoderFunc("bool", func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		switch iter.WhatIsNext() {

		case jsoniter.BoolValue:
			*(*bool)(ptr) = iter.ReadBool()

		case jsoniter.NumberValue:
			*(*bool)(ptr) = iter.ReadInt() != 0

		case jsoniter.StringValue:
			str := iter.ReadString()

			switch str {
			case "true", "1":
				*(*bool)(ptr) = true
			case "false", "0":
				*(*bool)(ptr) = false
			default:
				handleError(iter, "bool", "invalid bool: "+str)
			}

		default:
			iter.Skip()
		}
	})
}

//
// ======================== DECIMAL ========================
//

func registerDecimal() {
	jsoniter.RegisterTypeDecoderFunc("decimal.Decimal", func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		switch iter.WhatIsNext() {

		case jsoniter.NumberValue:
			str := iter.ReadNumber().String()
			d, err := decimal.NewFromString(str)
			if err != nil {
				handleError(iter, "decimal", "invalid decimal: "+str)
				return
			}
			*(*decimal.Decimal)(ptr) = d

		case jsoniter.StringValue:
			str := iter.ReadString()
			d, err := decimal.NewFromString(str)
			if err != nil {
				handleError(iter, "decimal", "invalid decimal: "+str)
				return
			}
			*(*decimal.Decimal)(ptr) = d

		case jsoniter.BoolValue:
			if iter.ReadBool() {
				*(*decimal.Decimal)(ptr) = decimal.NewFromInt(1)
			} else {
				*(*decimal.Decimal)(ptr) = decimal.Zero
			}

		default:
			iter.Skip()
		}
	})
}

//
// ======================== TIME ========================
//

func registerTime() {
	jsoniter.RegisterTypeDecoderFunc("time.Time", func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		switch iter.WhatIsNext() {

		case jsoniter.NumberValue:
			num := iter.ReadInt64()
			if num > 1e12 {
				*(*time.Time)(ptr) = time.UnixMilli(num).UTC()
			} else {
				*(*time.Time)(ptr) = time.Unix(num, 0).UTC()
			}

		case jsoniter.StringValue:
			str := iter.ReadString()

			// ✅ 你要求：传了就更新 → "" 也要写
			if str == "" {
				*(*time.Time)(ptr) = time.Time{}
				return
			}

			t, err := parseTime(str)
			if err != nil {
				handleError(iter, "time", "invalid time: "+str)
				return
			}

			*(*time.Time)(ptr) = t

		default:
			iter.Skip()
		}
	})
}

func parseTime(str string) (time.Time, error) {
	// 时间戳字符串
	if ts, err := strconv.ParseInt(str, 10, 64); err == nil {
		if ts > 1e12 {
			return time.UnixMilli(ts).UTC(), nil
		}
		return time.Unix(ts, 0).UTC(), nil
	}

	// RFC3339
	if t, err := time.Parse(time.RFC3339, str); err == nil {
		return t.UTC(), nil
	}

	// 常见格式
	loc := time.Local

	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006/01/02 15:04:05",
		"2006/01/02",
	}

	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, str, loc); err == nil {
			return t.UTC(), nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid time: %s", str)
}

//
// ======================== COMMON ========================
//

func handleError(iter *jsoniter.Iterator, typ, msg string) {
	if StrictMode {
		iter.ReportError(typ, msg)
	}
	// ❗非严格模式：直接跳过（不写）
}

func set(ptr unsafe.Pointer, val any) {
	switch v := val.(type) {
	case int:
		*(*int)(ptr) = v
	case int8:
		*(*int8)(ptr) = v
	case int16:
		*(*int16)(ptr) = v
	case int32:
		*(*int32)(ptr) = v
	case int64:
		*(*int64)(ptr) = v
	case float32:
		*(*float32)(ptr) = v
	case float64:
		*(*float64)(ptr) = v
	}
}
