<div dir="rtl">

# کتابخانه اتصال به سامانه پرداخت اینترنتی [ به پرداخت ملت](http://www.behpardakht.com/) 

روش کار:

<div dir="ltr">

```go

bpClient := mellatBank.Client{
	TerminalId: "TerminalId",
	UserName:   "UserName",
	Password:   "Password",
}

    // payment request
refId, err := bpClient.PaymentRequest(
	orderId,
	amount,
	localDate,
	localTime,
	additionalData,
	callBackUrl,
)
if err != nil {
	return err
}
	
// use of mellatBank.GatewayURL + refId to redirect user to payment page
    
// verify request
err := bpClient.VerifyRequest(
	orderId,
	saleOrderId,
	saleReferenceId,
	resCode,
)
if err != nil {
	return err
}

```

</div>

در صورت خطا یا مشکلی حتما اعلام بفرمایید

<div dir="ltr">

`github.com/tiaguinho/gosoap`

**{SOAP package for Go}**

</div>

</div>

