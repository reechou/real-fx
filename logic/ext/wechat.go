package ext

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/real-fx/config"
)

const (
	WECHAT_WITHDRAWAL = "/index.php?r=weixinpay/pay"
)

type WechatExt struct {
	cfg *config.Config

	client *http.Client
}

func NewWechatExt(cfg *config.Config) *WechatExt {
	return &WechatExt{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (we *WechatExt) Withdrawal(info *WithdrawalReq) error {
	u := we.cfg.WechatExtInfo.HostURL + WECHAT_WITHDRAWAL
	body, err := json.Marshal(info)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequest("POST", u, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	rsp, err := we.client.Do(httpReq)
	defer func() {
		if rsp != nil {
			rsp.Body.Close()
		}
	}()
	if err != nil {
		return err
	}
	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	var response WechatResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		return err
	}
	if response.Code != WECHAT_RESPONSE_OK {
		logrus.Errorf("wechat withdrawal error: %v", response)
		return fmt.Errorf("wechat withdrawal error: %s", response.Msg)
	}

	return nil
}
