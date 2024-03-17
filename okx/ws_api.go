package okx

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ZYKJShadow/go-okx/common/utils"
	"github.com/ZYKJShadow/recws"
	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type WebSocket struct {
	sync.RWMutex

	wsURL string

	ctx    context.Context
	cancel context.CancelFunc
	conn   recws.RecConn

	subscriptions map[string]interface{}
}

type Trade struct {
	SprdId  string `json:"sprdId"`
	Side    string `json:"side"`
	Sz      string `json:"sz"`
	Px      string `json:"px"`
	TradeId string `json:"tradeId"`
	Ts      string `json:"ts"`
}

type SocketMsg struct {
	Event string `json:"event"`
	Code  string `json:"code,omitempty"`
	Msg   string `json:"msg,omitempty"`
}

type SubscribeMsg struct {
	Op   string         `json:"op"`
	Args []SubscribeArg `json:"args"`
}

type SubscribeArg struct {
	Channel string `json:"channel"`
	InstId  string `json:"instId,omitempty"`
	SprdId  string `json:"sprdId,omitempty"`
}

func InitWebSocket(url string) *WebSocket {
	ws := &WebSocket{
		wsURL:         url,
		subscriptions: make(map[string]interface{}),
	}
	ws.ctx, ws.cancel = context.WithCancel(context.Background())
	ws.conn = recws.RecConn{
		KeepAliveTimeout: 10 * time.Second,
	}

	ws.conn.SubscribeHandler = ws.subscribeHandler
	return ws
}

func (w *WebSocket) SetProxy(proxyURL string) (err error) {
	if proxyURL == "" {
		return
	}
	var purl *url.URL
	purl, err = url.Parse(proxyURL)
	if err != nil {
		return
	}
	logx.Infof("[ws][%s] proxy url:%s", proxyURL, purl)
	w.conn.Proxy = http.ProxyURL(purl)
	return
}

func (w *WebSocket) Start() {
	logx.Infof("wsURL: %v", w.wsURL)
	w.conn.Dial(w.wsURL, nil)
	go w.run()
}

func (w *WebSocket) Subscribe(channel string, sub *SubscribeMsg) error {
	w.Lock()
	defer w.Unlock()
	w.subscriptions[channel] = sub
	return w.sendWSMessage(sub)
}

func (w *WebSocket) subscribeHandler() error {
	for _, v := range w.subscriptions {
		err := w.sendWSMessage(v)
		if err != nil {
			logx.Errorf("%v", err)
			return err
		}
	}
	return nil
}

func (w *WebSocket) sendWSMessage(msg interface{}) error {
	return w.conn.WriteJSON(msg)
}

func (w *WebSocket) run() {
	ctx := context.Background()
	for {
		select {
		case <-ctx.Done():
			go w.conn.Close()
			logx.Infof("Websocket closed %s", w.conn.GetURL())
			return
		default:
			_, msg, err := w.conn.ReadMessage()
			if err != nil {
				logx.Errorf("Read error: %v", err)
				time.Sleep(time.Millisecond * 200)
				continue
			}
			w.handle(msg)
		}
	}
}

func (w *WebSocket) handle(b []byte) {
	var msg SocketMsg
	err := json.Unmarshal(b, &msg)
	if err != nil {
		logx.Error(err)
		return
	}

	switch msg.Event {
	case "error":
		logx.Error(msg.Msg)
	case "subscribe":
		logx.Infof("subscribe success: %v", string(b))
	case "unsubscribe":
		logx.Infof("unsubscribe success: %v", string(b))
	default:
		ret := gjson.ParseBytes(b)
		channel := ret.Get("arg").Get("channel").String()
		action := ret.Get("action").String()
		w.subscribe(channel, action, ret.Get("data").String())
	}
}

func (w *WebSocket) subscribe(channel string, action string, data string) {
	switch channel {
	case "trades":
		var trades []*Trade
		err := json.Unmarshal([]byte(data), &trades)
		if err != nil {
			logx.Error(err)
			return
		}

		trade := trades[0]
		sz := utils.MustParseFloat64(trade.Sz)
		ts := utils.MustParseInt64(trade.Ts)
		if sz >= 1000 {
			fmt.Println(channel, time.UnixMilli(ts).Format("2006-01-02 15:04:05"), trade.Side, trade.Px, sz)
		}
	case "books":
		var books []*Book
		err := json.Unmarshal([]byte(data), &books)
		if err != nil {
			logx.Error(err)
			return
		}

		if action == "snapshot" {

		} else if action == "update" {

		}
	}
}
