package wechat

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"encoding/xml"
	"github.com/rdlucklib/rdluck_tools/nethttp"
	"github.com/rdlucklib/rdluck_tools/str"
	"github.com/rdlucklib/rdluck_tools/uuid"
	"runtime"
	"strconv"

	"github.com/pkg/errors"
)

//微信支付接口地址
const (
	WxUnifiedorder      = "https://api.mch.weixin.qq.com/pay/unifiedorder"                //统一下单
	WxOrderquery        = "https://api.mch.weixin.qq.com/pay/orderquery"                  //查询订单
	WxCloseorder        = "https://api.mch.weixin.qq.com/pay/closeorder"                  //关闭订单
	WxRefund            = "https://api.mch.weixin.qq.com/secapi/pay/refund"               //申请退款
	WxRefundquery       = "https://api.mch.weixin.qq.com/pay/refundquery"                 //查询退款
	WxDownloadbill      = "https://api.mch.weixin.qq.com/pay/downloadbill"                //下载对账单
	WxReport            = "https://api.mch.weixin.qq.com/payitil/report"                  //交易保障
	WxBatchquerycomment = "https://api.mch.weixin.qq.com/billcommentsp/batchquerycomment" //拉取订单评价信息

	//支付回调地址，必须是http,否则回调不到
	WXNotify_url = "http://wrhb.weiyunjinrong.cn/v1/callback/unifiedordercallback"
	//企业付款到用户零钱
	WXTransfers = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"
)

//微信支付配置
const (
	WXPayMchId  = "1496554202"         //支付商户号
	WXPayAppId  = "wxc6c84e8998fd22a0" //支付对应应用标识
	WXPaySecret = "2cb47ae1d364cb28e9a7012a7c886dad"
	WXPayApiKey = "57353dd828374b45be7c72334260645e"
)

//微信证书
const (
	WxCertPath = "./cert/apiclient_cert.pem"
	WxKeyPath  = "./cert/apiclient_key.pem"
	WxCAPath   = "./cert/rootca.pem"
)

type WxPay struct {
	OpenId         string
	TotalFee       int
	Body           string
	CallbackResult []byte
	Ip             string
	Desc           string
}

type Unifiedorder struct {
	Appid        string `xml:"appid"`        //小程序ID
	Mch_id       string `xml:"mch_id"`       //商户号
	Nonce_str    string `xml:"nonce_str"`    //随机字符串
	Sign         string `xml:"sign"`         //签名
	Out_trade_no string `xml:"out_trade_no"` //商户订单号
	Body         string `xml:"body"`         //商品描述
	Openid       string `xml:"openid"`       //用户标识
	Total_fee    int    `xml:"total_fee"`    //总金额
	Trade_type   string `xml:"trade_type"`   //交易类型
	Notify_url   string `xml:"notify_url"`   //通知地址
}

type WxPayResult struct {
	Return_code  string `xml:"return_code"`
	Return_msg   string `xml:"return_msg"`
	Appid        string `xml:"appid"`
	Mch_id       string `xml:"mch_id"`
	Nonce_str    string `xml:"nonce_str"`
	Sign         string `xml:"sign"`
	Result_code  string `xml:"result_code"`
	Prepay_id    string `xml:"prepay_id"`
	Trade_type   string `xml:"trade_type"`
	Out_trade_no string
}

type Notify struct {
	Appid          string  `xml:"appid"`
	Bank_type      string  `xml:"bank_type"`
	Cash_fee       float64 `xml:"cash_fee"`
	Fee_type       string  `xml:"fee_type"`
	Is_subscribe   string  `xml:"is_subscribe"`
	Mch_id         string  `xml:"mch_id"`
	Nonce_str      string  `xml:"nonce_str"`
	Openid         string  `xml:"openid"`
	Out_trade_no   string  `xml:"out_trade_no"`
	Result_code    string  `xml:"result_code"`
	Return_code    string  `xml:"return_code"`
	Sign           string  `xml:"sign"`
	Time_end       string  `xml:"time_end"`
	Total_fee      int     `xml:"total_fee"`
	Trade_type     string  `xml:"trade_type"`
	Transaction_id string  `xml:"transaction_id"`
}

type Transfers struct {
	Mch_appid        string `xml:"mch_appid"`
	Mchid            string `xml:"mchid"`
	Nonce_str        string `xml:"nonce_str"`
	Sign             string `xml:"sign"`
	Partner_trade_no string `xml:"partner_trade_no"`
	Openid           string `xml:"openid"`
	Check_name       string `xml:"check_name"`
	Amount           int    `xml:"amount"`
	Desc             string `xml:"desc"`
	Spbill_create_ip string `xml:"spbill_create_ip"`
}

type TransfersResult struct {
	Return_code      string `xml:"return_code"`
	Return_msg       string `xml:"return_msg"`
	Mchid            string `xml:"mchid"`
	Nonce_str        string `xml:"nonce_str"`
	Result_code      string `xml:"result_code"`
	Partner_trade_no string `xml:"partner_trade_no"`
	Payment_no       string `xml:"payment_no"`
	Payment_time     string `xml:"payment_time"`
}

func NewWxPay() *WxPay {
	return &WxPay{}
}

//微信支付统一下单
func (pay *WxPay) WxUnifiedorder() (*WxPayResult, error) {
	outTradeNo := uuid.NewUUID().Hex32()
	randStr := str.GetRandString(16)
	o := new(Unifiedorder)
	o.Appid = WXPayAppId
	o.Mch_id = WXPayMchId
	o.Nonce_str = randStr
	o.Out_trade_no = outTradeNo
	o.Body = pay.Body
	o.Openid = pay.OpenId
	o.Total_fee = pay.TotalFee
	o.Trade_type = "JSAPI"
	o.Notify_url = WXNotify_url

	m := make(map[string]interface{}, 0)
	m["appid"] = o.Appid
	m["mch_id"] = o.Mch_id
	m["nonce_str"] = o.Nonce_str
	m["out_trade_no"] = o.Out_trade_no
	m["body"] = o.Body
	m["openid"] = o.Openid
	m["total_fee"] = o.Total_fee
	m["trade_type"] = o.Trade_type
	m["notify_url"] = o.Notify_url
	sign := WxSign(m, WXPayApiKey)
	fmt.Println("sign:", sign)
	o.Sign = sign
	reqData, err := xml.Marshal(o)
	if err != nil {
		fmt.Println("xml编码失败，原因:", err)
		return nil, err
	}
	reqStr := string(reqData)
	reqStr = strings.Replace(reqStr, "Unifiedorder", "xml", -1)
	reqResult, err := nethttp.HttpPost(WxUnifiedorder, reqStr, nethttp.AcceptXml)
	//返回结果
	var r WxPayResult
	err = xml.Unmarshal(reqResult, &r)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	r.Out_trade_no = outTradeNo
	return &r, nil
}

//微信支付统一下单回调
func (pay *WxPay) WxUnifiedorderCallback() (bool, error) {
	if string(pay.CallbackResult) == "" {
		return false, errors.New("回调结果不能为空")
	}
	var notify Notify
	err := xml.Unmarshal(pay.CallbackResult, &notify)
	if err != nil {
		return false, errors.New("序列化回调结果失败:err:" + err.Error())
	}
	var reqMap map[string]interface{}
	reqMap = make(map[string]interface{}, 0)
	reqMap["appid"] = notify.Appid
	reqMap["bank_type"] = notify.Bank_type
	reqMap["cash_fee"] = notify.Cash_fee
	reqMap["fee_type"] = notify.Fee_type
	reqMap["is_subscribe"] = notify.Is_subscribe
	reqMap["mch_id"] = notify.Mch_id
	reqMap["nonce_str"] = notify.Nonce_str
	reqMap["openid"] = notify.Openid
	reqMap["out_trade_no"] = notify.Out_trade_no
	reqMap["result_code"] = notify.Result_code
	reqMap["return_code"] = notify.Return_code
	reqMap["time_end"] = notify.Time_end
	reqMap["total_fee"] = notify.Total_fee
	reqMap["trade_type"] = notify.Trade_type
	reqMap["transaction_id"] = notify.Transaction_id
	state := 0
	fmt.Println(state)
	//进行签名校验
	if WxVerifySign(reqMap, notify.Sign, WXPayApiKey) {
		if notify.Return_code != "SUCCESS" {
			return false, errors.New("请求通讯失败")
		}
		if notify.Result_code != "SUCCESS" {
			return false, errors.New("支付失败:" + notify.Result_code)
		}
		if notify.Total_fee != pay.TotalFee {
			return false, errors.New("支付金额错误:订单支付金额:" + strconv.Itoa(pay.TotalFee) + ";回调返回金额:" + strconv.Itoa(notify.Total_fee))
		}
		if notify.Return_code == "SUCCESS" && notify.Result_code == "SUCCESS" && notify.Total_fee != pay.TotalFee {
			return true, nil
		} else {
			return false, errors.New("未知错误:callback data:" + string(pay.CallbackResult))
		}
	} else {
		return false, errors.New("签名错误")
	}
	return false, nil
}

//企业付款到用户零钱
func (pay *WxPay) WxWithdrawals() (*TransfersResult, error) {
	outTradeNo := uuid.NewUUID().Hex32()
	randStr := str.GetRandString(16)
	t := new(Transfers)
	t.Mch_appid = WXPayAppId
	t.Mchid = WXPayMchId
	t.Nonce_str = randStr
	t.Partner_trade_no = outTradeNo
	t.Openid = pay.OpenId
	t.Check_name = "NO_CHECK"
	t.Amount = pay.TotalFee
	t.Desc = pay.Desc
	t.Spbill_create_ip = pay.Ip
	m := make(map[string]interface{}, 0)
	m["mch_appid"] = t.Mch_appid
	m["mchid"] = t.Mchid
	m["nonce_str"] = t.Nonce_str
	m["partner_trade_no"] = t.Partner_trade_no
	m["openid"] = t.Openid
	m["check_name"] = t.Check_name
	m["amount"] = t.Amount
	m["desc"] = t.Desc
	m["spbill_create_ip"] = t.Spbill_create_ip
	t.Sign = WxSign(m, WXPayApiKey)
	byte_req, err := xml.Marshal(t)
	if err != nil {
		fmt.Println("xml编码错误, 原因:", err)
		return nil, err
	}
	str_req := string(byte_req)
	str_req = strings.Replace(str_req, "Transfers", "xml", -1)
	bytes_req := []byte(str_req)
	result, err := nethttp.HttpPostHasCertificate(WXTransfers, nethttp.Content_Type_Text_Xml, WxCAPath, WxCertPath, WxKeyPath, bytes_req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var tr TransfersResult
	err = xml.Unmarshal(result, &tr)
	return &tr, err
}

func WxSign(req map[string]interface{}, key string) (sign string) {
	//对key进行升序排序
	sorted_keys := make([]string, 0)
	for k, _ := range req {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)
	//对key=value的键值对用&连接起来，略过空值
	var signString string
	for _, k := range sorted_keys {
		fmt.Printf("k=%v,v=%v\n", k, req[k])
		value := fmt.Sprintf("%v", req[k])
		if value != "" {
			signString = signString + k + "=" + value + "&"
		}
	}
	//在键值对的最后加上key=API_KEY
	if key != "" {
		signString = signString + "key=" + key
	}
	//进行MD5签名并且将所有字符转为大写.
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signString))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	fmt.Println("微信支付签名结果:", upperSign)
	return upperSign
}

func WxVerifySign(needVerify map[string]interface{}, sign, key string) bool {
	pc, _, line, _ := runtime.Caller(0)
	fc := runtime.FuncForPC(pc)
	signCalc := WxSign(needVerify, key)
	fmt.Println(fc.Name(), line, "计算出来的sign: ", signCalc)
	fmt.Println(fc.Name(), line, "微信异步通知sign: ", sign)
	if signCalc == sign {
		fmt.Println(fc.Name(), line, "签名通过")
		return true
	}
	fmt.Println(fc.Name(), line, "签名失败")
	return false
}
