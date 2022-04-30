package dingbot

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const defaultDingTalkURL = `https://oapi.dingtalk.com/robot/send`

type Dingbot struct {
	Token  string `yaml:"token" json:"token"`
	Secret string `yaml:"secret" json:"secret"`
}

func New(token, secret string) *Dingbot {
	if token == "" || secret == "" {
		return nil
	}
	return &Dingbot{
		Token:  token,
		Secret: secret,
	}
}

func (d *Dingbot) Notify(title, text string) error {
	url, err := d.makeSign()
	if err != nil {
		return err
	}
	msg, err := d.makeMarkdownMsg(title, text)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewReader(msg))
	if err != nil {
		return err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d, body: %v", resp.StatusCode, string(data))
	}

	return nil
}

func (d *Dingbot) makeSign() (string, error) {
	var (
		timestamp, sign string
	)
	timestamp = strconv.FormatInt(time.Now().Unix()*1000, 10)
	h := hmac.New(sha256.New, []byte(d.Secret))
	if _, err := h.Write([]byte(fmt.Sprintf("%s\n%s", timestamp, d.Secret))); err != nil {
		return "", err
	}
	sign = base64.StdEncoding.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("%s?access_token=%s&timestamp=%s&sign=%s", defaultDingTalkURL, d.Token, timestamp, sign), nil
}

type MarkdownMsg struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
}

func (d *Dingbot) makeMarkdownMsg(title, text string) ([]byte, error) {
	msg := &MarkdownMsg{
		MsgType: "markdown",
	}
	msg.Markdown.Title = title
	msg.Markdown.Text = fmt.Sprintf("### %s\n---\n%s", title, text)
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return data, nil
}
