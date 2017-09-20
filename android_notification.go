package umeng

import (
	"fmt"
	"errors"
	"encoding/json"
	"crypto/md5"
	"io"
	"strings"
	"time"
	"net/http"
	"io/ioutil"
)

const (
	DISPLAY_TYPE_NOTIFICATION = "notification"
	DISPLAY_TYPE_MESSAGE = "message"

	AFTER_OPEN_GO_APP = "go_app"
	AFTER_OPEN_GO_URL = "go_url"
	AFTER_OPEN_GO_ACTIVITY = "go_activity"
	AFTER_OPEN_GO_CUSTOM = "go_custom"
)

type AndroidNotification struct {
	UmengNotification
	Payload AndroidPayload `json:"payload" validate:"required"`
	Policy AndroidPolicy `json:"policy,omitempty"`
}

type AndroidPayload struct {
	DisplayType string `json:"display_type" validate:"required"`
	Body  AndroidPayloadBody `json:"body" validate:"required"`
	Extra map[string]string `json:"extra,omitempty"`
}

type AndroidPayloadBody struct {
	Ticker string `json:"ticker,omitempty"`
	Title string `json:"title,omitempty"`
	Text string `json:"text,omitempty"`
	Icon string `json:"icon,omitempty"`
	LargeIcon string `json:"largeIcon,omitempty"`
	Img string `json:"img,omitempty"`
	Sound string `json:"sound,omitempty"`
	BuilderId int `json:"builder_id,omitempty"`
	PlayVibrate string `json:"play_vibrate,omitempty"` // 可选 收到通知是否震动,默认为"true".
	PlayLights string `json:"play_lights,omitempty"`
	PlaySound string `json:"play_sound,omitempty"`
	AfterOpen string `json:"after_open,omitempty"`
	Url string `json:"url,omitempty"`
	Activity string `json:"activity,omitempty"`
	Custom interface{} `json:"custom,omitempty"` // 可选 display_type=message, 或者display_type=notification且"after_open"为"go_custom"时，该字段必填。用户自定义内容, 可以为字符串或者JSON格式。
}

type AndroidPolicy struct {
	StartTime string `json:"start_time,omitempty"`
	ExpireTime string `json:"expire_time,omitempty"`
	MaxSendNum int `json:"max_send_num,omitempty"`
	OutBizNo string `json:"out_biz_no,omitempty"`
}

func NewAndroidNotification(key, secret string) *AndroidNotification {
	notifier := new(AndroidNotification)
	notifier.AppKey = key
	notifier.AppMasterSecret = secret
	return notifier
}

func (n *AndroidPayload) validate() error {
	// validate payload display_type
	switch n.DisplayType {
	case DISPLAY_TYPE_NOTIFICATION:
		if n.Body.Ticker == "" || n.Body.Title == "" || n.Body.Text == "" || n.Body.AfterOpen == "" {
			return errors.New("You need to set ticker and title and text and after_open for display_type notification!")
		}
		// validate after_open
		switch n.Body.AfterOpen {
		case AFTER_OPEN_GO_APP:
		case AFTER_OPEN_GO_URL:
			if n.Body.Url == "" {
				return errors.New("You need to set url for after_open go_url!")
			}
		case AFTER_OPEN_GO_ACTIVITY:
			if n.Body.Activity == "" {
				return errors.New("You need to set activity for after_open go_url!")
			}
		case AFTER_OPEN_GO_CUSTOM:
			if n.Body.Custom == nil {
				return errors.New("You need to set custom for after_open go_custom!")
			}
		default:
			return fmt.Errorf("After open %s unsupported!", n.Body.AfterOpen)
		}

	case DISPLAY_TYPE_MESSAGE: // body的内容只需填写custom字段
		if n.Body.Custom == nil {
			return errors.New("You need to set custom for display_type message!")
		}
	default:
		return fmt.Errorf("Display type %s unsupported!", n.DisplayType)
	}

	return nil
}

// 验证数据是否合法
func (n *AndroidNotification) validate() error {
	err := validate.Struct(n)
	if err != nil {
		return err
	}

	err = n.UmengNotification.validate()
	if err != nil {
		return err
	}

	// validate payload
	err = n.Payload.validate()
	if err != nil {
		return err
	}

	return nil
}

// send the notification to umeng
func (n *AndroidNotification) send() error {
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
	fmt.Println("Umeng android notification url: ", url)
	fmt.Println("Umeng android notification body: ", string(pushBodyByte))

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
	fmt.Println("Umeng android notification return: ", string(body))

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

func (n *AndroidNotification) SendCustomizedcast() error {
	n.Timestamp = time.Now().Unix() * 1000
	n.Type = TYPE_CUSTOMIZEDCAST

	return n.send()
}