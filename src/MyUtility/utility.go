package MyUtility

import (
	"encoding/json"
	"log"
	"reflect"
	"strconv"
)

func GetAllMapKey(obj map[string]interface{}) []string {
	keyList := make([]string, 0)
	for _, k := range reflect.ValueOf(obj).MapKeys() {
		keyList = append(keyList, k.String())
	}
	return keyList
}

func GetAllMapValue(obj map[string]interface{}) []string {
	valueList := make([]string, 0)

	for _, v := range obj {

		switch reflect.TypeOf(v).Kind() {
		case reflect.Map:
			arr, err := json.Marshal(v)
			if err != nil {
				log.Println("GetAllMapValue err:", err)
			} else {
				valueList = append(valueList, string(arr))
			}
		case reflect.Float64:
			valueList = append(valueList, strconv.FormatFloat(v.(float64), 'f', -1, 64))
		case reflect.String:
			valueList = append(valueList, v.(string))
		}
	}

	return valueList
}

//适用 slice,array,map类型
func ContainInTarget(obj interface{}, target interface{}) bool {

	targetValue := reflect.ValueOf(target)

	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false

}

func GetString2Int(str string) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		log.Println("GetString2Int err:", err)
		return 0
	}
	return value
}

//对传进来的int和string内容求和
func GetSum(list ...interface{}) int {
	sum := 0
	if list != nil {
		for _, v := range list {
			switch v.(type) {
			case int:
				sum += v.(int)
			case string:
				sum += GetString2Int(v.(string))
			}
		}
	}
	return sum
}

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
