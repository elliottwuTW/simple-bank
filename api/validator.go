package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/simple_bank/util"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	// 呼叫 Field() 取得欄位的值，使用 Interface() 把 reflection 的值轉換成 interface{}
	// 再嘗試進行 string conversion
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}
