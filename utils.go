package wxwork_message_sdk

import (
    "strings"
)

// 正确获取中文下标，从0开始
func unicodeIndex(str, substr string) int {
    result := strings.Index(str, substr)
    if result >= 0 {
        prefix := []byte(str)[0:result]
        rs := []rune(string(prefix))
        result = len(rs)
    }

    return result
}

// 获取当前句子的标记符下标(针对多个标记符的情况)
func getDelimiterIndexForMultiple(content string, delimiters []string) int {
    index := len([]rune(content))

    for _, delimiter := range delimiters {
        currentDelimiterIndex := unicodeIndex(content, delimiter) + 1

        // 开始和最后的标识符不理会
        // 如果没有标记符也不理会
        if currentDelimiterIndex > 1 && currentDelimiterIndex < index && index != 0 {
            index = currentDelimiterIndex
        }
    }

    if index == len([]rune(content)) {
        return 0
    }

    return index
}

// 获取当前句子的标记符下标(针对单个标记符的情况)
func getDelimiterIndex(content string, delimiter string) int {
    index := unicodeIndex(content, delimiter) + 1
    if index > 1 {
        return index
    }
    return index
}

// 格式化句子
func formatContext(content string, delimiters []string) (string, string) {
    index := 0

    if len(delimiters) == 1 {
        index = getDelimiterIndex(content, delimiters[0])
    } else {
        index = getDelimiterIndexForMultiple(content, delimiters)
    }

    if index == 0 {
        return content, content
    }

    // 当前句子的标识符
    delimiter := string([]rune(content)[index-1 : index])

    // 后面是否还有
    behind := strings.Index(string([]rune(content)[:index]), delimiter)

    if behind == -1 {
        return strings.ToLower(string([]rune(content)[:index-1])), content
    }

    prefix := string([]rune(content)[:index - 1])
    content = string([]rune(content)[index:])

    nextPrefix, nextContent := formatContext(content, []string{delimiter})

    content = nextContent
    if nextPrefix != nextContent {
        prefix = prefix + " " + nextPrefix
        content = nextContent
    }

    return strings.Trim(prefix, " "), strings.Trim(content, " ")
}
