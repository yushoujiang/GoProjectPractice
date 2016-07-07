package MyUtility

import (
// "log"
// "reflect"
)

func GetMapWithDefault(mapValue map[string]interface{}, key string, defaultValue interface{}) interface{} {

	ret, ok := mapValue[key]
	if ok {
		return ret
	} else {
		return defaultValue
	}
}

func GetMapRetInt(mapValue map[string]interface{}, key string, defaultValue interface{}) int {

	tempValue := GetMapWithDefault(mapValue, key, defaultValue)
	// log.Println("type=", reflect.TypeOf(tempValue))

	ret, ok := tempValue.(float64)
	if ok {
		var value int = int(ret)
		return value
	} else {
		newRet, success := defaultValue.(int)
		if success {
			return newRet
		} else {
			return 0
		}
	}
}

func GetMapRetString(mapValue map[string]interface{}, key string, defaultValue interface{}) string {

	tempValue := GetMapWithDefault(mapValue, key, defaultValue)
	ret, ok := tempValue.(string)
	if ok {
		return ret
	} else {
		ret, ok = defaultValue.(string)
		if ok {
			return ret
		} else {
			return "get error"
		}
	}
}
