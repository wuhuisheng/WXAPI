package wxapi

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/robfig/cron"
	"net/http"
	"net/url"
)

type WXmanager interface {

	//获取access_token
	//InitWXManager(appid,appseceret,token,EncodingAESKey string)WXManager
	CheckSignature(c*gin.Context)bool    //签名验证
	//网页授权
	WebAuth(c * gin.Context,f func(auth AuthResp,  userinfo AuthuserResp,state string)(redicturl string))
	//事件消息处理
	HandleEventPush(ctx * gin.Context, f func(msg ReqMsg)(Isdefreply bool ,replymsg interface{}))
    //网页授权发起页连接
    GetAuthMenuurl(redirect_uri,scope,state string)string
    //自定义菜单管理
	CreatMenu(param gin.H,responsehanler func(resp BaseResp))
	QuerytMenu(responsehanler func(resp BaseResp))
    //用户管理
	GetuserInfo(openid string,respomsehandler func(resp UserResp))//用户详细信息
	GetUserlist(respomsehandler func(resp UserListResp),nestopenid ...string)//用户列表
	//账号管理  带参数二维码生成

	//微信服务器IP列表
	GetWXIPlist( respomsehandler func(iplist IpResp))
	//素材管理
	CreatNews(param gin.H)UPloadMaterialResp
	CreatMaterial(param gin.H,typ string)UPloadMaterialResp
	QuerymaterialList(param gin.H)QueryMaterialistlResp

}

//实现了上述接口
type WXManager struct {
	Appid          string
	Appsecret      string
	token          string
	Accesstoken    string
	EncodingAESKey string

}
//初始化微信公众平台管理器,返回引用类型，确保内外变化一致
func InitWXManager(appid,appseceret,token,EncodingAESKey string)*WXManager {

	var wx WXManager
	wx.Appid=appid
	wx.token=token
	wx.Appsecret=appseceret
	wx.EncodingAESKey=EncodingAESKey
    wx.getAceeesToken()
	//每小时刷新一次token
	c := cron.New()
	spec := "0 0 */1 * *  ?"
	err :=c.AddFunc(spec,wx.getAceeesToken)
	c.Start()
    fmt.Println("初始化公众平台管理器定时任务",err,spec)
	return &wx
}

//用户网页授权后
func (wx *WXManager) HandleAuth(c * gin.Context,handler func(auth AuthResp,authorization_code ,state string)(redicturl string))  {
	code := c.Query("code")
	state := c.Query("state")
	wx.GetAuthAccesstoken(code, func(resp AuthResp) {
		c.Redirect(http.StatusMovedPermanently,handler(resp,code,state))
	})
}


//处理自定义菜单等事件replymsg 回复的消息结构体 defreply 1为默认回复，0为自定义回复
func (wx *WXManager)HandleEventPush(ctx * gin.Context, f func(msg ReqMsg)(Isdefreply bool ,replymsg interface{})){


	wx.ParseReq(ctx, func(CheckSign bool, Orignmsg ReqMsg, Decrptmsg ReqMsg, safe bool) {

		if !CheckSign {
			ctx.String(http.StatusForbidden, "验证签名错误")
			return
		}
		if ctx.Request.Method!= http.MethodPost {
			echostr := ctx.Query("echostr")
			ctx.String(http.StatusOK,echostr)
			return
		}
		def,replymsg := f(Decrptmsg)
		if def {
			ctx.String(http.StatusOK,"success")
			return
		}
		if safe {
            //
			ctx.String(http.StatusOK,string(ReplyMsgData(wx.msgEncrept(replymsg))))
		}else {
			ctx.String(http.StatusOK, string(ReplyMsgData(replymsg)))
		}

	})
}

func (wx *WXManager) msgEncrept(msg interface{})EncryptMsg  {
   return  CreatEncryptMsg(ReplyMsgData(msg),DecodeAESKey(wx.EncodingAESKey),wx.Appid,wx.token)
}


//公众号网页授权使用
//redirect_uri 直接写授权后的重定向URL 不需要URLencode
//scope  snsapi_base,snsapi_userinfo
//redirect_uri 授权后重定向的回调链接地址
func (wx *WXManager)GetAuthMenuurl(redirect_uri,scope,state string)string  {
	str := url.QueryEscape(redirect_uri) //urlencode
	result :="https://open.weixin.qq.com/connect/oauth2/authorize?appid="
	result = result+wx.Appid+"&redirect_uri="+str+"&response_type=code&scope="+scope+"&state="+state
	result = result+"#wechat_redirect"
	return result
}

//code 用户同意授权后，页面将跳转至 redirect_uri/?code=CODE&state=STATE 微信发来的
//通过code换取网页授权access_token
func (wx *WXManager) GetAuthAccesstoken(authorization_code string,hander ...func(AuthResp))  {
	usr := "https://api.weixin.qq.com/sns/oauth2/access_token?appid="
	usr =usr+wx.Appid+"&secret="+wx.Appsecret+"&code="+authorization_code+"&grant_type=authorization_code"
	fmt.Println("获取用户授权accesstoken",usr)
	resp := Get(usr)
	var result AuthResp
	mapstructure.Decode(resp,&result)
	if len(hander)>0 {
		hander[0](result)
	}
}
//刷新token
func (wx *WXManager) RefreshAuthAccesstoken(refresh string,hander ...func(AuthResp))  {

	usr :="https://api.weixin.qq.com/sns/oauth2/refresh_token?appid="+wx.Appid
	usr =usr+"&grant_type=refresh_token&refresh_token="+refresh
	resp := Get(usr)
	var respauth AuthResp
	mapstructure.Decode(resp,&respauth)
	if len(hander)>0 {
		hander[0](respauth)
	}
}
func (wx *WXManager) CheckAuthAcesstoken(access_token,openid string)  {
	resp:=Get("https://api.weixin.qq.com/sns/auth?access_token="+access_token+"&openid="+openid)
	fmt.Println(resp)
}

func (wx *WXManager)GetAuthuserInfoBycode(authorization_code string,hander ...func(AuthuserResp))  {
	
	wx.GetAuthAccesstoken(authorization_code, func(resp AuthResp) {
		wx.GetAuthuserInfo(resp.Access_token,resp.Openid, func(resp AuthuserResp) {
			if len(hander)>0 {
				hander[0](resp)
			}
		})
	})
	
}

//传入参数为上边参数返回值
//网页授权拉取用户信息
func (wx *WXManager)GetAuthuserInfo(access_token ,openid string,hander ...func(AuthuserResp))  {
	url :="https://api.weixin.qq.com/sns/userinfo?access_token="
	url =url +access_token+"&openid="+openid+"&lang=zh_CN"

	var result AuthuserResp
	req :=Get(url)
	fmt.Println("用户授权后原始信息",req)
	mapstructure.Decode(req,&result)
	if len(hander)>0 {
		hander[0](result)
	}
}
//用户管理
//获取用户信息
func (wx *WXManager)GetuserInfo(openid string,respomsehandler ...func(resp UserResp))  {
	if len(respomsehandler)>0 {
		getuserInfo(wx.Accesstoken,openid,respomsehandler[0])
	}else {
		getuserInfo(wx.Accesstoken,openid)
	}

}
//获取用户列表
func (wx *WXManager)GetUserlist(respomsehandler func(resp UserListResp),nestopenid ...string)  {
	if len(nestopenid)>0 {
		getuserlist(wx.Accesstoken,respomsehandler,nestopenid[0])
	}else {
		getuserlist(wx.Accesstoken,respomsehandler)
	}

}
//todo黑名单管理



//数据统计接口


//菜单管理
//创建自定义菜单
func (wx *WXManager)CreatMenu(param gin.H,responsehanler func(resp BaseResp)) {
	POSTJson(" https://api.weixin.qq.com/cgi-bin/menu/create?access_token="+wx.Accesstoken, param, func(response JsonResponse) {
		var result BaseResp
		json.Unmarshal(response.Data,&result)
		result.JsonResponse =&response
		responsehanler(result)
	})


}
//自定义菜单查询
func (wx *WXManager)QuerytMenu(responsehanler func(resp BaseResp)) {
	resp  := Get(" https://api.weixin.qq.com/cgi-bin/menu/create?access_token="+wx.Accesstoken)
	var result BaseResp
	mapstructure.Decode(resp,&result)
	responsehanler(result)
}
//创建个性化菜单


//获取微信服务器IP地址

func (wx *WXManager)GetWXIPlist( respomsehandler func(iplist IpResp))  {

	resp :=Get("https://api.weixin.qq.com/cgi-bin/getcallbackip?access_token="+wx.Accesstoken)
	var result IpResp
	mapstructure.Decode(resp,&result)
	respomsehandler(result)
}

//素材管理
//添加图文素材
func (wx *WXManager)CreatNews(param gin.H, respomsehandler func(resp UPloadMaterialResp))  {

	POSTJson("https://api.weixin.qq.com/cgi-bin/material/add_news?access_token="+wx.Accesstoken,param, func(response JsonResponse) {
		var result UPloadMaterialResp
		json.Unmarshal(response.Data,&result)
		result.JsonResponse =&response
		respomsehandler(result)
	})

}
//添加其他素材
func (wx *WXManager)CreatMaterial(param gin.H,typ string,respomsehandler func(resp UPloadMaterialResp))  {
	POSTJson("https://api.weixin.qq.com/cgi-bin/material/add_material?access_token="+wx.Accesstoken+"&type="+typ,param, func(response JsonResponse) {
		var result UPloadMaterialResp
		json.Unmarshal(response.Data,&result)
		result.JsonResponse =&response
		respomsehandler(result)
	})

}

//查询素材列表

func (wx *WXManager)QuerymaterialList(param gin.H,respomsehandler func(resopne JsonResponse))  {
    url :="https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token="+wx.Accesstoken
   POSTJson(url,param, func(response JsonResponse) {
	   respomsehandler(response)
   })
}


//解析请求消息
func (wx *WXManager) ParseReq(context2 * gin.Context,f func(CheckSign bool,Orignmsg ReqMsg,Decrptmsg ReqMsg,safe bool)){

	parsereqMsg(context2,wx.token,wx.EncodingAESKey, func(CheckSign bool, Orignmsg ReqMsg, safe bool) {

		decreptMsg := Orignmsg
		if safe{
			decreptMsg,_ = decryptReqmsg(Orignmsg.Encrypt,wx.EncodingAESKey)
		}
		f(CheckSign,Orignmsg,decreptMsg,safe)
	})

}



//获取accesstoken
func (wx *WXManager)getAceeesToken()  {
	url :="https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid="
	url = url + wx.Appid+"&secret="+wx.Appsecret
	rep:=Get(url)
	wx.Accesstoken =fmt.Sprint(rep["access_token"])
	fmt.Println(rep)

}

//获取创建临时二维码ticket

func getTempticket(access_token ,expire_seconds,action_name,scene_str string)  {

	url :="https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token="
	url =url +access_token

	POSTJson(url,gin.H{"expire_seconds": expire_seconds,"action_name":action_name,"action_info":gin.H{"scene": gin.H{"scene_str": scene_str}}}, func(response JsonResponse) {

	})

}
