package wxapi

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"time"
)

//typ 事件类型
func (msg ReqMsg)createMsgBase(typ string)*MsgBase  {
	base := new(MsgBase)
	base.ToUserName.Value = msg.FromUserName
	base.FromUserName.Value = msg.ToUserName
	base.CreateTime.Value = fmt.Sprintf("%d", time.Now().Unix())
    base.MsgType.Value = typ
    return  base
}
// 创建被动回复文本消息
func (msg ReqMsg)CreatTextMsg(content string)TextMsg  {
	textmsg := TextMsg{}
	textmsg.Content.Value=content
	textmsg.MsgBase=msg.createMsgBase(MsgTypeText)
	return textmsg
}
// 创建被动回复图片消息
func (msg ReqMsg)CreatImageMsg(mediaID string)ImageMsg  {
	result := ImageMsg{}
	result.MediaID.Value=mediaID
	result.MsgBase=msg.createMsgBase(MsgTypeImage)
	return result
}
// 创建被动回复语音消息
func (msg ReqMsg)CreatVoiceMsg(mediaID string)VoiceMsg  {
	result := VoiceMsg{}
	result.MediaID.Value=mediaID
	result.MsgBase=msg.createMsgBase(MsgTypeVoice)
	return result
}
// 创建被动回复视频消息
func (msg ReqMsg)CreatTVideoMsg(title, descr, MediaID string)VideoMsg  {
	result := VideoMsg{}
	result.Video.Title.Value=title
	result.Video.Description.Value=descr
	result.Video.MediaID.Value=MediaID
	result.MsgBase=msg.createMsgBase(MsgTypeVideo)
	return result
}
// 创建被动回复音乐消息
func (msg ReqMsg)CreatTMusicMsg(title, descr, musicURL, hqMusicURL, thumbMediaID string)MusicMsg  {
	result := MusicMsg{}
	result.Music.Title.Value=title
	result.Music.Description.Value=descr
	result.Music.MusicURL.Value=musicURL
	result.Music.HQMusicURL.Value=hqMusicURL
	result.Music.ThumbMediaID.Value=thumbMediaID
	result.MsgBase=msg.createMsgBase(MsgTypeVideo)
	return result
}
// 创建被动回复图文消息
func (msg ReqMsg)CreatTNewsMsg(articles []*Article)newsMsg  {
	result := newsMsg{}
	result.Articles=articles
	result.ArticleCount.Value=fmt.Sprintf("%d", len(articles))
	result.MsgBase=msg.createMsgBase(MsgTypeVideo)
	return result
}

// 创建将消息转发到客服

func (msg ReqMsg)CreatTransfer2CustomerService(kfAccount ...string)Transfer2CustomerService  {
	result := Transfer2CustomerService{}

	if len(kfAccount) > 0 {
		result.KfAccount.Value=kfAccount[0]
	}
	
	result.MsgBase=msg.createMsgBase(MsgTypeVideo)
	return result
}
//加密消息体
func CreatEncryptMsg(data []byte,aesKey []byte,appid string,token string)EncryptMsg  {
	result := EncryptMsg{}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, int32(len(data)))

	length := buf.Bytes()
	random := RandAlnum(16)
	appID := []byte(appid)

	plain := bytes.Join([][]byte{random, length, data, appID}, nil)
	cipher, _ := AesEncrypt(plain, aesKey)

	result.Nonce.Value=RandNumStr(10)
	result.TimeStamp.Value=fmt.Sprintf("%d", time.Now().Unix())
    result.Encrypt.Value = Base64Encode(cipher)
    result.MsgSignature.Value=SignMsg(token,result.TimeStamp.Value,result.Nonce.Value,result.Encrypt.Value)
	return result
}

func ReplyMsgData(msg interface{}) []byte {

	data,err := xml.Marshal(msg)
	if err!=nil {
		return nil
	}
	return data
}
