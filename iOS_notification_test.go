package umeng

import (
	"fmt"
	"testing"
)

func TestIOSSendCustomizedcast(t *testing.T) {
	aNotifier := NewIOSNotification(AppKey, AppMasterSecret)

	// 必填项
	aNotifier.AliasType = "xxxxx"
	aNotifier.Alias = "30000372" // user_id

	//// test notification
	aNotifier.Payload = make(map[string]interface{})
	var aps IOSPayloadAps
	aps.Alert = "恭喜你获得了一个新红包"
	aNotifier.Payload["aps"] = aps // 必填
	aNotifier.Payload["url"] = "xxxxx://xxxxxxxxxxxxx"

	err := aNotifier.SendCustomizedcast()
	if err != nil {
		fmt.Println(err)
	}
}
