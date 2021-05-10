package msg_body

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"fmt"
	"gf_workchat/config"
	"github.com/gogf/gf/encoding/gbase64"
	"sort"
)

type MsgBody struct {
	config *config.Config
}

type WXBizJsonMsg4Recv struct {
	Tousername string `json:"tousername"`
	Encrypt    string `json:"encrypt"`
	AgentId    string `json:"agentid"`
}

type WXBizJsonMsg4Send struct {
	Encrypt      string `json:"encrypt"`
	MsgSignature string `json:"msg_signature"`
	TimeStamp    string `json:"timestamp"`
	Nonce        string `json:"nonce"`
}

const blockSize = 32

func NewMsgBody(cfg *config.Config) *MsgBody {
	return &MsgBody{
		config: cfg,
	}
}

func (m *MsgBody) Decrypt(msgSignature, timestamp, nonce, base64EncryptMsg string) ([]byte, error) {
	localSignature := m.Signature(m.config.Token, timestamp, nonce)
	if msgSignature != localSignature {
		return nil, fmt.Errorf("msg signature is invalid,msgSignature = %s, localSignature=%s", msgSignature, localSignature)
	}
	aesKey, err := gbase64.DecodeString(m.config.EncodingAESKey)
	if nil != err {
		return nil, fmt.Errorf("EncodingAESKey解析失败, err=%v", err)
	}
	encryptMsg, err := gbase64.DecodeString(base64EncryptMsg)
	if nil != err {
		return nil, fmt.Errorf("base64EncryptMsg解析失败, err=%v", err)
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("NewCipher失败, err=%v", err)
	}
	if len(encryptMsg) < aes.BlockSize {
		return nil, fmt.Errorf("encrypt_msg size is not valid, err=%v", err)
	}
	iv := aesKey[:aes.BlockSize]
	if len(encryptMsg)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("encrypt_msg not a multiple of the block size, err=%v", err)
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptMsg, encryptMsg)
	return encryptMsg, nil
}
func (m *MsgBody) Signature(token, timestamp, nonce string) string {
	strs := sort.StringSlice{token, timestamp, nonce}
	sort.Strings(strs)
	str := ""
	for _, s := range strs {
		str += s
	}
	h := sha1.New()
	h.Write([]byte(str))
	signatureNow := fmt.Sprintf("%x", h.Sum(nil))
	return signatureNow
}
func (m *MsgBody) Encrypt(plainText string) ([]byte, error) {
	aesKey, err := gbase64.DecodeString(m.config.EncodingAESKey)
	if nil != err {
		return nil, fmt.Errorf("EncodingAESKey解析失败, err=%v", err)
	}
	padMsg := m.pKCS7Padding(plainText)
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("aes NewCipher Error, err=%v", err)
	}
	cipherText := make([]byte, len(padMsg))
	iv := aesKey[:aes.BlockSize]
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, padMsg)
	base64Msg := gbase64.Encode(cipherText)
	return base64Msg, nil
}

func (m *MsgBody) pKCS7Padding(plainText string) []byte {
	padding := blockSize - (len(plainText) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	var buffer bytes.Buffer
	buffer.WriteString(plainText)
	buffer.Write(padText)
	return buffer.Bytes()
}
