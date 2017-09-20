package umeng

import (
	"testing"
)

const (
	AppKey          = "umeng-app-key"
	AppMasterSecret = "umeng-app-master-secret"
)

type NotifyCustomMsg struct {
	PType string `json:"ptype"` // 客户端自定义类型
	Url   string `json:"body"`  // 客户端自定义协议
}

func TestAndroidSendCustomizedcast(t *testing.T) {
	aNotifier := NewAndroidNotification(AppKey, AppMasterSecret)

	// 必填项
	aNotifier.AliasType = "xxxxx"
	aNotifier.Alias = "30000504" // user_id

	//// test message
	aNotifier.Payload.DisplayType = DISPLAY_TYPE_MESSAGE // 消息 只需custom字段
	var custom NotifyCustomMsg
	custom.PType = "reward_msg"
	custom.Url = "xxxxx"
	aNotifier.Payload.Body.Custom = custom

	// test notification
	aNotifier.Payload.DisplayType = DISPLAY_TYPE_NOTIFICATION
	aNotifier.Payload.Body.Ticker = "通知栏提示文字"
	aNotifier.Payload.Body.Title = "通知标题"
	aNotifier.Payload.Body.Text = "通知文字描述 通知文字描述 "
	aNotifier.Payload.Body.BuilderId = 1
	aNotifier.ProductionMode = "true"
	aNotifier.Payload.Body.AfterOpen = "go_custom"
	//var custom NotifyCustomMsg
	custom.PType = "action_open"
	custom.Url = "xxxxxxxxxxxxxxx"
	aNotifier.Payload.Body.Custom = custom

	aNotifier.SendCustomizedcast()
}
