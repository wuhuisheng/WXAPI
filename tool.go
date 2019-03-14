package wxapi

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"net/http"
	"strconv"
	 crand "crypto/rand"
	"math/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

const (
	alnum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)
var (
	stdRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)
// AES解密
func AesDecrypt(ciphertext []byte, aesKey []byte) (plaintext []byte, err error) {
	if len(ciphertext)%len(aesKey) != 0 {
		err = errors.New("ciphertext is not a multiple of the block size")
		return
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return
	}
	iv := make([]byte, aes.BlockSize)
	if _, err = io.ReadFull(crand.Reader, iv); err != nil {
		return
	}
	plaintext = make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	return
}



// AES加密
func AesEncrypt(plaintext []byte, aesKey []byte) (ciphertext []byte, err error) {
	if len(plaintext)%len(aesKey) != 0 {
		plaintext = pkcs7Pad(plaintext, len(aesKey))
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return
	}
	iv := make([]byte, aes.BlockSize)
	if _, err = io.ReadFull(crand.Reader, iv); err != nil {
		panic(err)
	}
	ciphertext = make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return
}

// PKCS#7填充
func pkcs7Pad(plaintext []byte, blockSize int) []byte {
	size := blockSize - len(plaintext)%blockSize
	pads := bytes.Repeat([]byte{byte(size)}, size)
	return append(plaintext, pads...)
}
// AESKey解码
func DecodeAESKey(encodingAESKey string) []byte {
	return Base64Decode(encodingAESKey + "=")
}
// Base64解码
func Base64Decode(str string) []byte {
	data, _ := base64.StdEncoding.DecodeString(str)
	return data
}
// Base64编码
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// 对消息签名
func SignMsg(token, timestamp, nonce string, encrypt ...string) string {
	ss := sort.StringSlice{
		token,
		timestamp,
		nonce,
	}
	if len(encrypt)>0 {
		ss = append(ss,encrypt[0])
	}
	ss.Sort()

	s := sha1.New()
	io.WriteString(s,strings.Join(ss,""))
	sgin :=fmt.Sprintf("%x",s.Sum(nil))
	return sgin

}

func Get(url string) map[string]interface{} {
	client := http.Client{}
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	resp,err := client.Get(url)
	var re map[string]interface{}
	if err !=nil{
		fmt.Println(err,url)
		return map[string]interface{}{}
	}
	json.NewDecoder(resp.Body).Decode(&re)

	return re
}




//发送post,数据json请求
func POSTJson(url string,params map[string]interface{},completHandler ...func(response JsonResponse))  {


	ll,_:= json.Marshal(params)
	Post(url,ll, func(response JsonResponse) {
		if len(completHandler)>0 {
			completHandler[0](response)
		}
	})

}

func Post(url string,data[]byte,completHandler func(response JsonResponse))  {
	client := http.Client{}
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	resp,err :=client.Post(url,"application/json",bytes.NewBuffer(data))
	respdata,_:=ioutil.ReadAll(resp.Body)
	var dic map[string]interface{}
	json.Unmarshal(respdata,&dic)
	completHandler(JsonResponse{respdata,err,dic})

}

type JsonResponse struct {
	 Data  []byte
	 Err   error
	 Dic   map[string]interface{}
}


func RandNumStr(n int) string {
	return string(RandNum(n))
}
func RandNum(n int) []byte {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = strconv.Itoa(stdRand.Intn(10))[0]
	}
	return b
}
func RandAlnumStr(n int) string {
	return string(RandAlnum(n))
}

func RandAlnum(n int) []byte {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = alnum[stdRand.Intn(len(alnum))]
	}
	return b
}

/**解密密文到对应XML结构体*/
//Encrypt 密文
//req 解密后对应的xml结构体
func decryptAPPAuthMsg(Encrypt ,Aeskey string)(APPAuthMsg,error)  {

	aeskey := DecodeAESKey(Aeskey)
	var req APPAuthMsg
	data,err :=  AesDecrypt(Base64Decode(Encrypt),aeskey)
	if err != nil {

		return req,err
	}
	var msgLen int32
	err = binary.Read(bytes.NewReader(data[16:20]), binary.BigEndian, &msgLen)
	if err != nil {
		return req,err
	}

	err = xml.Unmarshal(data[20:20+msgLen], &req)
	return req,err
}
//将密文解析为ReqMsg结构体
func decryptReqmsg(Encrypt,Aeskey string)(ReqMsg,error)  {
	aeskey := DecodeAESKey(Aeskey)
	var req ReqMsg
	data,err :=  AesDecrypt(Base64Decode(Encrypt),aeskey)
	if err != nil {
		return req,err
	}
	var msgLen int32
	err = binary.Read(bytes.NewReader(data[16:20]), binary.BigEndian, &msgLen)
	if err != nil {
		return req,err
	}

	err = xml.Unmarshal(data[20:20+msgLen], &req)
	return req,err

}

/**公众号用户管理*/
//获取用户信息列表
func getuserlist(access_token string,respomsehandler func(resp UserListResp),nestopenid ...string)  {
	url :="https://api.weixin.qq.com/cgi-bin/user/get?access_token="+access_token

	if len(nestopenid)>0 {
		url=url+"&next_openid="+nestopenid[0]
	}
	var result UserListResp

	mapstructure.Decode(Get(url),&result)
	respomsehandler(result)
}
//获取用户信息
func getuserInfo(access_token,openid string,respomsehandler ...func(resp UserResp))  {
	url :="https://api.weixin.qq.com/cgi-bin/user/info?access_token="
	url =url +access_token+"&openid="+openid+"&lang=zh_CN"
	var result UserResp
	req :=Get(url)
	mapstructure.Decode(req,&result)
	if len(respomsehandler)>0 {
		respomsehandler[0](result)
	}
}


/**第三方平台授权使用*/
//获取第三方平台access_token 令牌
func getcomponent_token(copentAppid,compentAppsecret,compentticket string) string {
	param:=gin.H{"component_appid": copentAppid, "component_appsecret": compentAppsecret, "component_verify_ticket": compentticket}
     result := ""
	 POSTJson("https://api.weixin.qq.com/cgi-bin/component/api_component_token",param, func(response JsonResponse) {
		 result=fmt.Sprint(response.Dic["component_access_token"])
	 })
	return result
}
//获取第三方平台预授权码
func getCompent_pre_authcode(token string,param map[string]interface{}) string {

	url := "https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?component_access_token="
	url = url+token
	result := ""
	 POSTJson(url,param, func(response JsonResponse) {
		 result=fmt.Sprint(response.Dic["pre_auth_code"])
	 })

	return result
}
//获取第三方平台移动端授权连接
func getCompentAuthUrl(appid,pre_auth_code,redicturi string)string {

	result := "https://mp.weixin.qq.com/safe/bindcomponent?action=bindcomponent&auth_type=3&no_scan=1&component_appid="
	result = result+appid
	result = result+"&pre_auth_code="+pre_auth_code
	result = result+"&redirect_uri="+redicturi
	result = result+"&auth_type=1#wechat_redirect"
	fmt.Println("第三方平台移动端授权连接是",result)
	return result
}