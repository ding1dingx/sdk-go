package corp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tidwall/gjson"

	"github.com/shenghui0779/gochat/event"
	"github.com/shenghui0779/gochat/urls"
	"github.com/shenghui0779/gochat/wx"
)

type Corp struct {
	corpid string
	token  string
	aeskey string
	nonce  func() string
	client wx.HTTPClient
}

func (corp *Corp) CorpID() string {
	return corp.corpid
}

// OAuth2URL 生成网页授权URL（请使用 URLEncode 对 redirectURL 进行处理）
// [参考](https://open.work.weixin.qq.com/api/doc/90000/90135/91020)
func (corp *Corp) OAuth2URL(scope AuthScope, redirectURL, state string) string {
	return fmt.Sprintf("%s?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect", urls.Oauth2Authorize, corp.corpid, redirectURL, scope, state)
}

// QRCodeAuthURL 生成扫码授权URL（请使用 URLEncode 对 redirectURL 进行处理）
// [参考](https://open.work.weixin.qq.com/api/doc/90000/90135/90988)
func (corp *Corp) QRCodeAuthURL(agentID, redirectURL, state string) string {
	return fmt.Sprintf("%s?appid=%s&agentid=%s&redirect_uri=%s&state=%s", urls.QRCodeAuthorize, corp.corpid, agentID, redirectURL, state)
}

func (corp *Corp) AccessToken(ctx context.Context, secret string, options ...wx.HTTPOption) (*AccessToken, error) {
	resp, err := corp.client.Do(ctx, http.MethodGet, fmt.Sprintf("%s?corpid=%s&corpsecret=%s", urls.CorpCgiBinAccessToken, corp.corpid, secret), nil, options...)

	if err != nil {
		return nil, err
	}

	r := gjson.ParseBytes(resp)

	if code := r.Get("errcode").Int(); code != 0 {
		return nil, fmt.Errorf("%d|%s", code, r.Get("errmsg").String())
	}

	token := new(AccessToken)

	if err = json.Unmarshal(resp, token); err != nil {
		return nil, err
	}

	return token, nil
}

// Do exec action
func (corp *Corp) Do(ctx context.Context, accessToken string, action wx.Action, options ...wx.HTTPOption) error {
	var (
		resp []byte
		err  error
	)

	if action.IsUpload() {
		form, ferr := action.UploadForm()

		if ferr != nil {
			return ferr
		}

		resp, err = corp.client.Upload(ctx, action.URL(accessToken), form, options...)
	} else {
		body, berr := action.Body()

		if berr != nil {
			return berr
		}

		resp, err = corp.client.Do(ctx, action.Method(), action.URL(accessToken), body, options...)

		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	r := gjson.ParseBytes(resp)

	if code := r.Get("errcode").Int(); code != 0 {
		return fmt.Errorf("%d|%s", code, r.Get("errmsg").String())
	}

	return action.Decode(resp)
}

// VerifyEventSign 验证事件消息签名
// 验证消息来自微信服务器，使用：msg_signature、timestamp、nonce、echostr（若验证成功，解密echostr后返回msg字段内容）
// 验证事件消息签名，使用：msg_signature、timestamp、nonce、msg_encrypt
// [参考](https://developer.work.weixin.qq.com/document/path/90930)
func (corp *Corp) VerifyEventSign(signature string, items ...string) bool {
	signStr := event.SignWithSHA1(corp.token, items...)

	return signStr == signature
}

// DecryptEventMessage 事件消息解密
func (corp *Corp) DecryptEventMessage(encrypt string) (wx.WXML, error) {
	b, err := event.Decrypt(corp.corpid, corp.aeskey, encrypt)

	if err != nil {
		return nil, err
	}

	return wx.ParseXML2Map(b)
}

type CorpOption func(corp *Corp)

// WithServerConfig 设置服务器配置
// [参考](https://open.work.weixin.qq.com/api/doc/90000/90135/90930)
func WithServerConfig(token, aeskey string) CorpOption {
	return func(corp *Corp) {
		corp.token = token
		corp.aeskey = aeskey
	}
}

// WithNonce 设置 Nonce（加密随机串）
func WithNonce(f func() string) CorpOption {
	return func(corp *Corp) {
		corp.nonce = f
	}
}

// WithClient 设置 HTTP Client
func WithClient(c *http.Client) CorpOption {
	return func(corp *Corp) {
		corp.client = wx.NewHTTPClient(c)
	}
}

// WithMockClient 设置 Mock Client
func WithMockClient(c wx.HTTPClient) CorpOption {
	return func(corp *Corp) {
		corp.client = c
	}
}

func New(corpid string, options ...CorpOption) *Corp {
	corp := &Corp{
		corpid: corpid,
		nonce: func() string {
			return wx.Nonce(16)
		},
		client: wx.NewDefaultClient(),
	}

	for _, f := range options {
		f(corp)
	}

	return corp
}
