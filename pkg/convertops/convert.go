package convertops

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/max-gui/logagent/pkg/logagent"
	"gopkg.in/yaml.v2"
)

/**
获取val变量对应类型的字符串值
*/
func StrValOfType(val interface{}) string {
	return reflect.TypeOf(val).String()
}

func StrValOfInterface(val interface{}) string {
	return fmt.Sprintf("%v", val)
}

/**
将yaml字符串反序列化成map对象
*/
func ConvertYamlToMap(ymlString string, c context.Context) map[interface{}]interface{} {
	logger := logagent.Inst(c)
	m := make(map[interface{}]interface{})

	err := yaml.Unmarshal([]byte(ymlString), &m)
	if err != nil {
		logger.Panicf("convert yaml to map error: %v", err)
	}
	return m
}

/**
将map序列化成yaml字符串
*/
func ConvertMapToYaml(m *map[interface{}]interface{}, c context.Context) string {
	out, err := yaml.Marshal(m)
	log := logagent.Inst(c)

	if err != nil {
		log.Panicf("error: %v", err)
	}
	return string(out)
}

/**
将string map序列化成yaml字符串
*/
func ConvertStrMapToYaml(m *map[string]interface{}, c context.Context) string {
	out, err := yaml.Marshal(m)
	log := logagent.Inst(c)

	if err != nil {
		log.Panicf("error: %v", err)
	}
	return string(out)
}

func CompareTwoMapInterface(data1 map[string]interface{},
	data2 map[string]interface{}) bool {
	YekSlice := make([]string, 0)
	dataSlice1 := make([]interface{}, 0)
	dataSlice2 := make([]interface{}, 0)
	for Yek, value := range data1 {
		YekSlice = append(YekSlice, Yek)
		dataSlice1 = append(dataSlice1, value)
	}
	for _, Yek := range YekSlice {
		if data, ok := data2[Yek]; ok {
			dataSlice2 = append(dataSlice2, data)
		} else {
			return false
		}
	}
	data1b, _ := json.Marshal(dataSlice1)
	data2b, _ := json.Marshal(dataSlice2)

	dataStr1 := string(data1b)
	dataStr2 := string(data2b)

	return strings.Compare(dataStr1, dataStr2) == 0
}

func Rndintstr(lens int) string {
	// var result string
	// rand.Seed(time.Now().UnixNano())

	// for i := 0; i < lens; i++ {
	// 	result += strconv.Itoa(rand.Intn(8) + 1)
	// }

	return RndRangestr(lens, 1, 9)
}

func RndRangestr(lens, low, high int) string {
	var result string
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < lens; i++ {
		result += strconv.Itoa(rand.Intn(high-low) + low)
	}

	return result
}
