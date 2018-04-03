package util

import (
	"background/common/logger"
	"fmt"

	"github.com/shopspring/decimal"
)

func Equal(price1, price2 interface{}) bool {
	var decimalPrice1, decimalPrice2 decimal.Decimal
	var err error
	decimalPrice1, err = getDecimal(price1)
	if err != nil {
		logger.Error(fmt.Sprintf("[%s] : [%s]", price1, err))
		return false
	}

	decimalPrice2, err = getDecimal(price2)
	if err != nil {
		logger.Error(fmt.Sprintf("[%s] : [%s]", price2, err))
		return false
	}
	return decimalPrice1.Equals(decimalPrice2)
}

func Add(args ...interface{}) (string, error) {
	var decimalVals []decimal.Decimal
	for _, val := range args {
		decimalVal, err := getDecimal(val)
		if err != nil {
			logger.Error(fmt.Sprintf("[%s] : [%s]", val, err))
			return "", err
		}
		decimalVals = append(decimalVals, decimalVal)
	}

	var result decimal.Decimal
	for _, decimalVal := range decimalVals {
		result = result.Add(decimalVal)
	}
	return result.String(), nil
}

func Sub(args ...interface{}) (string, error) {
	var decimalVals []decimal.Decimal
	for _, val := range args {
		decimalVal, err := getDecimal(val)
		if err != nil {
			logger.Error(fmt.Sprintf("[%s] : [%s]", val, err))
			return "", err
		}
		decimalVals = append(decimalVals, decimalVal)
	}

	var result decimal.Decimal
	result = decimalVals[0]
	for idx, decimalVal := range decimalVals {
		if idx != 0 {
			result = result.Sub(decimalVal)
		}
	}
	return result.String(), nil
}

func Mul(args ...interface{}) (string, error) {
	var decimalVals []decimal.Decimal
	for _, val := range args {
		decimalVal, err := getDecimal(val)
		if err != nil {
			logger.Error(fmt.Sprintf("[%s] : [%s]", val, err))
			return "", err
		}
		decimalVals = append(decimalVals, decimalVal)
	}

	var result decimal.Decimal
	result = decimalVals[0]
	for i := 1; i < len(decimalVals); i++ {
		result = result.Mul(decimalVals[i])
	}
	return result.String(), nil
}

func Div(args ...interface{}) (string, error) {
	var decimalVals []decimal.Decimal
	for _, val := range args {
		decimalVal, err := getDecimal(val)
		if err != nil {
			logger.Error(fmt.Sprintf("[%s] : [%s]", val, err))
			return "", err
		}
		decimalVals = append(decimalVals, decimalVal)
	}

	var result decimal.Decimal
	result = decimalVals[0]
	for i := 1; i < len(decimalVals); i++ {
		result = result.Div(decimalVals[i])
	}
	return result.String(), nil
}

func getDecimal(val interface{}) (decimal.Decimal, error) {
	str := fmt.Sprint(val)
	if str == "" {
		return decimal.Zero, nil
	}
	return decimal.NewFromString(str)
}
