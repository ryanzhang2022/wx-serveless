package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type MsgReply struct {
	ToUserName   string
	FromUserName string
	MsgType      string
	Content      string
	CreateTime   string
	Action       string `json:"action"`
}

func ReplyHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("请求到了====")

	decoder := json.NewDecoder(r.Body)
	body := make(map[string]any)

	if err := decoder.Decode(&body); err != nil {
		fmt.Fprintf(w, "解析请求体失败")
		return
	}
	defer r.Body.Close()

	reply := map[string]any{
		"ToUserName":   body["FromUserName"],
		"FromUserName": body["ToUserName"],
		"CreateTime":   strconv.Itoa(int(time.Now().Unix())),
		"MsgType":      "text",
		"Content":      "狐狐是垫的",
	}

	msg, err := json.Marshal(reply)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}
