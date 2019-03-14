package wxapi

import (
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/robfig/cron"
	"io/ioutil"
	"net/http"
	"net/url"
)

type WXKfManager interface {

	//小程序或者公众号授权给第三方平台


	HandleCompentAuthEventPush(context * gin.Context,responsehandler ...func(appmsg APPAuthMsg))              //授权事件推送,根据信息更新授权状态，对应授权事件接收URL
	HanleCompentAuth(context * gin.Context, responsehandler func(authinfo APPAuthInfoResp)(redicturl string)) //授权回调并得到授权方的账号信息
	//授权方信息管理
	GetCompentAuthOptionInfo(authorizer_appid,option_name string,responsehandler ...func(APPOptionResp))     //获取授权方的选项设置信息
	SetCompentAuthOption(authorizer_appid,option_name,option_value string,responsehandler ...func(BaseResp)) //设置授权方的选项设置信息
	//代公众号实现业务API
	//对应到消息与事件接收url中
	HandleAppEventPush(ctx * gin.Context, handler func(msg ReqMsg)(usedefult bool,replymsg interface{})) //接收处处公众号事件和消息
	//用户网页登录
	//对应用户授权后回调
	HanledAppAuth(context * gin.Context,completeHandler func(resp AuthResp,authuser AuthuserResp,state string)(redicturl string)) //获取第三方用户信息，并跳转到登录后的页面
	GetAppAuthurl(appid,scope,redirect_uri,state string) string                                                                   //用户登录发起页面
	//用户信息


}

type WXKFManager struct {
	Compenttoken    string //第三方平台的token 消息校验token
	CompentAppid    string //第三方平台的APPid
	CompentAeskey   string //第三方平台的秘钥
	Componentsecret string //第三方平台的appsecret
	Component_access_token string //第三方平台的component_access_token
	ComponentVerifyTicket string //微信服务器传输的第三方平台ComponentVerifyTicket
	Pre_auth_code     string  //预授权码
	Redircturl        string  //公众号授权完成后的回调URL用来接收授权码auth_code
	Compentauthurl    string  //当前第三方平台授权移动端连接
	AppAuthinfos     []APPAuthInfoResp //第三方公众号令牌数组，自动刷新
}

//初始化开放平台管理器
func InitWXKFManager(token,appid,EncodingAESKey,appseceret,Redircturl string,refreshAppAuthHanlder func(appauth ...APPAuthInfoResp),AppAuthinfos...[]APPAuthInfoResp) *WXKFManager {
      var  wx WXKFManager
      wx.Compenttoken = token
      wx.CompentAeskey = EncodingAESKey
      wx.CompentAppid = appid
      wx.Componentsecret = appseceret
      wx.Redircturl=Redircturl

	//每小时刷新一次component_accesstoken
	cr := cron.New()
	spec := "0 0 */1 * * ?"
	err:=cr.AddFunc(spec,wx.getComponent_access_token)

	err= cr.AddFunc(spec, func() {
		if len(AppAuthinfos)>0 {
			wx.AppAuthinfos=AppAuthinfos[0]
			for _,value := range wx.AppAuthinfos{

				wx.RefreshCompentAuthAccessToken(value.AuthorizationInfo.AuthorizerAppid,value.AuthorizationInfo.AuthorizerRefreshToken, func(APPAuthInfo APPAuthInfoResp) {
					refreshAppAuthHanlder(APPAuthInfo)
				})
			}
		}else {
			refreshAppAuthHanlder()
		}
	})
	fmt.Println("开房平台定时任务初始化",err,spec)

    spe :="0 */10 * * * ?"
    err=cr.AddFunc(spe,wx.getCompent_pre_auth_code)
	fmt.Println("开房平台定时任务初始化",err,spe)
	cr.Start()

      return  &wx
}
//授权流程

//授权回调 获取授权码并根据授权码获取授权信息
func (wx *WXKFManager) HanleCompentAuth(context * gin.Context, responsehandler func(authinfo APPAuthInfoResp)(redicturl string)){



	authcode := context.Query("auth_code")
	//查询公众号授权第三方平台的权限
	wx.getCompentAuthAccesstoken(authcode, func(appAuthInfo APPAuthInfoResp) {
		wx.AppAuthinfos=append(wx.AppAuthinfos, appAuthInfo)
		redicturl := responsehandler(appAuthInfo)
		context.Redirect(http.StatusMovedPermanently,redicturl)
	})
}

//获取授权方公众号账号基本信息
func (wx *WXKFManager) GetCompentAuthorizerInfo(authorizer_appid string,response ...func(resp APPUserInfoResp)) {
	url :="https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token="
	url =url +wx.Component_access_token
	POSTJson(url,gin.H{"component_appid": wx.CompentAppid,"authorizer_appid":authorizer_appid}, func(resp JsonResponse) {
		var result APPUserInfoResp
		mapstructure.Decode(resp.Dic,&result)
		result.JsonResponse=&resp
		if len(response)>0 {
			response[0](result)
		}
	})

}

//获取授权方选项设置信息

func (wx *WXKFManager) GetCompentAuthOptionInfo(authorizer_appid,option_name string,responsehandler ...func(APPOptionResp))  {

	url :="https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token="
	url =url +wx.Component_access_token
	 POSTJson(url,gin.H{"component_appid": wx.CompentAppid,"authorizer_appid":authorizer_appid,"option_name":option_name}, func(resp JsonResponse) {
		 var result APPOptionResp
		 mapstructure.Decode(resp.Dic,&result)
		 result.JsonResponse=&resp
		 if len(responsehandler)>0 {
			 responsehandler[0](result)
		 }
	 })



}
//设置授权方选项信息

func (wx *WXKFManager) SetCompentAuthOption(authorizer_appid,option_name,option_value string,responsehandler ...func(BaseResp))  {
	url :="https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token="
	url =url +wx.Component_access_token
	POSTJson(url,gin.H{"component_appid": wx.CompentAppid,"authorizer_appid":authorizer_appid,"option_name":option_name,
		"option_value":option_value}, func(resp JsonResponse) {
		var result BaseResp
		mapstructure.Decode(resp.Dic,&result)
		result.JsonResponse=&resp
		if len(responsehandler)>0 {
			responsehandler[0](result)
		}
	})

}

//授权通知处理
func (wx *WXKFManager) HandleCompentAuthEventPush(context * gin.Context,responsehandler ...func(appmsg APPAuthMsg)){

	  wx.parsereqToAPPAuthMsg(context, func(CheckSign bool, Orignmsg ReqMsg, Decrptmsg APPAuthMsg, safe bool) {
		  if CheckSign {
			  if len(responsehandler)>0 {
				  responsehandler[0](Decrptmsg)
			  }
		  }
	  })

	  context.String(http.StatusOK,"success")

}



//代公众号实现网页授权


/****************代公众号实现业务*******************/
//1.代公众号调用接口
//获取用户信息
func (wx *WXKFManager)GetUserInfo( authorizer_access_token,authorizer_appid string ,hanlder ...func(user UserResp) ){

	if len(hanlder)>0 {
		getuserInfo(authorizer_access_token,authorizer_appid,hanlder[0])
	}else {
		getuserInfo(authorizer_access_token,authorizer_appid)
	}
}
//获取用户列表
func (wx *WXKFManager)GetUserList( authorizer_access_token string,hanlder func(user UserListResp),nestopid ...string ){

	if len(nestopid)>0 {
		getuserlist(authorizer_access_token,hanlder,nestopid[0] )
	}else {
		getuserlist(authorizer_access_token,hanlder)
	}
}

//2.代公众号处理消息和事件
func (wx *WXKFManager) HandleAppEventPush(ctx * gin.Context, handler func(msg ReqMsg)(usedefult bool,replymsg interface{})){
	wx.parsereqToReqMsg(ctx, func(CheckSign bool, Orignmsg ReqMsg, Decrptmsg ReqMsg, safe bool) {
		if !CheckSign {
			ctx.String(http.StatusForbidden, "验证签名错误")
			return
		}
		if ctx.Request.Method!= http.MethodPost {
			echostr := ctx.Query("echostr")
			ctx.String(http.StatusOK,echostr)
			return
		}
		def,replymsg := handler(Decrptmsg)
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

//3.代公众号发起网页授权
//获取代公众号发起网页授权url
func (wx *WXKFManager) GetAppAuthurl(appid,scope,redirect_uri,state string) string {
	str := url.QueryEscape(redirect_uri)
	url:="https://open.weixin.qq.com/connect/oauth2/authorize?appid="
	url=url+appid+"&redirect_uri="+str
	url =url+"&response_type=code&scope="+scope+"&state="+state+"&component_appid="+wx.CompentAppid
	url =url+"#wechat_redirect"

	return url
}
//网页授权后回调
func (wx *WXKFManager) HanledAppAuth(context * gin.Context,completeHandler func(resp AuthResp,authuser AuthuserResp,state string)(redicturl string))  {

	code := context.Query("code")
	appid := context.Query("appid")
	state := context.Query("state")
	fmt.Println(code,appid,"接收到的微信信息时")
	wx.getAppAuthUserAccesstoken(code,appid, func(authresp AuthResp) {
		wx.getAppAuthuserInfo(authresp, func(authuser AuthuserResp) {
			redicturl:=completeHandler(authresp,authuser,state)
			context.Redirect(http.StatusMovedPermanently,redicturl)
		})
	})
}


//4代公众号调用jssdk


/*...........私有方法.................*/

//加密XML结构体消息体
func (wx *WXKFManager)msgEncrept(msg interface{})EncryptMsg  {
	return  CreatEncryptMsg(ReplyMsgData(msg),DecodeAESKey(wx.CompentAeskey),wx.CompentAppid,wx.Compenttoken)
}

/*授权流程使用*/
//1.获取第三方平台access_token
func (wx *WXKFManager) getComponent_access_token() {
	if len(wx.ComponentVerifyTicket)>0 {
		wx.Component_access_token = getcomponent_token(wx.CompentAppid,wx.Componentsecret,wx.ComponentVerifyTicket)

	}else {
		wx.Component_access_token=""
	}

}
//2.获取预授权码
func (wx *WXKFManager) getCompent_pre_auth_code()  {
	if len(wx.Component_access_token)>0 {
		wx.Pre_auth_code= getCompent_pre_authcode(wx.Component_access_token,gin.H{"component_appid": wx.CompentAppid})
		wx.Compentauthurl= getCompentAuthUrl(wx.CompentAppid,wx.Pre_auth_code,wx.Redircturl)
	}else {
		wx.Pre_auth_code=""
	}
}
//3.使用授权码换取公众号或小程序的接口调用凭据和授权信息
func (wx *WXKFManager) getCompentAuthAccesstoken(authorization_code string,handler ...func(appAuthInfo APPAuthInfoResp)) {

	url :="https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token="
	url =url +wx.Component_access_token
	POSTJson(url,gin.H{"component_appid": wx.CompentAppid,"authorization_code":authorization_code}, func(response JsonResponse) {
		var result APPAuthInfoResp
		mapstructure.Decode(response.Dic,&result)
		result.JsonResponse=&response
		if len(handler)>0 {
			handler[0](result)
		}
	})


}

//4.（刷新）授权公众号或小程序的接口调用凭据
func (wx *WXKFManager) RefreshCompentAuthAccessToken(authorizer_appid,authorizer_refresh_token string,response ...func(APPAuthInfo APPAuthInfoResp))  {
	url :="https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token="
	url =url +wx.Component_access_token
	POSTJson(url,gin.H{"component_appid": wx.CompentAppid,"authorizer_appid":authorizer_appid,"authorizer_refresh_token":authorizer_refresh_token}, func(res JsonResponse) {
		var result APPAuthInfoResp
		mapstructure.Decode(res.Dic,&result)
		result.JsonResponse=&res
		if len(response)>0 {
			response[0](result)
		}
	})

}

/*代公众号实现业务使用*/
//代公众号获取useraccess_token
func (wx *WXKFManager) getAppAuthUserAccesstoken(code ,appid string,handler ...func(authresp AuthResp)) {
	url :="https://api.weixin.qq.com/sns/oauth2/component/access_token?appid="+appid
	url = url +"&code="+code+"&grant_type=authorization_code&component_appid="+wx.CompentAppid+"&component_access_token="
	url = url + wx.Component_access_token
	resp := Get(url)

	var result AuthResp
	mapstructure.Decode(resp,&result)
	fmt.Println("代公众号获取usertoken",resp,result,appid,url)
	if len(handler)>0 {
		handler[0](result)
	}
}

//代公众号获取网页登录用户信息
func (wx *WXKFManager) getAppAuthuserInfo(auth AuthResp,handler ...func(authuser AuthuserResp))  {
	url :="https://api.weixin.qq.com/sns/userinfo?access_token="
	url =url +auth.Access_token+"&openid="+auth.Openid+"&lang=zh_CN"
	resp := Get(url)
	var result AuthuserResp
	mapstructure.Decode(resp,&result)
	fmt.Println("代公众号获取网页用户信息",resp,result)
	if len(handler)>0 {
		handler[0](result)
	}
}

/*解析回电xml使用*/
//解析推送授权事件的XML
func (wx *WXKFManager) parsereqToAPPAuthMsg(context2 * gin.Context,f func(CheckSign bool,Orignmsg ReqMsg,Decrptmsg APPAuthMsg,safe bool)){

	parsereqMsg(context2,wx.Compenttoken,wx.CompentAeskey, func(CheckSign bool, Orignmsg ReqMsg, safe bool) {

		var  decreptMsg APPAuthMsg
		if safe{
			decreptMsg,_ = decryptAPPAuthMsg(Orignmsg.Encrypt,wx.CompentAeskey)
			wx.ComponentVerifyTicket = decreptMsg.ComponentVerifyTicket
			if wx.Component_access_token=="" {
				wx.getComponent_access_token()
				wx.getCompent_pre_auth_code()
			}
		}
		f(CheckSign,Orignmsg,decreptMsg,safe)

	})


}
//解析代公众号实现事件消息XML
func (wx *WXKFManager) parsereqToReqMsg(context2 * gin.Context,f func(CheckSign bool,Orignmsg ReqMsg,Decrptmsg ReqMsg,safe bool)){


	parsereqMsg(context2,wx.Compenttoken,wx.CompentAeskey, func(CheckSign bool, Orignmsg ReqMsg, safe bool) {

		var  decreptMsg ReqMsg
		if safe{
			decreptMsg,_ = decryptReqmsg(Orignmsg.Encrypt,wx.CompentAeskey)
		}
		f(CheckSign,Orignmsg,decreptMsg,safe)
	})
}

//解析请求
func parsereqMsg(context2 * gin.Context,token string,aeskey string,handler func(CheckSign bool,Orignmsg ReqMsg,safe bool))  {

	sign := context2.Query("signature")
	timestamp := context2.Query("timestamp")
	nonce := context2.Query("nonce")
	encrypt_type := context2.Query("encrypt_type")
	msgsign := context2.Query("msg_signature")

	var  event ReqMsg
	s,_:=ioutil.ReadAll(context2.Request.Body)
	xml.Unmarshal(s,&event)

	safe := encrypt_type=="aes"
	checksign :=false

	if safe{
		checksign = SignMsg(token,timestamp,nonce,event.Encrypt)==msgsign
	}else {
		checksign =SignMsg(token,timestamp,nonce)==sign
	}
	handler(checksign,event,safe)
}

