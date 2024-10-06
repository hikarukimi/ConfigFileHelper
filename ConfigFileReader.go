package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// getConfigFileContext (将application.yaml中与configName相关的配置以map[string][string]的方式返回)
func getConfigFileContext(filePath string, configName string) map[string]string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("文件打开失败")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("文件关闭异常")
		}
	}(file)

	result := make(map[string]string)
	//表示目标配置前的空格个数，以这个指标判别哪一些配置是目标配置包含的配置
	targetBlankCount := 0
	//表示当前行是否应该应该将key:value键值对读取到result中
	resultAppendFlag := false

	reader := bufio.NewReader(file)
	for {
		content, err := reader.ReadString('\n')

		//通过当前行的内容是否与传入的配置项名称一致和缩进空格个数切换是否读入key:value
		//每次读取到目标配置项或者判断为已经读取到另一配置项时切换resultAppendFlag的状态实现result只包含目标配置项的内容而排除其他配置项
		isTargetConfig, currentBlankCount := isSameConfig(content, configName)
		if isTargetConfig {
			//找到目标配置项，将targetBlankCount设置为当前配置项的空格缩进个数
			targetBlankCount = currentBlankCount
		}
		//开始读取目标配置项和读取完毕目标配置项时切换resultAppendFlag的状态
		if (isTargetConfig) || (resultAppendFlag == true && currentBlankCount <= targetBlankCount) {
			resultAppendFlag = !resultAppendFlag
		}
		//后续条件是为了排除配置项本身 即排除mysql:这样的内容
		if resultAppendFlag && currentBlankCount > targetBlankCount {
			splitContent := strings.Split(content, ":")
			//value从1开始是为了去除:之后的空格
			key, value := strings.Trim(splitContent[0], " "), splitContent[1][1:]
			result[key] = value
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
		}
	}
	return result
}

func isSameConfig(fileContent string, configName string) (bool, int) {
	//以缩进的空格个数来识别哪一些配置是目标配置
	contentBlankCount := blankCount(fileContent)
	fileContent = strings.TrimLeft(fileContent, " ")

	if len(fileContent) < len(configName) {
		return false, contentBlankCount
	}

	//如果当前行长度为零表示遍历到了文件最末尾
	currentConfig := fileContent[:len(configName)]
	for i := 0; i < len(configName); i++ {
		if currentConfig[i] != configName[i] {
			return false, contentBlankCount
		}
	}
	return true, contentBlankCount
}
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
