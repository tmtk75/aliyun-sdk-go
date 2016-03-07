package api

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type Config struct {
	RegionId string
}

type Service interface {
	Version() string
	Endpoint(name string) string
	Method(name string) string
}

func Request(conf *Config, srv Service, x interface{}) {
	//fmt.Printf("%v\n", x)
	a := x
	body := requestBody(srv, x)

	name := reflect.ValueOf(a).Type().Name()
	method := srv.Method(name)
	req, err := http.NewRequest(method, srv.Endpoint(name), strings.NewReader(body))
	if err != nil {
		log.Fatalf("%v", err)
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded;charset=UTF-8")

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatalf("%v", err)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("%v", string(b))
}

func Fill(conf *Config, x interface{}) interface{} {
	t := reflect.ValueOf(x).Elem()
	fn := t.FieldByName("RegionId")
	if fn.CanSet() && fn.Kind() == reflect.String && fn.Interface() == "" {
		fn.SetString(conf.RegionId)
	}
	st := t.FieldByName("StartTime")
	et := t.FieldByName("EndTime")
	//log.Printf("StartTime: %v\n", st)
	//log.Printf("EndTime: %v\n", et)
	f := func(v reflect.Value) bool {
		return v.CanSet() && v.Kind() == reflect.Struct && v.Interface() == time.Time{}
	}
	if f(st) && f(et) {
		now := time.Now().UTC()
		st.Set(reflect.ValueOf(now.Add(-1 * time.Hour)))
		et.Set(reflect.ValueOf(now))
	}
	return t.Interface()
}

func requestBody(srv Service, a interface{}) string {
	v := reflect.ValueOf(a)
	t := v.Type()

	q := url.Values{}
	q.Add("Action", t.Name())
	//log.Printf("%v\n", t.Name())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		//log.Printf("%v\n", f.Name)
		x := v.FieldByName(f.Name).Interface()
		if v, ok := x.(string); ok {
			if v != "" {
				q.Add(f.Name, v)
			}
		} else if v, ok := x.(time.Time); ok {
			fmt := f.Tag.Get("format")
			if fmt == "" {
				fmt = "2006-01-02T15:04:03Z"
			}
			q.Add(f.Name, time.Time(v).Format(fmt))
		}
	}

	q.Add("Format", "JSON")
	q.Add("Version", srv.Version())
	q.Add("AccessKeyId", os.Getenv("ALY_ACCESS_KEY_ID"))
	q.Add("SignatureVersion", "1.0")
	q.Add("SignatureMethod", "HMAC-SHA1")
	q.Add("SignatureNonce", strconv.Itoa(int(rand.Float64()*10000000)))
	q.Add("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))

	h := make([]string, 0, len(q))
	for k, _ := range q {
		h = append(h, k)
	}
	sort.Strings(h)

	cqs := ""
	for _, k := range h {
		cqs += "&" + k + "=" + url.QueryEscape(q.Get(k))
	}

	stringToSign := srv.Method(t.Name()) + "&" + "%2F" + "&" + url.QueryEscape(cqs[1:len(cqs)])
	//log.Printf("stringToSign: %v\n", stringToSign)

	secret := os.Getenv("ALY_ACCESS_KEY_SECRET")
	//log.Printf("%v\n", secret)
	hm := hmac.New(sha1.New, []byte(secret+"&"))
	hm.Write([]byte(stringToSign))
	sign := base64.StdEncoding.EncodeToString(hm.Sum(nil))

	return q.Encode() + "&Signature=" + url.QueryEscape(sign)
}
