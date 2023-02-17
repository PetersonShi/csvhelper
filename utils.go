package csvhelper

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

const TAG = "csv-field"

func StructBind(obj interface{}, data map[string]interface{}) error {
	objValue := valueElem(obj)

	fieldMap := make(map[string]reflect.Value)
	for i := 0; i < objValue.NumField(); i++ {
		fieldInfo := objValue.Type().Field(i)
		tag := fieldInfo.Tag.Get(TAG)
		if tag == "" {
			tag = fieldInfo.Name
		}
		if objValue.Field(i).CanSet() {
			fieldMap[strings.ToLower(tag)] = objValue.Field(i)
		}
	}

	for key, value := range data {
		key = strings.ToLower(key)
		field, ok := fieldMap[key]
		if ok == false {
			continue
		}
		if field.IsValid() == false {
			continue
		}
		set(field, ToString(value))
	}
	return nil
}

func set(field reflect.Value, data string) {
	switch field.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(data, 10, 64)
		if err == nil {
			field.SetUint(v)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(data, 10, 64)
		if err == nil {
			field.SetInt(v)
		}
		break
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(data, 64)
		if err == nil {
			field.SetFloat(v)
		}
		break
	case reflect.String:
		field.SetString(data)
		break
	case reflect.Bool:
		v, err := strconv.ParseBool(data)
		if err == nil {
			field.SetBool(v)
		}
	}
}

func ReorganizeKeyField(keyField string) string {
	keyField = strings.TrimSpace(keyField)
	keyNameData := []byte(keyField)
	keyField = strings.ToUpper(string(keyNameData[0])) + string(keyNameData[1:])
	return keyField
}

func valueElem(data interface{}) reflect.Value {
	valueData := reflect.ValueOf(data)
	if valueData.Kind() != reflect.Ptr {
		panic("obj must be ptr")
	}

	for {
		if valueData.Kind() == reflect.Ptr {
			valueData = valueData.Elem()
		} else {
			break
		}
	}
	return valueData
}

func makeCSVTitle(data interface{}) string {
	retStr := ""
	valueData := valueElem(data)

	for i := 0; i < valueData.NumField(); i++ {
		retStr += strings.ToLower(valueData.Type().Field(i).Name) + ","
	}
	retStr = strings.TrimRight(retStr, ",")
	retStr += "\n"
	return retStr
}

func makeCSVContent(data interface{}) string {
	retStr := ""
	valueData := valueElem(data)

	for i := 0; i < valueData.NumField(); i++ {
		retStr += ToString(valueData.Field(i).Interface()) + ","
	}

	retStr = strings.TrimRight(retStr, ",")
	retStr += "\n"
	return retStr
}

func ToString(value interface{}) string {
	switch value.(type) {
	case int8:
		return strconv.FormatInt(int64(value.(int8)), 10)
	case uint8:
		return strconv.FormatUint(uint64(value.(uint8)), 10)
	case int16:
		return strconv.FormatInt(int64(value.(int16)), 10)
	case uint16:
		return strconv.FormatUint(uint64(value.(uint16)), 10)
	case int32:
		return strconv.FormatInt(int64(value.(int32)), 10)
	case uint32:
		return strconv.FormatUint(uint64(value.(uint32)), 10)
	case int64:
		return strconv.FormatInt(value.(int64), 10)
	case uint64:
		return strconv.FormatUint(value.(uint64), 10)
	case int:
		return strconv.FormatInt(int64(value.(int)), 10)
	case uint:
		return strconv.FormatUint(uint64(value.(uint)), 10)
	case float32:
		return strconv.FormatFloat(float64(value.(float32)), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	case string:
		return value.(string)
	case []byte:
		return string(value.([]byte))
	case bool:
		if value.(bool) == true {
			return "true"
		}
		return "false"
	default:
		newValue, _ := json.Marshal(value)
		return string(newValue)
	}
}

func ToInt(value interface{}) int {
	switch value.(type) {
	case int8:
		return int(value.(int8))
	case uint8: //& byte
		return int(value.(uint8))
	case int16:
		return int(value.(int16))
	case uint16:
		return int(value.(uint16))
	case int32:
		return int(value.(int32))
	case uint32:
		return int(value.(uint32))
	case int64:
		return int(value.(int64))
	case uint64:
		return int(value.(uint64))
	case int:
		return value.(int)
	case uint:
		return int(value.(uint))
	case float32:
		return int(value.(float32))
	case float64:
		return int(value.(float64))
	case string:
		v, _ := strconv.Atoi(value.(string))
		return v
	case []byte:
		v, _ := strconv.Atoi(string(value.([]byte)))
		return v
	case bool:
		if value.(bool) == true {
			return 1
		}
		return 0
	default:
		return 0
	}

	return 0
}

func ToFloat32(value interface{}) float32 {
	switch value.(type) {
	case int8:
		return float32(value.(int8))
	case uint8: //& byte
		return float32(value.(uint8))
	case int16:
		return float32(value.(int16))
	case uint16:
		return float32(value.(uint16))
	case int32:
		return float32(value.(int32))
	case uint32:
		return float32(value.(uint32))
	case int64:
		return float32(value.(int64))
	case uint64:
		return float32(value.(uint64))
	case int:
		return float32(value.(int))
	case uint:
		return float32(value.(uint))
	case float32:
		return value.(float32)
	case float64:
		return float32(value.(float64))
	case string:
		v, _ := strconv.ParseFloat(value.(string), 32)
		return float32(v)
	case []byte:
		v, _ := strconv.ParseFloat(string(value.(string)), 32)
		return float32(v)
	case bool:
		if value.(bool) == true {
			return 1
		}
		return 0
	default:
		return 0
	}
	return 0
}

func ToFloat64(value interface{}) float64 {
	switch value.(type) {
	case int8:
		return float64(value.(int8))
	case uint8: //& byte
		return float64(value.(uint8))
	case int16:
		return float64(value.(int16))
	case uint16:
		return float64(value.(uint16))
	case int32:
		return float64(value.(int32))
	case uint32:
		return float64(value.(uint32))
	case int64:
		return float64(value.(int64))
	case uint64:
		return float64(value.(uint64))
	case int:
		return float64(value.(int))
	case uint:
		return float64(value.(uint))
	case float32:
		return float64(value.(float32))
	case float64:
		return value.(float64)
	case string:
		v, _ := strconv.ParseFloat(value.(string), 64)
		return v
	case []byte:
		v, _ := strconv.ParseFloat(string(value.(string)), 64)
		return v
	case bool:
		if value.(bool) == true {
			return 1
		}
		return 0
	default:
		return 0
	}
	return 0
}
