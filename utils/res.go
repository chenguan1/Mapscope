package utils

import (
	"github.com/kataras/iris/v12/context"
	"net/http"
)

// Res 返回数据模板
type Res struct {
	ctx  context.Context `json:-`
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data interface{}     `json:"data"`
}

var codes = map[int]string{
	0:                              "Error",
	http.StatusOK:                  "OK",                  // 200
	http.StatusBadRequest:          "BadRequest",          // 400
	http.StatusInternalServerError: "InternalServerError", // 500
}

// NewRes Create ResModel
func NewRes(c context.Context) *Res {
	return &Res{
		ctx:  c,
		Code: http.StatusOK,
		Msg:  codes[http.StatusOK],
	}
}

// reset
func (res *Res) Reset(c context.Context) {
	res.ctx = c
	res.Code = http.StatusOK
	res.Msg = codes[http.StatusOK]
	res.Data = nil
}

// Fail
func (res *Res) Fail() {
	res.Code = 0
	res.Msg = codes[0]
	res.ctx.JSON(res)
}

// FailCode
func (res *Res) FailCode(code int) {
	res.Code = code
	res.Msg = codes[code]
	res.ctx.JSON(res)
}

// FailErr
func (res *Res) FailErr(err error) {
	res.Code = 0
	res.Msg = codes[0]
	if err != nil {
		res.Msg = err.Error()
	}
	res.ctx.JSON(res)
}

// FailMsg
func (res *Res) FailMsg(msg string) {
	res.Code = 0
	res.Msg = msg
	res.ctx.JSON(res)
}

// FailData
func (res *Res) FailData(data interface{}) {
	res.Code = 0
	res.Msg = codes[0]
	res.Data = data
	res.ctx.JSON(res)
}

// Done
func (res *Res) Done() {
	res.Code = http.StatusOK
	res.Msg = codes[http.StatusOK]
	res.ctx.JSON(res)
}

// DoneMsg
func (res *Res) DoneMsg(msg string) {
	res.Code = http.StatusOK
	res.Msg = msg
	res.ctx.JSON(res)
}

// DoneData
func (res *Res) DoneData(data interface{}) {
	res.Code = http.StatusOK
	res.Msg = codes[http.StatusOK]
	res.Data = data
	res.ctx.JSON(res)
}
