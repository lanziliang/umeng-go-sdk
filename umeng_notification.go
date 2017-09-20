package umeng

import (
	"errors"
	"fmt"
)

var (
	Host       = "http://msg.umeng.com"
	UploadPath = "/upload"
	PostPath   = "/api/send"
)

const (
	TYPE_UNICAST        = "unicast"
	TYPE_LISTCAST       = "listcast"
	TYPE_FILECAST       = "filecast"
	TYPE_BROADCAST      = "broadcast"
	TYPE_GROUPCAST      = "groupcast"
	TYPE_CUSTOMIZEDCAST = "customizedcast"
)

type UmengNotification struct {
	AppMasterSecret string `json:"-" validate:"required"`
	AppKey          string `json:"appkey" validate:"required"`    // 必填 应用唯一标识
	Timestamp       int64  `json:"timestamp" validate:"required"` // 必填 时间戳，10位或者13位均可，时间戳有效期为10分钟
	/*
		必填 消息发送类型,其值可以为:
		unicast-单播
		listcast-列播(要求不超过500个device_token)
		filecast-文件播(多个device_token可通过文件形式批量发送）
		broadcast-广播
		groupcast-组播(按照filter条件筛选特定用户群, 具体请参照filter参数)
		customizedcast(通过开发者自有的alias进行推送),
		  包括以下两种case:
		   - alias: 对单个或者多个alias进行推送
		   - file_id: 将alias存放到文件后，根据file_id来推送
	*/
	Type string `json:"type" validate:"required"`
	/*
		可选 设备唯一表示
		当type=unicast时,必填, 表示指定的单个设备
		当type=listcast时,必填,要求不超过500个,
		多个device_token以英文逗号间隔
	*/
	DeviceTokens   string      `json:"device_tokens,omitempty"`
	AliasType      string      `json:"alias_type,omitempty"` // 可选 当type=customizedcast时，必填，alias的类型,alias_type可由开发者自定义,开发者在SDK中调用setAlias(alias, alias_type)时所设置的alias_type
	Alias          string      `json:"alias,omitempty"`      // 可选 当type=customizedcast时, 开发者填写自己的alias。要求不超过50个alias,多个alias以英文逗号间隔。在SDK中调用setAlias(alias, alias_type)时所设置的alias
	FileId         string      `json:"file_id,omitempty"`
	Filter         interface{} `json:"filter,omitempty"`
	ProductionMode string      `json:"production_mode,omitempty"`
	Description    string      `json:"description,omitempty"`
	ThirdPartyId   string      `json:"thirdparty_id,omitempty"`
}

type UmengNotificationRet struct {
	Ret  string            `json:"ret"`
	Data map[string]string `json:"data"`
}

func (n *UmengNotification) GetApiUrl() string {
	return fmt.Sprintf("%s%s", Host, PostPath)
}

func (n *UmengNotification) SetAppConfig(key, secret string) {
	n.AppKey = key
	n.AppMasterSecret = secret
}

func (n *UmengNotification) validate() error {
	// TODO all type validate
	switch n.Type {
	case TYPE_CUSTOMIZEDCAST:
		if n.AliasType == "" || (n.Alias == "" && n.FileId == "") {
			return errors.New("You need to set alias or upload file for customizedcast!")
		}
	case TYPE_UNICAST:
	case TYPE_LISTCAST:
	case TYPE_FILECAST:
	case TYPE_BROADCAST:
	case TYPE_GROUPCAST:
	default:
		return fmt.Errorf("Type %s unsupported!", n.Type)
	}

	return nil
}
