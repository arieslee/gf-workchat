package message

// MsgType 参考 https://work.weixin.qq.com/api/doc/90000/90135/90239
type MsgType string
type EventType string

// 消息类型
const (
	MsgTypeText     MsgType = "text"
	MsgTypeImage    MsgType = "image"
	MsgTypeVoice    MsgType = "voice"
	MsgTypeVideo    MsgType = "video"
	MsgTypeLocation MsgType = "location"
	MsgTypeLink     MsgType = "link"
	MsgTypeEvent    MsgType = "event"
)

// 事件类型
const (
	EventTypeSubcrible      EventType = "subscribe"
	EventTypeUnSubcrible    EventType = "unsubscribe"
	EventTypeEnterAgent     EventType = "enter_agent"
	EventTypeReportLocation EventType = "location"
)

// MsgContent 用户文本消息内容
type MsgContent struct {
	ToUsername   string `xml:"ToUserName"`
	FromUsername string `xml:"FromUserName"`
	CreateTime   uint32 `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	Msgid        string `xml:"MsgId"`
	Agentid      uint32 `xml:"AgentId"`
}

// MsgContentReply 回复用户文本消息内容
type MsgContentReply struct {
	ToUsername   string `xml:"ToUserName"`
	FromUsername string `xml:"FromUserName"`
	CreateTime   uint32 `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
}

// MsgImageReply 回复用户图片消息内容
type MsgImageReply struct {
	ToUsername    string          `xml:"ToUserName"`
	FromUsername  string          `xml:"FromUserName"`
	CreateTime    uint32          `xml:"CreateTime"`
	MsgType       string          `xml:"MsgType"` // 消息类型，此时固定为：image
	ImageResource []ImageResource `xml:"Image"`
}

// ImageResource 图片资源
type ImageResource struct {
	MediaId string `xml:"MediaId"` // 图片媒体文件id，可以调用获取媒体文件接口拉取
}

// MsgVoiceReply 回复用户语音消息内容
type MsgVoiceReply struct {
	ToUsername    string          `xml:"ToUserName"`
	FromUsername  string          `xml:"FromUserName"`
	CreateTime    uint32          `xml:"CreateTime"`
	MsgType       string          `xml:"MsgType"` // 消息类型，此时固定为：voice
	VoiceResource []VoiceResource `xml:"Voice"`
}

// VoiceResource 语音资源
type VoiceResource struct {
	MediaId string `xml:"MediaId"` // 语音文件id，可以调用获取媒体文件接口拉取
}

// MsgVideoReply 回复用户视频消息内容
type MsgVideoReply struct {
	ToUsername    string          `xml:"ToUserName"`
	FromUsername  string          `xml:"FromUserName"`
	CreateTime    uint32          `xml:"CreateTime"`
	MsgType       string          `xml:"MsgType"` // 消息类型，此时固定为：video
	VideoResource []VideoResource `xml:"Video"`
}

// VideoResource 视频资源
type VideoResource struct {
	MediaId     string `xml:"MediaId"`     // 视频文件id，可以调用获取媒体文件接口拉取
	Title       string `xml:"Title"`       // 视频消息的标题,不超过128个字节，超过会自动截断
	Description string `xml:"Description"` // 视频消息的描述,不超过512个字节，超过会自动截断
}

// MsgNewsReply 回复用户图文消息内容
type MsgNewsReply struct {
	ToUsername   string         `xml:"ToUserName"`
	FromUsername string         `xml:"FromUserName"`
	CreateTime   uint32         `xml:"CreateTime"`
	MsgType      string         `xml:"MsgType"`      // 消息类型，此时固定为：news
	ArticleCount uint32         `xml:"ArticleCount"` // 图文消息的数量
	NewsResource []NewsResource `xml:"Articles"`     // 图文资源
}

// NewsResource 图文资源
type NewsResource struct {
	Items []News `xml:"item"` // 图文列表项
}

// News 图文信息
type News struct {
	Title       string `xml:"Title"`       // 视频消息的标题,不超过128个字节，超过会自动截断
	Description string `xml:"Description"` // 视频消息的描述,不超过512个字节，超过会自动截断
	Url         string `xml:"Url"`         // 点击后跳转的链接。
	PicUrl      string `xml:"PicUrl"`      // 图文消息的图片链接，支持JPG、PNG格式，较好的效果为大图640320，小图8080。
}

// MsgTaskCardReply 回复用户任务卡片
type MsgTaskCardReply struct {
	ToUsername       string             `xml:"ToUserName"`
	FromUsername     string             `xml:"FromUserName"`
	CreateTime       uint32             `xml:"CreateTime"`
	MsgType          string             `xml:"MsgType"` // 消息类型，此时固定为：update_taskcard
	TaskCardResource []TaskCardResource `xml:"TaskCard"`
}

// TaskCardResource 任务卡片资源
type TaskCardResource struct {
	ReplaceName string `xml:"ReplaceName"` // 点击任务卡片按钮后显示的按钮名称
}

// ==========================接受的消息===========================================

// MsgImage 用户图片消息内容
type MsgImage struct {
	ToUsername   string `xml:"ToUserName"`   // 企业微信CorpID
	FromUsername string `xml:"FromUserName"` // 成员UserID
	CreateTime   uint32 `xml:"CreateTime"`   // 消息创建时间（整型）
	MsgType      string `xml:"MsgType"`      // 消息类型，此时固定为：image
	PicUrl       string `xml:"PicUrl"`       // 图片链接
	MediaId      string `xml:"MediaId"`      // 图片媒体文件id，可以调用获取媒体文件接口拉取，仅三天内有效
	Msgid        string `xml:"MsgId"`        // 消息id，64位整型
	Agentid      uint32 `xml:"AgentId"`      // 企业应用的id，整型。可在应用的设置页面查看
}

// MsgVoice 用户语音消息内容
type MsgVoice struct {
	ToUsername   string `xml:"ToUserName"`   // 企业微信CorpID
	FromUsername string `xml:"FromUserName"` // 成员UserID
	CreateTime   uint32 `xml:"CreateTime"`   // 消息创建时间（整型）
	MsgType      string `xml:"MsgType"`      // 消息类型，此时固定为：voice
	Format       string `xml:"Format"`       // 语音格式，如amr，speex等
	MediaId      string `xml:"MediaId"`      // 语音媒体文件id，可以调用获取媒体文件接口拉取数据，仅三天内有效
	Msgid        string `xml:"MsgId"`        // 消息id，64位整型
	Agentid      uint32 `xml:"AgentId"`      // 企业应用的id，整型。可在应用的设置页面查看
}

// MsgVideo 用户视频消息内容
type MsgVideo struct {
	ToUsername   string `xml:"ToUserName"`   // 企业微信CorpID
	FromUsername string `xml:"FromUserName"` // 成员UserID
	CreateTime   uint32 `xml:"CreateTime"`   // 消息创建时间（整型）
	MsgType      string `xml:"MsgType"`      // 消息类型，此时固定为：video
	ThumbMediaId string `xml:"ThumbMediaId"` // 视频消息缩略图的媒体id，可以调用获取媒体文件接口拉取数据，仅三天内有效
	MediaId      string `xml:"MediaId"`      // 视频媒体文件id，可以调用获取媒体文件接口拉取数据，仅三天内有效
	Msgid        string `xml:"MsgId"`        // 消息id，64位整型
	Agentid      uint32 `xml:"AgentId"`      // 企业应用的id，整型。可在应用的设置页面查看
}

// MsgLocation 用户位置消息内容
type MsgLocation struct {
	ToUsername   string `xml:"ToUserName"`   // 企业微信CorpID
	FromUsername string `xml:"FromUserName"` // 成员UserID
	CreateTime   uint32 `xml:"CreateTime"`   // 消息创建时间（整型）
	MsgType      string `xml:"MsgType"`      // 消息类型，此时固定为： location
	Location_X   string `xml:"Location_X"`   // 地理位置纬度
	Location_Y   string `xml:"Location_Y"`   // 地理位置经度
	Scale        uint32 `xml:"Scale"`        // 地图缩放大小
	Label        string `xml:"Label"`        //地理位置信息
	Msgid        string `xml:"MsgId"`        // 消息id，64位整型
	Agentid      uint32 `xml:"AgentId"`      // 企业应用的id，整型。可在应用的设置页面查看
}

// MsgLink 用户链接消息内容
type MsgLink struct {
	ToUsername   string `xml:"ToUserName"`   // 企业微信CorpID
	FromUsername string `xml:"FromUserName"` // 成员UserID
	CreateTime   uint32 `xml:"CreateTime"`   // 消息创建时间（整型）
	MsgType      string `xml:"MsgType"`      // 消息类型，此时固定为：link
	Title        string `xml:"Title"`        // 标题
	Description  string `xml:"Description"`  // 描述
	Url          string `xml:"Url"`          // 链接跳转的url
	PicUrl       string `xml:"PicUrl"`       // 封面缩略图的url
	Msgid        string `xml:"MsgId"`        // 消息id，64位整型
	Agentid      uint32 `xml:"AgentId"`      // 企业应用的id，整型。可在应用的设置页面查看
}

// EventSub 用户订阅事件
type EventSub struct {
	ToUsername   string `xml:"ToUserName"`   //企业微信CorpID
	FromUsername string `xml:"FromUserName"` //成员UserID
	CreateTime   uint32 `xml:"CreateTime"`   //消息创建时间（整型）
	MsgType      string `xml:"MsgType"`      //消息类型，此时固定为：event
	Event        string `xml:"Event"`        //事件类型，subscribe(关注)、unsubscribe(取消关注)
	EventKey     string `xml:"EventKey"`     //事件KEY值，此事件该值为空
	Agentid      uint32 `xml:"AgentId"`      //企业应用的id，整型。可在应用的设置页面查看
}

// EventEnterAgent 用户进入应用事件
type EventEnterAgent struct {
	ToUsername   string `xml:"ToUserName"`   //企业微信CorpID
	FromUsername string `xml:"FromUserName"` //成员UserID
	CreateTime   uint32 `xml:"CreateTime"`   //消息创建时间（整型）
	MsgType      string `xml:"MsgType"`      //消息类型，此时固定为：event
	Event        string `xml:"Event"`        //事件类型：enter_agent
	EventKey     string `xml:"EventKey"`     //事件KEY值，此事件该值为空
	Agentid      uint32 `xml:"AgentId"`      //企业应用的id，整型。可在应用的设置页面查看
}

// EventReportLocation 上报地理位置事件
type EventReportLocation struct {
	ToUsername   string `xml:"ToUserName"`   //企业微信CorpID
	FromUsername string `xml:"FromUserName"` //成员UserID
	CreateTime   uint32 `xml:"CreateTime"`   //消息创建时间（整型）
	MsgType      string `xml:"MsgType"`      //消息类型，此时固定为：event
	Event        string `xml:"Event"`        //事件类型：LOCATION
	Latitude     string `xml:"Latitude"`     // 地理位置纬度
	Longitude    string `xml:"Longitude"`    // 地理位置经度
	Precision    string `xml:"Precision"`    //地理位置精度
	EventKey     string `xml:"EventKey"`     //事件KEY值，此事件该值为空
	Agentid      uint32 `xml:"AgentId"`      //企业应用的id，整型。可在应用的设置页面查看
	AppType      string `xml:"AppType"`      //app类型，在企业微信固定返回wxwork，在微信不返回该字段
}

// MsgTypeResp 获取用户消息类型响应
type MsgTypeResp struct {
	MsgType string `xml:"MsgType"` // 消息类型
}

// EventTypeResp 获取用户事件类型响应，在MsgType为event时查询
type EventTypeResp struct {
	Event string `xml:"Event"` // 事件类型
}
