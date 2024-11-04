package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// getConfigFileContext 将 application.yaml 中与 configName 相关的配置以 map[string]string 的方式返回
func getConfigFileContext(filePath string, configName string) map[string]string {
	fileContent := readFile(filePath)

	result := make(map[string]string)
	// 表示目标配置前的空格个数，以这个指标判别哪些配置是目标配置包含的配置
	targetBlankCount := 0
	// 表示当前行是否应该将 key:value 键值对读取到 result 中
	resultAppendFlag := false

	for _, content := range fileContent {

		// 通过当前行的内容是否与传入的配置项名称一致和缩进空格个数切换是否读入 key:value
		// 每次读取到目标配置项或者判断为已经读取到另一配置项时切换 resultAppendFlag 的状态实现 result 只包含目标配置项的内容而排除其他配置项
		isTargetConfig, currentBlankCount := isSameConfig(content, configName)
		if isTargetConfig {
			// 找到目标配置项，将 targetBlankCount 设置为当前配置项的空格缩进个数
			targetBlankCount = currentBlankCount
		}
		// 开始读取目标配置项和读取完毕目标配置项时切换 resultAppendFlag 的状态
		if (isTargetConfig) || (resultAppendFlag == true && currentBlankCount <= targetBlankCount) {
			resultAppendFlag = !resultAppendFlag
		}
		// 后续条件是为了排除配置项本身，即排除 mysql: 这样的内容
		if resultAppendFlag && currentBlankCount > targetBlankCount {
			splitContent := strings.Split(content, ":")
			// value 从 1 开始是为了去除 : 之后的空格
			key, value := strings.Trim(splitContent[0], " "), splitContent[1][1:]
			result[key] = value
		}
	}
	return result
}

// getSingleConfig 获取单个配置项的值
func getSingleConfig(filePath string, configName string) (string, error) {
	fileContent := readFile(filePath)
	for _, content := range fileContent {

		currentConfig := strings.Split(strings.Trim(content, " "), ":")
		currentConfigKey := currentConfig[0]
		currentConfigValue := currentConfig[1]
		if currentConfigKey == configName {
			return currentConfigValue, nil
		}

	}
	return "", errors.New("not found")
}

// isSameConfig 判断当前行是否为目标配置项，同时返回当前行的缩进空格个数
func isSameConfig(fileContent string, configName string) (bool, int) {
	// 以缩进的空格个数来识别哪些配置是目标配置
	contentBlankCount := blankCount(fileContent)
	fileContent = strings.TrimLeft(fileContent, " ")

	if len(fileContent) < len(configName) {
		return false, contentBlankCount
	}

	// 如果当前行长度为零表示遍历到了文件最末尾
	currentConfig := fileContent[:len(configName)]
	for i := 0; i < len(configName); i++ {
		if currentConfig[i] != configName[i] {
			return false, contentBlankCount
		}
	}
	return true, contentBlankCount
}

// blankCount 返回当前行缩进的空格个数
func blankCount(fileContent string) int {
	count := 0
	for _, char := range fileContent {
		if char != ' ' {
			break
		}
		count++
	}
	return count
}

// readFile 读取整个文件的内容，返回每一行作为一个元素的切片
func readFile(filePath string) []string {
	result := make([]string, 0)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("文件打开失败:", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println("文件关闭异常:", err)
		}
	}()
	reader := bufio.NewReader(file)
	for {
		content, err := reader.ReadString('\n')
		content = strings.TrimSuffix(content, "\n") // 去除每行末尾的换行符
		result = append(result, content)
		if err != nil {
			break
		}
	}
	return result
}
