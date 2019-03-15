package wxapi

//请求微信网页授权access_token 返回值
type AuthResp struct {
	Access_token  string
	Expires_in    string
	Refresh_token string
	Openid        string
	Scope         string
}
type BaseResp struct {
	Errcode string
	Errmsg  string
	*JsonResponse
}
//网页授权返回用户结构体
type AuthuserResp struct {
	Openid    string   //用户的唯一标识
	Nickname  string   //用户昵称
	Sex       int   //用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	Province  string   //用户个人资料填写的省份
	City      string   //普通用户个人资料填写的城市
	Country   string   //国家，如中国为CN
	Headimgurl string  //用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效
	Privilege []string //用户特权信息，json 数组，
	Unionid    string  //只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段
	*BaseResp
}





type TicketReso struct {
	ticket   string
	expire_seconds string
	url          string
}

type QueryMaterialistlResp struct {
	Total_count   int
	Item_count    int
	Item []Materia
	*BaseResp
}
type  Materia struct {
	title	string//图文消息的标题
	thumb_media_id	string//图文消息的封面图片素材id（必须是永久mediaID）
	show_cover_pic	string//是否显示封面，0为false，即不显示，1为true，即显示
	author	string//作者
	digest	string//图文消息的摘要，仅有单图文消息才有摘要，多图文此处为空
	content	string//图文消息的具体内容，支持HTML标签，必须少于2万字符，小于1M，且此处会去除JS
	url	string//图文页的URL，或者，当获取的列表是图片素材列表时，该字段是图片的URL
	content_source_url	string//图文消息的原文地址，即点击“阅读原文”后的URL
	update_time	string//这篇图文消息素材的最后更新时间
	name	string//文件名称
}

//素材请求todo
type ArticleReq struct {

}

// JSSDKSignature JSSDK 签名对象
type JSSDKSignature struct {
	AppID, Noncestr, Sign string
	Timestamp             int64
}