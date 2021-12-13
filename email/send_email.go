package email

import (
	"gopkg.in/gomail.v2"
	"strings"
	"time"
	"math/rand"
)

//发送邮件
func SendEmail(title, content string, touser string)bool {
	var arr []string
	sub := strings.Index(touser, ";")
	if sub >= 0 {
		spArr := strings.Split(touser, ";")
		for _, v := range spArr {
			arr = append(arr, v)
		}
	}else{
		arr = append(arr, touser)
	}
	m := gomail.NewMessage()
	m.SetHeader("From", "amen017@qq.com")
	m.SetHeader("To", arr...)
	m.SetHeader("Subject", title+" "+GetRandString(16))
	m.SetBody("text/html", content)
	d := gomail.NewDialer("smtp.qq.com", 587, "amen017@qq.com", "fhhyacegtznabidh")
	if err := d.DialAndSend(m); err != nil {
		return false
	}
	return true
}

//随机数种子
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetRandString(size int) string {
	allLetterDigit := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "!", "@", "#", "$", "%", "^", "&", "*"}
	randomSb := ""
	digitSize := len(allLetterDigit)
	for i := 0; i < size; i++ {
		randomSb += allLetterDigit[rnd.Intn(digitSize)]
	}
	return randomSb
}