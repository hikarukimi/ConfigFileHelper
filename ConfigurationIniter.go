package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ConfigReader 用于读取配置文件并将其内容赋值给结构体或字符串
type ConfigReader struct {
	configFilePath string // 配置文件的路径
}

// assignConfigToStructHelper 是一个辅助函数，用于将配置映射到指定的结构体实例中
func assignConfigToStructHelper(config map[string]string, structType interface{}) interface{} {
	reflectVal := reflect.ValueOf(structType)
	reflectVal = reflect.ValueOf(structType).Elem() // 获取指针指向的实际值
	reflectModel := reflectVal.Type()               // 获取类型信息

	// 遍历结构体的所有字段
	for i := 0; i < reflectModel.NumField(); i++ {
		field := reflectModel.Field(i)    // 获取字段的信息
		key := field.Name                 // 字段名作为配置键
		if reflectVal.Field(i).CanSet() { // 检查是否可以直接设置该字段
			value := strings.TrimRight(config[key], "\r") // 去除值右侧的回车符

			// 根据字段类型设置对应的值
			switch field.Type.Kind() {
			case reflect.String:
				reflectVal.Field(i).SetString(value)
			case reflect.Int:
				intValue, err := strconv.Atoi(value)
				if err != nil {
					fmt.Printf("Error setting field %s: %v\n", key, err)
					continue
				}
				reflectVal.Field(i).SetInt(int64(intValue))
			case reflect.Bool:
				boolValue, err := strconv.ParseBool(value)
				if err != nil {
					fmt.Printf("Error setting field %s: %v\n", key, err)
					continue
				}
				reflectVal.Field(i).SetBool(boolValue)
			default:
				fmt.Printf("Unsupported field type for field %s\n", key)
			}
		}
	}
	return reflectVal.Interface() // 返回设置好值的结构体
}

// AssignMapConfigToStruct 从配置文件中读取指定名称的配置，并将其映射到给定的结构体中
func (cr *ConfigReader) AssignMapConfigToStruct(configName string, structType interface{}) interface{} {
	return assignConfigToStructHelper(getConfigFileContext(cr.configFilePath, configName), structType)
}

// AssignSingleConfigToString 从配置文件中读取单个配置项，并将其值赋给提供的字符串指针
func (cr *ConfigReader) AssignSingleConfigToString(configName string, result *string) {
	configValue, _ := getSingleConfig(cr.configFilePath, configName)
	*result = configValue
}
