package config

type Config struct {
	CropId         string `json:"corp_id"`          // corpid
	AppSecret      string `json:"app_secret"`       // appsecret
	Token          string `json:"token"`            // token
	EncodingAESKey string `json:"encoding_aes_key"` // EncodingAESKey
	ReceiverId     string `json:"receiver_id"`      // 企业应用的回调，表示corpid 第三方事件的回调，表示suiteid
}
