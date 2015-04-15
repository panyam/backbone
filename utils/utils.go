package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strconv"
)

func RandSeq(chars string, n int) string {
	letters := []rune(chars)
	numLetters := len(letters)
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(numLetters)]
	}
	return string(b)
}

func RandDigits(n int) string {
	return RandSeq("0123456789", n)
}

func RandAlnum(n int) string {
	return RandSeq("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", n)
}

func RandAlpha(n int) string {
	return RandSeq("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", n)
}

func ID2String(id int64) string {
	return fmt.Sprintf("%d", id)
}

func String2ID(strid string) int64 {
	val, err := strconv.ParseInt(strid, 10, 64)
	if err != nil {
		return 0
	}
	return val
}

func JsonNumberToInt(value interface{}) (int, error) {
	out, err := JsonNumberToInt64(value)
	return int(out), err
}

func JsonNumberToInt64(value interface{}) (int64, error) {
	out, ok := value.(int64)
	if ok {
		return out, nil
	}

	outNum, ok := value.(json.Number)
	if ok {
		return outNum.Int64()
	}
	return 0, errors.New("Invalid int64")
}

func GetMethod(i interface{}, methodName string) reflect.Value {
	var ptr reflect.Value
	var value reflect.Value
	var finalMethod reflect.Value

	value = reflect.ValueOf(i)

	// if we start with a pointer, we need to get value pointed to
	// if we start with a value, we need to get a pointer to that value
	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	} else {
		ptr = reflect.New(reflect.TypeOf(i))
		temp := ptr.Elem()
		temp.Set(value)
	}

	log.Println("Ptr: ", ptr)
	log.Println("Value: ", value)

	// check for method on value
	method := value.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}
	// check for method on pointer
	method = ptr.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}
	return finalMethod
}

func CallMethod(method reflect.Value, args []reflect.Value) []interface{} {
	if method.IsValid() {
		return_values := method.Call(args)
		output := make([]interface{}, 0, len(return_values))
		for _, value := range return_values {
			output = append(output, value.Interface())
		}
		return output
	}

	// return or panic, method not found of either type
	return nil
}

func GetParamType(method reflect.Value, paramIndex int) (bool, reflect.Type) {
	if method.IsValid() {
		methodType := method.Type()
		paramType := methodType.In(paramIndex)
		if paramType.Kind() == reflect.Ptr {
			return true, paramType.Elem()
		} else {
			return false, paramType
		}
	}
	return false, nil
}
