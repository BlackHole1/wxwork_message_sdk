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

## Screenshots

![Imgur](https://i.imgur.com/ui8IUsG.png)

![Imgur](https://i.imgur.com/8rhMYzG.png)

## 相关资料

[企业微信API - 接收消息与事件](https://work.weixin.qq.com/api/doc#10514)

[企业微信API - golang API](https://github.com/sbzhu/weworkapi_golang)