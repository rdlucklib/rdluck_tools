package jpushv2

import (
	"fmt"

	"strings"

	"github.com/ylywyn/jpush-api-go-client"
)

var (
	appKey string
	secret string
)

const (
	//==============推送平台=============
	Platform_IOS     = jpushclient.IOS     //ios客户端
	Platform_ANDROID = jpushclient.ANDROID //安卓客户端
	Platform_ALL     = "all"               //所有客户端
	// WINPHONE = "winphone"
	//==============推送目标=============
	Audience_TAG     = jpushclient.TAG     //标签 多个标签之间是 OR 关系，即取并集
	Audience_TAG_AND = jpushclient.TAG_AND //标签 多个标签之间是 AND 关系，即取交集
	Audience_ALIAS   = jpushclient.ALIAS   //别名 多个别名之间是 OR 关系，即取并集
	Audience_ID      = jpushclient.ID      //注册id 多个注册ID之间是 OR 关系，即取并集
	Audience_ALL     = "all"               //所有用户
)

type JPush struct {
	pf           *jpushclient.Platform
	ad           *jpushclient.Audience
	notice       *jpushclient.Notice
	msg          *jpushclient.Message
	option       *jpushclient.Option
	platform     string //平台
	badge        int    //ios角标
	title        string
	content      string
	alert        string
	users        []string //推送目标
	audienceType string   //推送目标类型
	appKey       string
	secret       string
}

//初始化jpush
//批量参数 以,分隔
func InitJPush(jpush_appKey, jpush_secret string) {
	appKey = jpush_appKey
	secret = jpush_secret
}

//实例化一个jpush
func NewJPush(config ...string) *JPush {
	platform := Platform_ALL
	audienceType := Audience_ALL
	s := secret
	k := appKey
	apns := true

	if len(config) > 4 {
		platform = config[0]
		audienceType = config[1]
		if config[2] == "false" {
			apns = false
		}
		s = config[3]
		k = config[4]
	} else if len(config) > 2 {
		platform = config[0]
		audienceType = config[1]
		if config[2] == "false" {
			apns = false
		}
	} else if len(config) > 1 {
		platform = config[0]
		audienceType = config[1]
	} else if len(config) > 0 {
		platform = config[0]
	}

	option := &jpushclient.Option{}
	option.SetApns(apns)
	fmt.Println(s, k)
	return &JPush{
		platform:     platform,
		badge:        1,
		audienceType: audienceType,
		pf:           &jpushclient.Platform{},
		ad:           &jpushclient.Audience{},
		notice:       &jpushclient.Notice{},
		msg:          &jpushclient.Message{},
		option:       option,
		secret:       s,
		appKey:       k,
	}
}

//设置环境 true：生产环境，false：开发环境
func (j *JPush) SetApns(apns bool) {
	j.option.SetApns(apns)
}

//设置推送平台
func (j *JPush) SetPlatform(platform string) {
	j.platform = platform
}

//设置角标 ios有效
func (j *JPush) SetBadge(badge int) {
	j.badge = badge
}

//设置推送目标类型
func (j *JPush) SetAudienceType(audienceType string) {
	j.audienceType = audienceType
}

//设置key
func (j *JPush) SetSK(s, k string) {
	j.secret = s
	j.appKey = k
}

func (j *JPush) buildMessage() {
	//判断系统
	if j.platform == Platform_ANDROID {
		j.pf.Add(Platform_ANDROID)
		j.notice.SetAndroidNotice(&jpushclient.AndroidNotice{Alert: j.alert})
	} else if j.platform == Platform_IOS {
		j.pf.Add(Platform_IOS)
		j.notice.SetIOSNotice(&jpushclient.IOSNotice{Alert: j.alert, Badge: j.badge})
	} else if j.platform == Platform_ALL {
		j.pf.All()
		j.notice.SetAndroidNotice(&jpushclient.AndroidNotice{Alert: j.alert})
		j.notice.SetIOSNotice(&jpushclient.IOSNotice{Alert: j.alert, Badge: j.badge})
	}
	j.msg.Title = j.title
	j.msg.Content = j.content
}

//构建PayLoad
func (j *JPush) buildPayLoad() []byte {
	payload := jpushclient.NewPushPayLoad()
	payload.SetPlatform(j.pf)
	payload.SetNotice(j.notice)
	payload.SetAudience(j.ad)
	payload.SetMessage(j.msg)
	payload.SetOptions(j.option)
	messageBytes, _ := payload.ToBytes()
	fmt.Printf("%s\r\n", string(messageBytes))
	return messageBytes
}

func (j *JPush) pushMessage() {
	j.buildMessage()

	switch j.audienceType {
	case Audience_ALL:
		j.ad.All()
	case Audience_ALIAS:
		j.ad.SetAlias(j.users)
	case Audience_ID:
		j.ad.SetID(j.users)
	case Audience_TAG:
		j.ad.SetTag(j.users)
	case Audience_TAG_AND:
		j.ad.SetTagAnd(j.users)
	}
	messageBytes := j.buildPayLoad()
	if j.secret == "" || j.appKey == "" {
		fmt.Println("secret appKey is null")
	}
	secrets := strings.Split(j.secret, ",")
	appKeys := strings.Split(j.appKey, ",")
	if len(secrets) == len(appKeys) {
		for i := 0; i < len(secrets); i++ {
			client := jpushclient.NewPushClient(secrets[i], appKeys[i])
			str, err := client.Send(messageBytes)
			if err != nil {
				fmt.Println("err: ", err.Error())
			} else {
				fmt.Println("ok: ", str)
			}
		}
	}

}

//发送消息
func (j *JPush) PushMessage(alert string, users ...string) {
	j.alert = alert
	j.users = users
	j.content = "1"
	j.title = "2"
	j.msg.SetTitle("1")
	j.msg.SetContent("2")
	j.pushMessage()
}
