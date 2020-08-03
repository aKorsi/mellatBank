package mellatBank

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tiaguinho/gosoap"
)

const (
	serviceURL = "https://bpm.shaparak.ir/pgwchannel/services/pgw?wsdl"
	GatewayURL = "https://bpm.shaparak.ir/pgwchannel/startpay.mellat"
	ymd        = "20060102"
	his        = "150405"
)

type Client struct {
	TerminalId string
	UserName   string
	Password   string
}

type responseReturn struct {
	Return string `xml:"return"`
}

func (c *Client) getBankError(result string) error {
	switch result {
	case "0":
		return nil
	case "11":
		return errors.New("شماره کارت نامعتبر است")
	case "12":
		return errors.New("موجودی کافی نیست")
	case "13":
		return errors.New("رمز نادرست است")
	case "14":
		return errors.New("تعداد دفعات وارد کردن رمز بیش از حد مجاز است")
	case "15":
		return errors.New("کارت نامعتبر است")
	case "16":
		return errors.New("دفعات برداشت وجه بیش از حد مجاز است")
	case "17":
		return errors.New("کاربر از انجام تراکنش منصرف شده است")
	case "18":
		return errors.New("تاریخ انقضای کارت گذشته است")
	case "19":
		return errors.New("مبلغ برداشت وجه بیش از حد مجاز است")
	case "111":
		return errors.New("صادر کننده کارت نامعتبر است")
	case "112":
		return errors.New("خطای سوییچ صادر کننده کارت")
	case "113":
		return errors.New("پاسخی از صادر کننده کارت دریافت نشد")
	case "114":
		return errors.New("دارنده این کارت مجاز به انجام این تراکنش نیست")
	case "21":
		return errors.New("پذیرنده نامعتبر است")
	case "23":
		return errors.New("خطای امنیتی رخ داده است")
	case "24":
		return errors.New("اطلاعات کاربری پذیرنده نامعتبر است")
	case "25":
		return errors.New("مبلغ نامعتبر است")
	case "31":
		return errors.New("پاسخ نامعتبر است")
	case "32":
		return errors.New("فرمت اطلاعات وارد شده صحیح نمی‌باشد")
	case "33":
		return errors.New("حساب نامعتبر است")
	case "34":
		return errors.New("خطای سیستمی")
	case "35":
		return errors.New("تاریخ نامعتبر است")
	case "41":
		return errors.New("شماره درخواست تکراری است")
	case "42":
		return errors.New("تراکنش Sale یافت نشد")
	case "43":
		return errors.New("قبلا درخواست Verify داده شده است")
	case "44":
		return errors.New("درخواست Verify یافت نشد")
	case "45":
		return errors.New("تراکنش Settle شده است")
	case "46":
		return errors.New("تراکنش Settle نشده است")
	case "47":
		return errors.New("تراکنش Settle یافت نشد")
	case "48":
		return errors.New("تراکنش Reverse شده است")
	case "49":
		return errors.New("تراکنش Refund یافت نشد")
	case "412":
		return errors.New("شناسه قبض نادرست است")
	case "413":
		return errors.New("شناسه پرداخت نادرست است")
	case "414":
		return errors.New("سازمان صادر کننده قبض نامعتبر است")
	case "415":
		return errors.New("مدت زمان مجاز برای انجام تراکنش به پایان رسیده است")
	case "416":
		return errors.New("خطا در ثبت اطلاعات")
	case "417":
		return errors.New("شناسه پرداخت کننده نامعتبر است")
	case "418":
		return errors.New("اشکال در تعریف اطلاعات مشتری")
	case "419":
		return errors.New("تعداد دفعات ورود اطلاعات از حد مجاز گذشته است")
	case "421":
		return errors.New("IP نامعتبر است")
	case "51":
		return errors.New("تراکنش تکراری است")
	case "54":
		return errors.New("تراکنش مرجع موجود نیست")
	case "55":
		return errors.New("تراکنش نامعتبر است")
	case "61":
		return errors.New("خطا در واریز")
	case "62":
		return errors.New("مسير بازگشت به سايت در دامنه ثبت شده برای پذيرنده قرار ندارد")
	case "98":
		return errors.New("سقف استفاده از رمز ايستا به پايان رسيده است")
	default:
		return errors.New("خطای ناشناخته")
	}
}

func (c *Client) PaymentRequest(orderId, amount int64, localDate, localTime time.Time, additionalData, callBackUrl string) (string, error) {
	soap, err := gosoap.SoapClient(serviceURL)
	if err != nil {
		return "", err
	}

	err = soap.Call("bpPayRequest", gosoap.Params{
		"terminalId":     c.TerminalId,
		"userName":       c.UserName,
		"userPassword":   c.Password,
		"callBackUrl":    callBackUrl,
		"orderId":        orderId,
		"amount":         amount,
		"localDate":      localDate.Format(ymd),
		"localTime":      localTime.Format(his),
		"additionalData": additionalData,
		"payerId":        0,
	})
	if err != nil {
		return "", err
	}

	res := new(responseReturn)
	err = soap.Unmarshal(res)
	if err != nil {
		return "", err
	}

	response := strings.Split(res.Return, ",")

	err = c.getBankError(response[0])
	if err != nil {
		return "", err
	}

	return response[1], nil
}

func (c *Client) VerifyRequest(orderId, saleOrderId, saleReferenceId int64, resCode string) error {
	err := c.getBankError(resCode)
	if err == nil {
		soap, err := gosoap.SoapClient(serviceURL)
		if err != nil {
			return err
		}

		params := gosoap.Params{
			"terminalId":      c.TerminalId,
			"userName":        c.UserName,
			"userPassword":    c.Password,
			"orderId":         orderId,
			"saleOrderId":     saleOrderId,
			"saleReferenceId": saleReferenceId,
		}

		err = soap.Call("bpVerifyRequest", params)
		if err != nil {
			return err
		}

		res := new(responseReturn)
		err = soap.Unmarshal(res)
		if err != nil {
			return err
		}

		return c.getBankError(res.Return)
	}

	return err
}

func (c *Client) MakeForm(refId, phoneNumber, body string) string {
	return fmt.Sprintf(
		`<form id="payForm" name="payForm" action="%s" method="post" target="_self">
					<input type="hidden" id="RefId" name="RefId" value='%s' />
					<input type="hidden" id="mobileNo" name="MobileNo" value='%s' />
					<div>%s</div>
					<input type="submit" value="پرداخت" />
				</form>`,
		GatewayURL,
		refId,
		phoneNumber,
		body,
	)
}

type Peeker interface {
	Peek(key string) []byte
}

func (c *Client) ParseCallBack(req Peeker) (SaleReferenceId, RefId, ResCode, SaleOrderId, CardHolderInfo, CardHolderPan string) {
	SaleReferenceId = string(req.Peek("SaleReferenceId"))
	RefId = string(req.Peek("RefId"))
	ResCode = string(req.Peek("ResCode"))
	SaleOrderId = string(req.Peek("SaleOrderId"))
	CardHolderInfo = string(req.Peek("CardHolderInfo"))
	CardHolderPan = string(req.Peek("CardHolderPan"))
	return
}
