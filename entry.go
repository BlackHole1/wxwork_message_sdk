package wxwork_message_sdk

import (
    "encoding/json"
    "encoding/xml"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"
    "strings"
)

import (
    "github.com/julienschmidt/httprouter"
    "github.com/sbzhu/weworkapi_golang/wxbizmsgcrypt"
)

type ErrorInfo struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
}

type MsgContent struct {
    ToUsername   string `xml:"ToUserName"`
    FromUsername string `xml:"FromUserName"`
    CreateTime   uint32 `xml:"CreateTime"`
    MsgType      string `xml:"MsgType"`
    Content      string `xml:"Content"`
    Msgid        string `xml:"MsgId"`
    Agentid      uint32 `xml:"AgentId"`
}

type Wx struct {
    WxCrypt        *wxbizmsgcrypt.WXBizMsgCrypt
    Path           string
    Port           string
    RegistryHandle map[string]func(textContent string) (content string, err error)
    Delimiters     []string
}

// 因为企业微信在校验消息接收/回复时，会发送校验请求，此函数是专门用来校验的
// 详情可见：https://github.com/BlackHole1/weworkapi_golang/blob/master/sample.go#L24-L38
func verifyController(wxCrypt *wxbizmsgcrypt.WXBizMsgCrypt) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        params := r.URL.Query()
        msgSignature := params.Get("msg_signature")
        timestamp := params.Get("timestamp")
        nonce := params.Get("nonce")
        echostr := params.Get("echostr")

        // 进行解密
        sEchoStr, cryptErr := wxCrypt.VerifyURL(msgSignature, timestamp, nonce, echostr)
        if nil != cryptErr {
            resultBytes, _ := json.Marshal(ErrorInfo{
                Code: cryptErr.ErrCode,
                Msg:  cryptErr.ErrMsg,
            })
            _, _ = fmt.Fprintf(w, string(resultBytes))
            return
        }

        _, _ = fmt.Fprintf(w, string(sEchoStr))
    }
}

// 接收用户发来的消息，并进行回复
func receiveController(wx *Wx) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        params := r.URL.Query()
        msgSignature := params.Get("msg_signature")
        timestamp := params.Get("timestamp")
        nonce := params.Get("nonce")
        xmlData, _ := ioutil.ReadAll(r.Body)

        // 进行解密
        msg, cryptErr := wx.WxCrypt.DecryptMsg(msgSignature, timestamp, nonce, xmlData)
        if nil != cryptErr {
            return
        }

        var msgContent MsgContent
        err := xml.Unmarshal(msg, &msgContent)
        if nil != err {
            return
        }
        content := strings.ToLower(msgContent.Content)
        result := "没有匹配到相应结果"

        prefix, content := formatContext(content, wx.Delimiters)

        switch {
        case nil != wx.RegistryHandle[prefix]:
            result, err = wx.RegistryHandle[prefix](content)
            if nil != err {
                result = err.Error()
            }
        }

        // 目前只支持本文消息
        responseData := "<xml><ToUserName><![CDATA[" + msgContent.ToUsername + "]]></ToUserName><FromUserName><![CDATA[" + msgContent.FromUsername + "]]></FromUserName><CreateTime>" + strconv.FormatUint(uint64(msgContent.CreateTime), 10) + "</CreateTime><MsgType><![CDATA[text]]></MsgType><Content><![CDATA[" + result + "]]></Content><MsgId>" + msgContent.Msgid + "</MsgId><AgentID>" + strconv.FormatUint(uint64(msgContent.Agentid), 10) + "</AgentID></xml>"

        sEncryptMsg, cryptErr := wx.WxCrypt.EncryptMsg(responseData, timestamp, nonce)
        if nil != cryptErr {
            return
        }

        _, _ = fmt.Fprint(w, string(sEncryptMsg))
    }
}

// 第一层函数创建企业微信所需的key
// 第二层函数创建此SDK所需要的配置
func Create(token string, corpId string, encodingAesKey string) func(path string, port string, delimiters []string) *Wx {
    return func(path string, port string, delimiters []string) *Wx {
        wxCrypt := wxbizmsgcrypt.NewWXBizMsgCrypt(token, encodingAesKey, corpId, wxbizmsgcrypt.XmlType)
        return &Wx{
            RegistryHandle: make(map[string]func(textContent string) (content string, err error)),
            WxCrypt:        wxCrypt,
            Path:           path,
            Port:           port,
            Delimiters:     delimiters,
        }
    }
}

// 配置好后启动服务
func (w *Wx) Run() {
    router := httprouter.New()
    router.GET(w.Path, verifyController(w.WxCrypt))
    router.POST(w.Path, receiveController(w))
    fmt.Println("Server Started!")
    log.Fatal(http.ListenAndServe(w.Port, router))
}

// 注册消息
func (w *Wx) Registry(handel func(textContent string) (content string, err error), levels ...string) {
    w.RegistryHandle[strings.ToLower(strings.Join(levels, " "))] = handel
}
