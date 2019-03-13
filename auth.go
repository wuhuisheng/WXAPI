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
//公众号获取用户基本信息
type UserResp struct {
	*BaseResp
	Subscribe	string//用户是否订阅该公众号标识，值为0时，代表此用户没有关注该公众号，拉取不到其余信息。
	Openid	string//用户的标识，对当前公众号唯一
	Nickname	string//用户的昵称
	Sex	int//用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	City	string//用户所在城市
	Country	string//用户所在国家
	Province	string//用户所在省份
	Language	string//用户的语言，简体中文为zh_CN
	Headimgurl	string//用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	Subscribe_time	int//用户关注时间，为时间戳。如果用户曾多次关注，则取最后关注时间
	Unionid	string//只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
	Remark	string//公众号运营者对粉丝的备注，公众号运营者可在微信公众平台用户管理界面对粉丝添加备注
	Groupid	int//用户所在的分组ID（兼容旧的用户分组接口）
	Tagid_list	[]int//用户被打上的标签ID列表
	Subscribe_scene	string//返回用户关注的渠道来源，ADD_SCENE_SEARCH 公众号搜索，ADD_SCENE_ACCOUNT_MIGRATION 公众号迁移，ADD_SCENE_PROFILE_CARD 名片分享，ADD_SCENE_QR_CODE 扫描二维码，ADD_SCENEPROFILE LINK 图文页内名称点击，ADD_SCENE_PROFILE_ITEM 图文页右上角菜单，ADD_SCENE_PAID 支付后关注，ADD_SCENE_OTHERS 其他
	Qr_scene	int//二维码扫码场景（开发者自定义）
	Qr_scene_str	string//二维码扫码场景描述（开发者自定义
}

type UserListResp struct {


	Total	string//关注该公众账号的总用户数
	Count	string// 拉取的OPENID个数，最大值为10000
	Data	struct{
		Openid  []string
	}   // 列表数据，OPENID的列表
	Next_openid	string// 拉取列表的最后一个用户的OPENID
	*BaseResp
}

type MenuResp struct {
	Button []Menu
	*BaseResp
}
type Menu struct {
    Type  string    //菜单的响应动作类型，view表示网页类型，click表示点击类型，miniprogram表示小程序类型
	Name  string    //菜单标题
	Key   string    //菜单KEY值，用于消息接口推送
	Url   string    //网页 链接，用户点击菜单可打开链接，不超过1024字节。 type为miniprogram时，不支持小程序的老版本客户端将打开本url
	Appid string    //小程序的appid
	Pagepath string //小程序的页面路径
	Media_id string //调用新增永久素材接口返回的合法media_id
	Sub_button []Menu //二级菜单数组
}
type IpResp struct {
	*BaseResp
	Ip_list []string //微信服务器IP地址列表
}

type UPloadMaterialResp struct {
	Media_id  string //新增的永久素材的media_id
	Url       string //新增的图片素材的图片URL（仅新增图片素材时会返回该字段）
	*BaseResp
}

type TicketReso struct {
	ticket   string
	expire_seconds string
	url          string
}

type QueryMaterialistlResp struct {
	Total_count   string
	Item_count    string
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