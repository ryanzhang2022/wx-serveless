package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type MsgReply struct {
	ToUserName   string
	FromUserName string
	MsgType      string
	Content      string
	Action       string `json:"action"`
}

func ReplyHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("请求到了====")
	data := MsgReply{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		fmt.Fprintf(w, "解析请求体失败")
		return
	}
	if data.Action != "" {
		w.Write([]byte("success"))
		return
	}

	reply := MsgReply{
		ToUserName:   data.FromUserName,
		FromUserName: data.ToUserName,
		MsgType:      data.MsgType,
		Content:      "狐狐是垫的",
	}

	msg, err := json.Marshal(reply)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}
