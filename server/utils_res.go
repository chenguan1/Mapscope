package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var codes = map[int]string{
	0: "检查消息",

	200: "成功",
	201: "已创建",
	202: "已接受",
	204: "无内容",

	400:  "请求无法解析",
	4001: "必填参数校验错误",
	4002: "达到最大尝试登录次数,稍后再试",
	4003: "瓦片请求格式错误",
	4004: "符号请求格式错误",
	4005: "字体请求格式错误",

	401:  "未授权",
	4011: "用户名或密码错误",
	4012: "用户名非法,请使用字母,数字,短划线,下划线组合或用户名需少于32个字符",
	4013: "邮箱非法,请使用能收到验证邮件的正确邮箱",
	4014: "密码非法,请使用至少4位以上密码字符",
	4015: "用户名已注册,请使用新的用户名",
	4016: "邮箱已注册,请使用新的邮箱",

	403:  "禁止访问",
	4031: "邮箱不存在",

	404:  "找不到资源",
	4041: "用户不存在",
	4042: "角色不存在",
	4043: "服务不存在",
	4044: "找不到样式",
	4045: "找不到瓦片集",
	4046: "找不到数据集",
	4047: "找不到字体库",
	4048: "找不到上传文件",
	4049: "地图不存在",

	408: "请求超时",

	422: "数据集不存在",

	500:  "系统错误",
	5001: "数据库错误",
	5002: "文件读写错误",
	5003: "IO读写错误",
	5004: "MBTiles读写错误",
	5005: "系统配置错误",

	501: "维护中",
	503: "服务不可用",
}

type Res struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

//NewRes Create Res
func NewRes() *Res {
	return &Res{
		Code: http.StatusOK,
		Msg:  codes[http.StatusOK],
	}
}

//Fail failed error
func (res *Res) Fail(c *gin.Context, code int) {
	res.Code = code
	res.Msg = codes[code]
	c.JSON(http.StatusOK, res)
}

//FailErr failed string
func (res *Res) FailErr(c *gin.Context, err error) {
	res.Code = 0
	if err != nil {
		res.Msg = err.Error()
	}
	c.JSON(http.StatusOK, res)
}

//FailMsg failed string
func (res *Res) FailMsg(c *gin.Context, msg string) {
	res.Code = 0
	res.Msg = msg
	c.JSON(http.StatusOK, res)
}

//Done done
func (res *Res) Done(c *gin.Context, msg string) {
	res.Code = http.StatusOK
	res.Msg = codes[http.StatusOK]
	if msg != "" {
		res.Msg = msg
	}
	c.JSON(http.StatusOK, res)
}

//DoneCode done
func (res *Res) DoneCode(c *gin.Context, code int) {
	res.Code = code
	res.Msg = codes[code]
	c.JSON(http.StatusOK, res)
}

//DoneData done
func (res *Res) DoneData(c *gin.Context, data interface{}) {
	res.Code = http.StatusOK
	res.Msg = codes[http.StatusOK]
	res.Data = data
	c.JSON(http.StatusOK, res)
}

//Reset reset to init
func (res *Res) Reset() {
	res.Code = http.StatusOK
	res.Msg = codes[http.StatusOK]
}