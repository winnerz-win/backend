package gmail

import (
	"fmt"
	"net/smtp"
	"strings"

	"txscheduler/brix/tools/dbg"
)

/*
	google 개인정보 - 보안 - 보안 수준이 낮은 앱의 액세스 (활성화)

	brickstream10001@gmail.com / brickstream1!
*/

const (
	_smptAddr = "smtp.gmail.com:587" //587
	_authAddr = "smtp.gmail.com"
)

//Postman : gmail Sender
type Postman struct {
	from string
	pwd  string
}

//New :
func New(from, pwd string) Postman {
	return Postman{
		from: from,
		pwd:  pwd,
	}
}

func (my Postman) send(tolist []string, contentType, subject, message string, fromName ...string) error {

	fromText := my.from
	if len(fromName) > 0 {
		fromText = fromName[0]
		fmt.Println("fromText :", fromText)
	}

	msg := "From: " + fromText + "\n" +
		"To: " + strings.Join(tolist, ",") + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-Version: 1.0" + "\n" +
		`Content-Type: ` + contentType + `; charset="utf-8"` + "\n" +
		`Content-Transfer-Encoding: quoted-printable` + "\n" +
		`Content-Disposition: inline` + "\n\n" +
		message

	err := smtp.SendMail(
		_smptAddr,
		smtp.PlainAuth("", my.from, my.pwd, _authAddr),
		fromText,
		tolist,
		[]byte(msg),
	)
	if err != nil {
		dbg.RedItalicBG("infra.postman@send error : ", my.from)
	}
	return err
}

//SendText :
func (my Postman) SendText(tolist []string, subject, message string, fromName ...string) error {
	return my.send(tolist, "text/plain", subject, message, fromName...)
}

//SendHTML :
func (my Postman) SendHTML(tolist []string, subject, htmlString string) error {
	return my.send(tolist, "text/html", subject, htmlString)
}

//SendMail :
func (my Postman) SendMail(tolist []string, subject, textMessage string) error {
	contentType := "text/plain"
	checkText := strings.ToLower(textMessage)
	htmlCnt := strings.Count(checkText, "html>")
	if htmlCnt >= 2 {
		if strings.Count(checkText, "</html>") >= 1 {
			contentType = "text/html"
		}
	}

	return my.send(tolist, contentType, subject, textMessage)
}
