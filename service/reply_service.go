package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ReplyReq struct {
	FromUserName string
	ToUserName   string
	MsgType      string
	Content      string
}

//func ReplyHandler(w http.ResponseWriter, r *http.Request) {
//	log.Print("请求到了====")
//
//	decoder := json.NewDecoder(r.Body)
//	body := &ReplyReq{}
//
//	if err := decoder.Decode(&body); err != nil {
//		fmt.Fprintf(w, "解析请求体失败")
//		return
//	}
//	defer r.Body.Close()
//
//	reply := map[string]interface{}{
//		"ToUserName":   body.FromUserName,
//		"FromUserName": body.ToUserName,
//		"CreateTime":   strconv.Itoa(int(time.Now().Unix())),
//		"MsgType":      "text",
//		"Content":      "狐狐是垫的",
//	}
//
//	msg, err := json.Marshal(reply)
//	if err != nil {
//		fmt.Fprint(w, "内部错误")
//		return
//	}
//	w.Header().Set("content-type", "application/json")
//	w.Write(msg)
//}

func ReplyHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("请求到了e====")

	decoder := json.NewDecoder(r.Body)
	body := &ReplyReq{}

	if err := decoder.Decode(&body); err != nil {
		fmt.Fprintf(w, "解析请求体失败")
		return
	}
	defer r.Body.Close()

	if strings.Contains(body.Content, "gpt") {
		go fetchChatGptSend(body.Content, body.FromUserName)
		w.Write([]byte("success"))
	} else {
		reply := map[string]interface{}{
			"ToUserName":   body.FromUserName,
			"FromUserName": body.ToUserName,
			"CreateTime":   strconv.Itoa(int(time.Now().Unix())),
			"MsgType":      "text",
			"Content":      "狐狐是垫的",
		}
		msg, _ := json.Marshal(reply)
		w.Header().Set("content-type", "application/json")
		w.Write(msg)
	}
}

type SendReplyReq struct {
	ToUser  string `json:"touser"`
	MsgType string `json:"msgtype"`
	Text    struct {
		Context string `json:"context"`
	} `json:"text"`
}
type ChatGptReq struct {
	Prompt           string  `json:"prompt"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float32 `json:"temperature"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
	Model            string  `json:"model"`
}

type ChatGptResp struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func fetchChatGptSend(question string, toUser string) {
	chatapi := "https://api.openai.com/v1/completions"
	wxReplyApi := "http://api.weixin.qq.com/cgi-bin/message/custom/send"

	gptReq := ChatGptReq{
		Prompt:           question,
		MaxTokens:        2048,
		Temperature:      0.1,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		Model:            "text-davinci-003",
	}

	gptResp := ChatGptResp{}
	if err := httpPost(chatapi, gptReq, &gptResp); err != nil {
		log.Print(err)
		return
	}
	log.Print(gptResp)

	send := SendReplyReq{
		ToUser:  toUser,
		MsgType: "text",
		Text: struct {
			Context string `json:"context"`
		}{gptResp.Choices[0].Text},
	}

	respMap := map[string]interface{}{}
	if err := httpPost(wxReplyApi, send, &respMap); err != nil {
		log.Print(err)
		return
	}

}

func httpPost(url string, req interface{}, resp interface{}) (err error) {
	client := http.Client{}

	reqByte, err := json.Marshal(req)
	if err != nil {
		log.Print(err)
		return
	}
	reqReader := bytes.NewReader(reqByte)

	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, reqReader)
	if err != nil {
		log.Print(err)
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if strings.Contains(url, "openai.com") {
		httpReq.Header.Set("Authorization", "Bearer "+genKey())
	}
	tempResp, err := client.Do(httpReq)
	if err != nil {
		log.Print(err)
		return
	}
	respBody, err := io.ReadAll(tempResp.Body)
	if err != nil {
		log.Print(err)
		return
	}
	defer tempResp.Body.Close()

	log.Printf("http-resp:%s\n", string(respBody))
	if err = json.Unmarshal(respBody, resp); err != nil {
		log.Print(err)
		return
	}
	return nil
}

func genKey() string {
	str := "u09q-ksige4jUF2nbCrxNblB3TaKhDjJJFkVU168wTit8FR5g8r"

	parts := 8
	length := len(str)
	partLength := int(math.Ceil(float64(length) / float64(parts)))

	strArr := make([]string, parts)
	k := 0
	for i := 0; i < length; i += partLength {
		end := int(math.Min(float64(i+partLength), float64(length)))
		strArr[k] = reverse(strings.TrimSpace(str[i:end]))
		k++
	}
	return strings.Join(strArr, "")
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
