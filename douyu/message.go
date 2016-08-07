package douyu

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

// 消息体类型
const (
	TypeMessageToServer   = 689
	TypeMessageFromServer = 690
)

// 弹幕服务器端相应消息类型
const (
	TypeLoginRes = "loginres" // 登录响应消息
	// 服务端返回登陆响应消息,完整的数据部分应包含的字段如下:
	// 字段说明
	//     type             表示为“登出”消息,固定为 loginres
	//     userid           用户 ID
	//     roomgroup        房间权限组
	//     pg               平台权限组
	//     sessionid        会话 ID
	//     username         用户名
	//     nickname         用户昵称
	//     is_signed        是否已在房间签到
	//     signed_count     日总签到次数
	//     live_stat        直播状态
	//     npv              是否需要手机验证
	//     best_dlev        最高酬勤等级
	//     cur_lev          酬勤等级
	TypeKeeplive = "keeplive" // 服务端心跳消息
	// 服务端响应客户端心跳的消息,完整的数据部分应包含的字段如下:
	// 字段说明
	//     type             表示为“心跳”消息,固定为 keeplive
	//     tick             对应客户端心跳请求中的 tick
	TypeChatMsg = "chatmsg" // 弹幕消息
	// 用户在房间发送弹幕时,服务端发此消息给客户端,完整的数据部分应包含的字 段如下:
	// 字段说明
	//     type             表示为“弹幕”消息,固定为 chatmsg
	//     gid              弹幕组 id
	//     rid              房间 id
	//     uid              发送者 uid
	//     nn               发送者昵称
	//     txt              弹幕文本内容
	//     cid              弹幕唯一 ID
	//     level            用户等级
	//     gt               礼物头衔:默认值 0(表示没有头衔)
	//     col              颜色:默认值 0(表示默认颜色弹幕)
	//     ct               客户端类型:默认值 0(表示 web 用户)
	//     rg               房间权限组:默认值 1(表示普通权限用户)
	//     pg               平台权限组:默认值 1(表示普通权限用户)
	//     dlv              酬勤等级:默认值 0(表示没有酬勤)
	//     dc               酬勤数量:默认值 0(表示没有酬勤数量)
	//     bdlv             最高酬勤等级:默认值 0(表示全站都没有酬勤)
	TypeLuckyGuy = "onlinegift" //领取在线鱼丸暴击消息
	// 在线领取鱼丸时,若出现暴击(鱼丸数大于等于 60)服务则发送领取暴击消息 到客户端。完整的数据部分应包含的字段如下:
	// 字段说明
	//     type             表示为“领取在线鱼丸”消息,固定为 onlinegift
	//     rid              房间 ID
	//     uid              用户 ID
	//     gid              弹幕分组 ID
	//     sil              鱼丸数
	//     if               领取鱼丸的等级
	//     ct               客户端类型标识
	//     nn               用户昵称
	TypeNewGift = "dgb" // 赠送礼物消息
	// 用户在房间赠送礼物时,服务端发送此消息给客户端。完整的数据部分应包含的 字段如下:
	// 字段说明
	//     type   表示为“赠送礼物”消息,固定为 dgb
	//     rid    房间 ID
	//     gid    弹幕分组 ID
	//     gfid   礼物 id
	//     gs     礼物显示样式
	//     uid    用户 id
	//     nn     用户昵称
	//     str    用户战斗力
	//     level  用户等级
	//     dw     主播体重
	//     gfcnt  礼物个数:默认值 1(表示 1 个礼物)
	//     hits   礼物连击次数:默认值 1(表示 1 连击)
	//     dlv    酬勤头衔:默认值 0(表示没有酬勤)
	//     dc     酬勤个数:默认值 0(表示没有酬勤数量)
	//     bdl    全站最高酬勤等级:默认值 0(表示全站都没有酬勤)
	//     rg     房间身份组:默认值 1(表示普通权限用户)
	//     pg     平台身份组:默认值 1(表示普通权限用户)
	//     rpid   红包 id:默认值 0(表示没有红包)
	//     slt    红包开启剩余时间:默认值 0(表示没有红包)
	//     elt    红包销毁剩余时间:默认值 0(表示没有红包)
	TypeUserEnter = "uenter" // 特殊用户进房通知消息
	// 具有特殊属性的用户进入直播间时,服务端发送此消息至客户端。完整的数据部 分应包含的字段如下:
	// 字段说明
	//     type   表示为“用户进房通知”消息,固定为 uenter
	//     rid    房间 ID
	//     gid    弹幕分组 ID
	//     uid    用户 ID
	//     nn     用户昵称
	//     str    战斗力
	//     level  新用户等级
	//     gt     礼物头衔:默认值 0(表示没有头衔)
	//     rg     房间权限组:默认值 1(表示普通权限用户)
	//     pg     平台身份组:默认值 1(表示普通权限用户)
	//     dlv    酬勤等级:默认值 0(表示没有酬勤)
	//     dc     酬勤数量:默认值 0(表示没有酬勤数量)
	//     bdlv   最高酬勤等级:默认值 0(表示全站都没有酬勤)
	TypeNewDeserve = "bc_buy_deserve" // 用户赠送酬勤通知消息
	// 用户赠送酬勤时,服务端发送此消息至客户端。完整的数据部分应包含的字段如 下:
	// 字段说明
	//     type   表示为“赠送酬勤通知”消息,固定为 bc_buy_deserve
	//     rid    房间 ID
	//     gid    弹幕分组 ID
	//     level  用户等级
	//     cnt    赠送数量
	//     hits   赠送连击次数
	//     lev    酬勤等级
	//     sui    用户信息序列化字符串,详见下文。注意,此处为嵌套序列化,需注 意符号的转义变换。(转义符号参见 2.2 序列化)
	TypeLiveStatusChange = "rss" // 房间开关播提醒消息
	// 房间开播提醒主要部分应包含的字段如下:
	// 字段说明
	//     type    表示为“房间开播提醒”消息,固定为 rss
	//     rid     房间 id
	//     gid     弹幕分组 id
	//     ss      直播状态,0-没有直播,1-正在直播
	//     code    类型
	//     rt      开关播原因:0-主播开关播,其他值-其他原因
	//     notify  通知类型
	//     endtime 关播时间(仅关播时有效)
	TypeRankList = "ranklist" // 广播排行榜消息
	TypeMsgToAll = "ssd"      // 超级弹幕消息(如，火箭弹幕)
	// 超级弹幕主要部分应包含的字段如下:
	// 字段说明
	//     type     表示为“超级弹幕”消息,固定为 ssd
	//     rid      房间 id
	//     gid      弹幕分组 id
	//     sdid     超级弹幕 id
	//     trid     跳转房间 id
	//     content  超级弹幕的内容
	TypeMsgToRoom = "spbc" // 房间内礼物广播
	// 房间内赠送礼物成功后效果主要部分应包含的字段如下:
	//  字段说明
	//     type   表示为“房间内礼物广播”,固定为 spbc
	//     rid    房间 id
	//     gid    弹幕分组 id
	//     sn     赠送者昵称
	//     dn     受赠者昵称
	//     gn     礼物名称
	//     gc     礼物数量
	//     drid   赠送房间
	//     gs     广播样式
	//     gb     是否有礼包(0-无礼包,1-有礼包)
	//     es     广播展现样式(1-火箭,2-飞机)
	//     gfid   礼物 id
	//     eid    特效 id
	TypeNewRedPacket = "ggbb" // 房间用户抢红包
	// 房间赠送礼物成功后效果(赠送礼物效果,连击数)主要部分应包含的字段如下:
	// 字段说明
	//     type  表示“房间用户抢红包”信息,固定为 ggbb
	//     rid   房间 id
	//     gid   弹幕分组 id
	//     sl    抢到的鱼丸数量
	//     sid   礼包产生者 id
	//     did   抢礼包者 id
	//     snk   礼包产生者昵称
	//     dnk   抢礼包者昵称
	//     rpt   礼包类型
	TypeRoomRankChange = "rankup" // 房间内top10变化消息
	// 房间内 top10 排行榜变化后,广播。主要部分应包含的字段如下:
	// 字段说明
	// type  表示为“房间 top10 排行榜变换”,固定为 rankup
	// rid   房间 id
	// gid   弹幕分组 id
	// uid   用户 id
	// drid  目标房间 id
	// rt    房间所属栏目类型
	// bt    广播类型:1-房间内广播,2-栏目广播,4-全站广播
	// sz    展示区域:1-聊天区展示,2-flash 展示,3-都显示
	// nk    用户昵称
	// rkt   top10 榜的类型 1-周榜 2-总榜 4-日榜
	// rn    上升后的排名
)

// Message 为客户端发送给弹幕服务器的消息体
type Message struct {
	BodyValues     map[string]interface{} // 消息正文map
	HeaderType     int                    // 消息类型，2字节，689为客户端发给服务器
	HeaderSecret   int8                   // 加密字段，1字节，暂时未用，默认为0
	HeaderReserved int8                   // 保留字段，1字节，暂时未用，默认为0
	Ending         int8                   // 结尾字段，1字节
}

// SetField 设置消息正文内容
func (msg *Message) SetField(name string, value interface{}) *Message {
	if msg.BodyValues == nil {
		msg.BodyValues = make(map[string]interface{})
	}
	msg.BodyValues[name] = value

	return msg
}

// Field 获取指定的字段值
func (msg *Message) Field(name string) (interface{}, bool) {
	value, ok := msg.BodyValues[name]
	return value, ok
}

// ContentString 返回正文内容字符串
func (msg *Message) ContentString() string {
	var items = make([]string, 0, len(msg.BodyValues))

	for field, value := range msg.BodyValues {
		items = append(items, fmt.Sprintf("%s@=%v/", field, value))
	}

	return strings.Join(items, "")
}

// Bytes 返回消息体的字节数组
func (msg *Message) Bytes() []byte {
	var content = msg.ContentString()
	var length = 9 + len(content) // 长度4字节 + 类型2字节 + 加密字段1字节 + 保留字段1字节 + 结尾字段1字节

	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.LittleEndian, int32(length))
	binary.Write(buffer, binary.LittleEndian, int32(length))
	binary.Write(buffer, binary.LittleEndian, int16(msg.HeaderType))
	binary.Write(buffer, binary.LittleEndian, msg.HeaderSecret)
	binary.Write(buffer, binary.LittleEndian, msg.HeaderReserved)
	binary.Write(buffer, binary.LittleEndian, []byte(content))
	binary.Write(buffer, binary.LittleEndian, msg.Ending)
	return buffer.Bytes()
}

// NewMessage 构建一个消息
func NewMessage(params ...map[string]interface{}) *Message {
	bodyValues := make(map[string]interface{})
	for _, param := range params {
		for k, v := range param {
			bodyValues[k] = v
		}
	}

	return &Message{
		BodyValues: bodyValues,
		HeaderType: TypeMessageToServer,
	}
}

// NewMessageToServer 构造一个新的客户端消息
func NewMessageToServer(params ...map[string]interface{}) *Message {
	msg := NewMessage(params...)
	msg.HeaderType = TypeMessageToServer
	return msg
}

// NewMessageFromServer 构造一个新的服务端消息
func NewMessageFromServer(content []byte) (*Message, error) {
	msg := NewMessage()
	msg.HeaderType = TypeMessageFromServer

	s := strings.Trim(string(content), "/")
	items := strings.Split(s, "/")

	for _, item := range items {
		kv := strings.SplitN(item, "@=", 2)
		msg.SetField(kv[0], kv[1])
	}

	return msg, nil
}
