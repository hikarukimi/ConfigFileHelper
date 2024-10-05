package ConfigFileHelper

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
	result := make(map[string]string)
	resultAppendFlag := false

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("文件关闭异常")
		}
	}(file)

	if err != nil {
		fmt.Println("文件打开失败")
	}
	reader := bufio.NewReader(file)
	for {
		content, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
		}
		//如果当前行开头带有空格表名他只是一个单独的配置子项,而不是希望的配置项
		//每次读取到配置项时或者判断为相同配置项时切换resultAppendFlag的状态实现result只包含目标配置项的内容而排除其他配置项
		if content[0] != ' ' || isSameConfig(content, configName) {
			resultAppendFlag = !resultAppendFlag
		}
		//后续条件是为了排除配置项本身 即排除mysql:这样的内容
		if resultAppendFlag && content[0] == ' ' {
			splitContent := strings.Split(content, ":")
			//value从1开始是为了去除:之后的空格,到len(splitContent[1])-1结束是为了删去多余的\n
			key, value := strings.Trim(splitContent[0], " "), splitContent[1][1:len(splitContent[1])-1]
			result[key] = value
		}
	}
	return result
}

func isSameConfig(fileContent string, configName string) bool {
	//如果当前行长度为零表示遍历到了文件最末尾
	currentConfig := fileContent[:len(configName)]
	for i := 0; i < len(configName); i++ {
		if currentConfig[i] != configName[i] {
			return false
		}
	}
	return true
}
