## wxwork_message_sdk

基于企业微信的消息加密/解密库进行的二次封装，使用方不用去关心如何通信、加密、解密。只需要注册相关的词语的层级就好。

*目前只支持文字消息*

> 消息加密/解密库：[weworkapi_golang](https://github.com/sbzhu/weworkapi_golang)

## Install

#### go.mod

添加`github.com/BlackHole1/wxwork_message_sdk v1.0.0`到`go.mod`文件里

```
require (
    //...
    github.com/BlackHole1/wxwork_message_sdk v1.0.0
)
```

安装依赖
```shell
$ go mod tidy
```

#### import

直接在项目中引入

```go
import (
    "github.com/BlackHole1/wxwork_message_sdk"
)
```

## Usage

```go
package main

import (
    "github.com/BlackHole1/wxwork_message_sdk"
)

func help(content string) (string, error) {
    return "I am help", nil
}

func aocTest(content string) (string, error) {
    return "hi" + content, nil
}

func main() {
    initWxKey := wxwork_message_sdk.Create("54DdTgvctgbIMwOpPl7", "wwfaa0e4cbfcdd3ef0", "4LRnApxTnSan9VuumDmBbp7F0z6ufnHzCd4FraB7IRz")
    company := initWxKey("/", ":8080", []string{" ", ":", "："})
    company.Registry(help, "help")
    company.Registry(aocTest, "AOC", "test")
    company.Run()
}
```

## API

### Create

创建一个消息实例，其用法为：

```go
wxwork_message_sdk.Create(token string, corpId string, encodingAesKey string)(path string, port string, delimiters []string)
```

`token`、`corpId`、`encodingAesKey`为企业微信提供：

* token 在建立第三方应用时，企业微信提供
* corpId 为企业号，详情可见：https://work.weixin.qq.com/api/doc#90000/90135/90665
* encodingAesKey 在建立第三方应用时，企业微信提供

以上是第一层的参数，`Create`会返回一个函数，这个函数也有三个参数：

* path 当前接受消息的url路径
* port 要接受消息的端口，非数值型，而是字符串，例如：`":8080"`
* delimiters 分隔符，用于把每句话安装分隔符进行分割，可有多个

### Registry

注册一些语句，根据不用的语句，来调用相关函数。：

单层级，且没有`content`

```go
func aoc(content: string) (string, error) {
    // 判断是否为 aoc help
    if (content == "help") {
        return "I'm Help"
    }
    return "没有找到匹配条件，可输入aoc help查看帮助"
}

// 注册aoc消息
company.Registry(AocHelp, "aoc")
```

多层级捕获

```go
// 此函数，将捕获：
// aoc 日志 login //=> content就为login
func aocLog(content: string) (string, error) {
    // 查询aoc下的login日志
    return "输出结果"
}

// 注册 aoc 日志 消息
company.Registry(aocLog, "aoc", "日志")
```

### Run

当消息注册完毕后，就可以通过此函数，进行启动

## Screenshots

![Imgur](https://i.imgur.com/ui8IUsG.png)

![Imgur](https://i.imgur.com/8rhMYzG.png)

## 相关资料

[企业微信API - 接收消息与事件](https://work.weixin.qq.com/api/doc#10514)

[企业微信API - golang API](https://github.com/sbzhu/weworkapi_golang)