package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"gf_workchat/config"
	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/gf/util/grand"
	"sort"
	"strings"
)

var (
	ErrorInvalidMsgSignature = errors.New("msg signature is invalid")
)

type Crypto struct {
	config *config.Config
}

type WXBizJsonMsg4Recv struct {
	ToUsername string `json:"tousername"`
	Encrypt    string `json:"encrypt"`
	Agentid    string `json:"agentid"`
}

type WXBizJsonMsg4Send struct {
	Encrypt   string `json:"encrypt"`
	Signature string `json:"msgsignature"`
	Timestamp string `json:"timestamp"`
	Nonce     string `json:"nonce"`
}
type WXBizParseBody struct {
	Random     string `json:"random"`
	Msg        string `json:"msg"`
	MsgLen     uint32 `json:"msg_len"`
	ReceiverId string `json:"receiver_id"`
}
type WXBizEchoStr struct {
}

const blockSize = 32

func New(cfg *config.Config) *Crypto {
	return &Crypto{
		config: cfg,
	}
}
func (c *Crypto) GetAesKey() ([]byte, error) {
	aesKey, err := gbase64.DecodeString(c.config.EncodingAESKey + "=")
	if nil != err {
		return nil, fmt.Errorf("EncodingAESKey解析失败, err=%v", err)
	}
	return aesKey, nil
}

// VerifyURL 验证url
// 企业开启回调模式时，企业微信会向验证url发送一个get请求
// 假设点击验证时，企业收到类似请求：
// GET /cgi-bin/wxpush?msg_signature=5c45ff5e21c57e6ad56bac8758b79b1d9ac89fd3&timestamp=1409659589&nonce=263014780&echostr=P9nAzCzyDtyTWESHep1vC5X9xho%2FqYX3Zpb4yKa9SKld1DsH3Iyt3tP3zNdtp%2B4RPcs8TgAE7OaBO%2BFZXvnaqQ%3D%3D
// HTTP/1.1 Host: qy.weixin.qq.com
// 接收到该请求时，企业应
// 1.解析出Get请求的参数，包括消息体签名(msg_signature)，时间戳(timestamp)，随机数字串(nonce)以及企业微信推送过来的随机加密字符串(echostr),
// 2.验证消息体签名的正确性
// 3. 解密出echostr原文，将原文当作Get请求的response，返回给企业微信
// 第2，3步可以用VerifyURL来实现。
// 原样输出VerifyURL中返回的string即可
func (c *Crypto) VerifyURL(msgSignature, timestamp, nonce, echoStr string) (string, error) {
	localSignature := c.GetSignature(c.config.Token, timestamp, nonce, echoStr)
	if msgSignature != localSignature {
		return "", fmt.Errorf(ErrorInvalidMsgSignature.Error()+",msgSignature = %s, localSignature = %s", msgSignature, localSignature)
	}
	plainText, err := c.Decrypt(msgSignature, timestamp, nonce, echoStr)
	if err != nil {
		return "", err
	}

	parsed, err := c.ParsePlainText(plainText)
	if nil != err {
		return "", err
	}
	receiverId := parsed.ReceiverId
	msg := parsed.Msg
	if len(c.config.ReceiverId) > 0 && strings.Compare(receiverId, c.config.ReceiverId) != 0 {
		return "", errors.New("receiver_id is not equal")
	}
	return msg, nil
}

// GetSignature 消息体的签名
func (c *Crypto) GetSignature(token, timestamp, nonce, msgEncrypt string) string {
	strs := sort.StringSlice{token, timestamp, nonce, msgEncrypt}
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

// Decrypt 用来解密的基本方法
func (c *Crypto) Decrypt(msgSignature, timestamp, nonce, base64EncryptMsg string) ([]byte, error) {
	localSignature := c.GetSignature(c.config.Token, timestamp, nonce, base64EncryptMsg)
	if msgSignature != localSignature {
		return nil, fmt.Errorf(ErrorInvalidMsgSignature.Error()+",msgSignature = %s, localSignature=%s", msgSignature, localSignature)
	}
	aesKey, err := c.GetAesKey()
	if nil != err {
		return nil, err
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

// DecryptMsg 解密消息
func (c *Crypto) DecryptMsg(msgSignature, timestamp, nonce, base64EncryptMsg string) (string, error) {
	localSignature := c.GetSignature(c.config.Token, timestamp, nonce, base64EncryptMsg)
	if msgSignature != localSignature {
		return "", fmt.Errorf(ErrorInvalidMsgSignature.Error()+",msgSignature = %s, localSignature=%s", msgSignature, localSignature)
	}
	echoStr, err := c.Decrypt(msgSignature, timestamp, nonce, base64EncryptMsg)
	if err != nil {
		return "", err
	}

	parsed, err := c.ParsePlainText(echoStr)
	if nil != err {
		return "", err
	}
	receiverId := parsed.ReceiverId
	msg := parsed.Msg
	if len(c.config.ReceiverId) > 0 && strings.Compare(string(receiverId), c.config.ReceiverId) != 0 {
		return "", errors.New("receiver_id is not eq")
	}
	return msg, nil
}

func (c *Crypto) pKCS7Padding(plainText string) []byte {
	padding := blockSize - (len(plainText) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	var buffer bytes.Buffer
	buffer.WriteString(plainText)
	buffer.Write(padText)
	return buffer.Bytes()
}

func (c *Crypto) pKCS7Unpadding(plaintext []byte) ([]byte, error) {
	plainTextLen := len(plaintext)
	if nil == plaintext || plainTextLen == 0 {
		return nil, errors.New("pKCS7Unpadding error nil or zero")
	}
	if plainTextLen%blockSize != 0 {
		return nil, errors.New("pKCS7Unpadding text not a multiple of the block size")
	}
	paddingLen := int(plaintext[plainTextLen-1])
	return plaintext[:plainTextLen-paddingLen], nil
}

func (c *Crypto) ParsePlainText(text []byte) (*WXBizParseBody, error) {
	plainText, err := c.pKCS7Unpadding(text)
	if nil != err {
		return nil, err
	}
	textLen := uint32(len(plainText))
	if textLen < 20 {
		return nil, errors.New("plain is to small 1")
	}
	random := plainText[:16]
	msgLen := binary.BigEndian.Uint32(plainText[16:20])
	if textLen < (20 + msgLen) {
		return nil, errors.New("plain is to small 2")
	}
	msg := plainText[20 : 20+msgLen]
	receiverId := plainText[20+msgLen:]
	parsed := &WXBizParseBody{
		Random:     string(random),
		MsgLen:     msgLen,
		Msg:        string(msg),
		ReceiverId: string(receiverId),
	}
	return parsed, nil
}

// Encrypt 用来加密字符串的基本方法
func (c *Crypto) Encrypt(plainText string) ([]byte, error) {
	aesKey, err := c.GetAesKey()
	if err != nil {
		return nil, err
	}
	padMsg := c.pKCS7Padding(plainText)
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	cipherText := make([]byte, len(padMsg))
	iv := aesKey[:aes.BlockSize]
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, padMsg)
	base64Msg := make([]byte, base64.StdEncoding.EncodedLen(len(cipherText)))
	base64.StdEncoding.Encode(base64Msg, cipherText)
	return base64Msg, nil
}

// EncryptMsg 加密消息
func (c *Crypto) EncryptMsg(replyMsg, timestamp, nonce string) (*WXBizJsonMsg4Send, error) {
	randStr := grand.S(16)
	var buffer bytes.Buffer
	buffer.WriteString(randStr)

	msgLenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLenBuf, uint32(len(replyMsg)))
	buffer.Write(msgLenBuf)
	buffer.WriteString(replyMsg)
	buffer.WriteString(c.config.ReceiverId)

	encodeText, err := c.Encrypt(buffer.String())
	if nil != err {
		return nil, err
	}
	cipherText := string(encodeText)
	signature := c.GetSignature(c.config.Token, timestamp, nonce, cipherText)
	send := &WXBizJsonMsg4Send{
		Signature: signature,
		Timestamp: timestamp,
		Nonce:     nonce,
		Encrypt:   cipherText,
	}
	return send, nil
}
