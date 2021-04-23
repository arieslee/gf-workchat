package config

type Config struct {
	CropId          string `json:"corp_id"`           //corpid
	AppSecret      string `json:"app_secret"`       //appsecret
	Token          string `json:"token"`            //token
	EncodingAESKey string `json:"encoding_aes_key"` //EncodingAESKey
}
