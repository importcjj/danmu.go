package main

import (
	"fmt"

	"github.com/importcjj/danmu.go/douyu"
)

// 默认弹幕服务器
const (
	DefaultDouyuDanmuHost = "openbarrage.douyutv.com"
	DefaultDouyuDanmuPort = 8601
)

// DanmuHandle 为自定义的弹幕处理
func DanmuHandle(message *douyu.Message) {
	contentType, ok := message.Field("type")
	if !ok {
		return
	}
	switch contentType {
	case douyu.TypeChatMsg:
		// 默认全部为string
		nick, _ := message.Field("nn")
		level, _ := message.Field("level")
		text, _ := message.Field("txt")
		fmt.Printf("<level %s> - %s >>> %s\n", level, nick, text)
	case douyu.TypeUserEnter:
		nick, _ := message.Field("nn")
		level, _ := message.Field("level")
		fmt.Printf("!!!!!欢迎<lv %s> %s 进入房间\n", level, nick)
	}
}

func main() {
	douyuClient := douyu.New()
	douyuClient.Connect(DefaultDouyuDanmuHost, DefaultDouyuDanmuPort)
	douyuClient.JoinRoom(288016)
	douyuClient.HandleFunc(DanmuHandle)
	douyuClient.Watch()
}
