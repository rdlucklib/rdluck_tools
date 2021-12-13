package common

import (
	"crypto/md5"
	cr "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	FormatDate     = "2006-01-02"
	FormatDateTime = "2006-01-02 15:04:05"
	FormatTime     = "15:04:05"
	Regular        = "^((13[0-9])|(14[5|7])|(15([0-3]|[5-9]))|(18[0-9])|(17[0-9]))\\d{8}$"
)

//随机数种子
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetRandString(size int) string {
	allLetterDigit := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	randomSb := ""
	digitSize := len(allLetterDigit)
	for i := 0; i < size; i++ {
		randomSb += allLetterDigit[rnd.Intn(digitSize)]
	}
	return randomSb
}

func StringsToJSON(str string) string {
	rs := []rune(str)
	jsons := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			jsons += string(r)
		} else {
			jsons += "\\u" + strconv.FormatInt(int64(rint), 16) // json
		}
	}
	return jsons
}

//序列化
func ToString(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

//md5加密
func MD5(data string) string {
	m := md5.Sum([]byte(data))
	return hex.EncodeToString(m[:])
}

// 获取数字随机字符
func GetRandDigit(n int) string {
	return fmt.Sprintf("%0"+strconv.Itoa(n)+"d", rnd.Intn(int(math.Pow10(n))))
}

// 获取随机数
func GetRandNumber(n int) int {
	return rnd.Intn(n)
}

func GetToday(format string) string {
	today := time.Now().Format(format)
	return today
}

//获取今天剩余秒数
func GetTodayLastSecond() time.Duration {
	today := GetToday(FormatDate) + " 23:59:59"
	end, _ := time.ParseInLocation(FormatDateTime, today, time.Local)
	return time.Duration(end.Unix()-time.Now().Local().Unix()) * time.Second
}

// 处理出生日期函数
func GetBrithDate(idcard string) string {
	l := len(idcard)
	var s string
	if l == 15 {
		s = "19" + idcard[6:8] + "-" + idcard[8:10] + "-" + idcard[10:12]
		return s
	}
	if l == 18 {
		s = idcard[6:10] + "-" + idcard[10:12] + "-" + idcard[12:14]
		return s
	}
	return GetToday(FormatDate)
}

//处理性别
func WhichSexByIdcard(idcard string) string {
	var sexs = [2]string{"女", "男"}
	length := len(idcard)
	if length == 18 {
		sex, _ := strconv.Atoi(string(idcard[16]))
		return sexs[sex%2]
	} else if length == 15 {
		sex, _ := strconv.Atoi(string(idcard[14]))
		return sexs[sex%2]
	}
	return "男"
}

//截取小数点后几位
func SubFloatToString(f float64, m int) string {
	n := strconv.FormatFloat(f, 'f', -1, 64)
	if n == "" {
		return ""
	}
	if m >= len(n) {
		return n
	}
	newn := strings.Split(n, ".")
	if m == 0 {
		return newn[0]
	}
	if len(newn) < 2 || m >= len(newn[1]) {
		return n
	}
	return newn[0] + "." + newn[1][:m]
}

//截取小数点后几位
func SubFloatToFloat(f float64, m int) float64 {
	newn := SubFloatToString(f, m)
	newf, _ := strconv.ParseFloat(newn, 64)
	return newf
}

//获取相差时间-年
func GetYearDiffer(start_time, end_time string) int64 {
	var Age int64

	t1, err := time.ParseInLocation("2006-01-02", start_time, time.Local)
	t2, err := time.ParseInLocation("2006-01-02", end_time, time.Local)

	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix()
		Age = diff / (3600 * 365 * 24)
		return Age
	} else {
		return Age
	}
}

//获取相差时间-秒
func GetSecondDifferByTime(start_time, end_time time.Time) int64 {
	diff := end_time.Unix() - start_time.Unix()
	return diff
}

func FixFloat(f float64, m int) float64 {
	newn := SubFloatToString(f+0.00000001, m)
	newf, _ := strconv.ParseFloat(newn, 64)
	return newf
}

//验证是否是手机号
func Validate(mobileNum string) bool {
	reg := regexp.MustCompile(Regular)
	return reg.MatchString(mobileNum)
}

func GetNowDayTime(paramTime time.Time) (has bool) {
	nowTime := time.Now().Format("2006-01-02")
	nowDay, _ := time.ParseInLocation("2006-01-02", nowTime, time.Local)
	has = true
	if nowDay.After(paramTime) {
		has = false
	}
	return
}

func UrlEncode(s string) string {
	return url.QueryEscape(s)
}

func VersionToInt(version string) int {
	version = strings.Replace(version, ".", "", -1)
	n, _ := strconv.Atoi(version)
	return n
}

type ToutiaoResult struct {
	Msg  string
	Code int
	Ret  int
}

//sort map by key
func SortMapByKey2Str(m map[string]interface{}) string {
	// To store the keys in slice in sorted order
	var keys []string
	var s string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// To perform the opertion you want
	for _, k := range keys {
		if m[k] != nil {
			s += k + "=" + fmt.Sprint(m[k]) + "&"
		}
	}
	return strings.TrimSuffix(s, "&")
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func Json2Map(jsonstr []byte) (s map[string]string, err error) {
	var result map[string]string
	if err := json.Unmarshal(jsonstr, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// 将字符串数组转化为逗号分割的字符串形式  ["str1","str2","str3"] >>> "str1,str2,str3"
func StrListToString(strList []string) (str string) {
	if len(strList) > 0 {
		for k, v := range strList {
			if k == 0 {
				str = v
			} else {
				str = str + "," + v
			}
		}
		return
	}
	return ""
}

func CalStageTerm(money float64, term, acc int) ([]float64, error) {
	if term < 1 {
		return nil, errors.New("invalid term, term must > 0")
	}
	stages := make([]float64, term)
	otherTerm := FixFloat(money/float64(term), acc)
	firstTerm := FixFloat(money-float64(term-1)*otherTerm, acc)
	for i := 0; i < term; i++ {
		if i == 0 {
			stages[i] = firstTerm
		} else {
			stages[i] = otherTerm
		}
	}
	return stages, nil
}

func RandIntNum(min, max int64) int64 {
	maxBigInt := big.NewInt(max)
	i, _ := cr.Int(cr.Reader, maxBigInt)
	iInt64 := i.Int64()
	if iInt64 < min {
		iInt64 = RandIntNum(min, max) //应该用参数接一下
	}
	return iInt64
}

func GetAgeByBirthdate(Birthdate string) int {
	now := time.Now()
	b, err := time.ParseInLocation("2006-01-02", Birthdate, time.Local)
	fmt.Println(err)
	age := now.Year() - b.Year()
	if now.Month() < b.Month() || (now.Month() == b.Month() && now.Day() < b.Day()) {
		age--
	}
	if age < 0 {
		age = 0
	}
	return age
}
