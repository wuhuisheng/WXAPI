package wxapi
//授权方授权信息
type APPAuthInfoResp struct {
	AuthorizationInfo struct{
		AuthorizerAppid    string        `json:"authorizer_appid"`//授权方appid
		AuthorizerAccessToken string     `json:"authorizer_access_token"`//授权方接口调用凭据
		ExpiresIn               string   `json:"expires_in"`//有效期
		AuthorizerRefreshToken string     `json:"authorizer_refresh_token"`//
		FuncInfo                []Fuc   `json:"func_info"`//授权给开发者的权限集列表

	}         `json:"authorization_info"`//授权信息
	*BaseResp
}
type Fuc struct {
	Funcscope_category string `json:"funcscope_category"` //权限集的ID
}
//授权方公众号详细信息
type APPUserInfoResp struct {
	Nick_name     string `json:"nick_name"`
	Head_img      string  `json:"head_img"`
	Service_type_info string `json:"service_type_info"`
	Verify_type_info string  `json:"verify_type_info"`
	User_name string         `json:"user_name"`
	Signature string         `json:"signature"`
	Principal_name string    `json:"principal_name"`
	Business_info  string    `json:"business_info"`
	Alias string             `json:"alias"`
	Qrcode_url string        `json:"qrcode_url"`
	Authorization_info string  `json:"authorization_info"`
	Authorization_appid string  `json:"authorization_appid"`
	Func_info           []Fuc   `json:"func_info"`
	*BaseResp

}

//第三方授权相关解析结构体
type APPAuthMsg struct {
	AppId                 string  `xml:"AppId"` //第三方平台appid
	CreateTime            string  `xml:"CreateTime"`//时间戳
	InfoType              string  `xml:"InfoType"`//unauthorized是取消授权，updateauthorized是更新授权，authorized是授权成功通知
	AuthorizerAppid       string    //公众号或小程序
	AuthorizationCode     string    //授权码，可用于换取公众号的接口调用凭据
	AuthorizationCodeExpiredTime string //授权码过期时间
	PreAuthCode                  string //预授权码
	ComponentVerifyTicket string  `xml:"ComponentVerifyTicket"`
}

type APPAuthtokenResp struct {
	component_access_token	string //第三方平台access_token
	expires_in	string //有效期
}

type APPPreResp struct {
	pre_auth_code	string //预授权码
	expires_in	string //有效期
}

type APPOptionResp struct {
	authorizer_appid	string //授权公众号或小程序的appid
	option_name	  string //选项名称
	option_value	string //选项值
}

