package ruleenginelib

import (
	"fmt"
)

func EvaluateOperator(identifier, value interface{}, operator string) (bool, error) {
	switch operator {
	case "anyof":
		switch valSlice := value.(type) {
		case []interface{}:
			factNum, err := assertIsNumber(identifier)
			if err == nil {
				for _, val := range valSlice {
					valueNum, err := assertIsNumber(val)
					if err == nil && valueNum == factNum {
						return true, nil
					}
				}
			} else {
				for _, val := range valSlice {
					if identifier == val {
						return true, nil
					}
				}
			}
			return false, nil
		default:
			factNum, err := assertIsNumber(identifier)
			if err == nil {
				valueNum, err := assertIsNumber(value)
				if err != nil {
					return false, err
				}
				return factNum == valueNum, nil
			}
			return identifier == value, nil
		}
	case "noneof":
		switch valSlice := value.(type) {
		case []interface{}:
			factNum, err := assertIsNumber(identifier)
			if err == nil {
				for _, val := range valSlice {
					valueNum, err := assertIsNumber(val)
					if err == nil && valueNum == factNum {
						return false, nil
					}
				}
			} else {
				for _, val := range valSlice {
					if identifier == val {
						return false, nil
					}
				}
			}
			return true, nil
		default:
			factNum, err := assertIsNumber(identifier)
			if err == nil {
				valueNum, err := assertIsNumber(value)
				if err != nil {
					return false, err
				}
				return factNum != valueNum, nil
			}
			return identifier != value, nil
		}
	case "=":
		fallthrough
	case "eq":
		factNum, err := assertIsNumber(identifier)
		if err == nil {
			valueNum, err := assertIsNumber(value)
			if err != nil {
				return false, err
			}
			return factNum == valueNum, nil
		}

		return identifier == value, nil
	case "!=":
		fallthrough
	case "neq":
		factNum, err := assertIsNumber(identifier)
		if err == nil {
			valueNum, err := assertIsNumber(value)
			if err != nil {
				return false, err
			}
			return factNum != valueNum, nil
		}

		return identifier != value, nil

	case "<":
		fallthrough
	case "lt":
		factNum, err := assertIsNumber(identifier)
		if err != nil {
			return false, err
		}
		valueNum, err := assertIsNumber(value)
		if err != nil {
			return false, err
		}

		return factNum < valueNum, nil

	case ">":
		fallthrough
	case "gt":
		factNum, err := assertIsNumber(identifier)
		if err != nil {
			return false, err
		}
		valueNum, err := assertIsNumber(value)
		if err != nil {
			return false, err
		}

		return factNum > valueNum, nil

	case ">=":
		fallthrough
	case "gte":
		factNum, err := assertIsNumber(identifier)
		if err != nil {
			return false, err
		}
		valueNum, err := assertIsNumber(value)
		if err != nil {
			return false, err
		}

		return factNum >= valueNum, nil

	case "<=":
		fallthrough
	case "lte":
		factNum, err := assertIsNumber(identifier)
		if err != nil {
			return false, err
		}
		valueNum, err := assertIsNumber(value)
		if err != nil {
			return false, err
		}

		return factNum <= valueNum, nil

	default:
		return false, fmt.Errorf("unrecognised operator %s", operator)
	}
	return false, fmt.Errorf("unrecognised condition kind %T", value)
}

func assertIsNumber(v interface{}) (float64, error) {
	isFloat := true
	var d int
	var f float64

	d, ok := v.(int)
	if !ok {
		f, ok = v.(float64)
		if !ok {
			return 0, fmt.Errorf("%s is not a number", v)
		}
	} else {
		isFloat = false
	}

	if isFloat {
		return f, nil
	}
	return float64(d), nil
}
