package utils

import (
	"reflect"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
)

func MapDecode(data map[string]string, result interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: stringToTimeHookFunc(),
		Result:     result,
	})
	if err != nil {
		return err
	}
	err = decoder.Decode(data)
	return err
}

func stringToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		str := data.(string)
		switch t.Kind() {
		case reflect.Struct:
			if t == reflect.TypeOf(time.Time{}) {
				if value, err := strconv.ParseInt(str, 10, 64); err == nil {
					return time.UnixMilli(value), nil
				}
			}
		case reflect.Bool:
			return str == "1" || str == "true", nil
		case reflect.Int:
			if val, err := strconv.Atoi(str); err == nil {
				return val, nil
			}
		}
		return data, nil
	}
}
