package umeng

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type IOSNotification struct {
	UmengNotification
	Payload map[string]interface{} `json:"payload" validate:"required"`
	Policy  IOSPolicy              `json:"policy"`
}

type IOSPayloadAps struct {
	Alert            string `json:"alert" validate:"required"`
	Badge            string `json:"badge,omitempty"`
	Sound            string `json:"sound,omitempty"`
	ContentAvailable string `json:"content-available,omitempty"`
	Category         string `json:"category,omitempty"`
}

type IOSPolicy struct {
	StartTime      string `json:"start_time,omitempty"`
	ExpireTime     string `json:"expire_time,omitempty"`
	MaxSendNum     int    `json:"max_send_num,omitempty"`
	ApnsCollapseId string `json:"apns-collapse-id,omitempty"` //可选，iOS10开始生效。
}

func NewIOSNotification(key, secret string) *IOSNotification {
	notifier := new(IOSNotification)
	notifier.AppKey = key
	notifier.AppMasterSecret = secret
	return notifier
}

// 验证数据是否合法
func (n *IOSNotification) validate() error {
	err := validate.Struct(n)
	if err != nil {
		return err
	}

	err = n.UmengNotification.validate()
	if err != nil {
		return err
	}

	// validate payload
	if _, ok := n.Payload["aps"]; !ok {
		return errors.New("You need to set aps for payload!")
	}

	return nil
}

// send the notification to umeng
func (n *IOSNotification) send() error {
	err := n.validate()
	if err != nil {
		return err
	}

	pushBodyByte, err := json.Marshal(n)
	if err != nil {
		return err
	}

	url := n.GetApiUrl()
	string2sign := fmt.Sprintf("%s%s%s%s", "POST", url, string(pushBodyByte), n.AppMasterSecret)

	h := md5.New()
	io.WriteString(h, string2sign)
	sign := fmt.Sprintf("%x", h.Sum(nil))

	timeout := 30 * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	url = fmt.Sprintf("%s?sign=%s", url, sign)

	// debug
	fmt.Println("Umeng iOS notification url: ", url)
	fmt.Println("Umeng iOS notification body: ", string(pushBodyByte))

	resp, err := client.Post(url, "application/json", strings.NewReader(string(pushBodyByte)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// debug
	fmt.Println("Umeng iOS notification return: ", string(body))

	var ret UmengNotificationRet
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return err
	}
	if ret.Ret == "FAIL" {
		return errors.New("ret fail")
	}

	return nil
}

func (n *IOSNotification) SendCustomizedcast() error {
	n.Timestamp = time.Now().Unix() * 1000
	n.Type = TYPE_CUSTOMIZEDCAST

	return n.send()
}
