package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type MsgReply struct {
	ToUserName   string
	FromUserName string
	MsgType      string
	Content      string
}

func ReplyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "请求入参:%+v", r.Host)
	data := MsgReply{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		fmt.Fprintf(w, "解析请求体失败")
		return
	}
	res := &JsonResult{}
	reply := MsgReply{
		ToUserName:   data.FromUserName,
		FromUserName: data.ToUserName,
		MsgType:      data.MsgType,
		Content:      "狐狐是垫的",
	}
	res.Data = reply

	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}
