package main

import (
	"reflect"
	"strings"
)

// assignConfigToStruct (将从application读取到的相关配置导入到传入的结构体中)
func assignConfigToStructHelper(config map[string]string, structType interface{}) interface{} {
	reflectVal := reflect.ValueOf(structType)
	reflectVal = reflect.ValueOf(structType).Elem()
	reflectModel := reflectVal.Type()

	for i := 0; i < reflectModel.NumField(); i++ {
		filed := reflectModel.Field(i)
		key := filed.Name
		if reflectVal.Field(i).CanSet() {
			value := strings.TrimRight(config[key], "\r")
			reflectVal.Field(i).SetString(value)
		}
	}
	return reflectVal.Interface()
}

// AssignConfigToStruct (向外界暴露的接口，传入配置项名字和结构体即可将对应的配置赋值给结构体)
func AssignConfigToStruct(configName string, structType interface{}) interface{} {
	return assignConfigToStructHelper(getConfigFileContext("application.yaml", configName), structType)
}
