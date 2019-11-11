package router

import (
	"net/http"
	"strconv"
	"time"
)

type Request struct {
	writer     http.ResponseWriter
	request    *http.Request
	Parameters map[string]value
	Get        map[string]value
	Post       map[string]value
	Session    Session
	Cookie     map[string]http.Cookie
	Matched    bool
}

func NewRequest(writer http.ResponseWriter, request *http.Request) Request {
	req := Request{
		writer:  writer,
		request: request,
	}

	return req
}

func (r *Request) Header(statusCode int) {
	r.writer.WriteHeader(statusCode)
}

func (r *Request) SetCookie(name,value string,expiration time.Duration) {
	cookie := http.Cookie{Name: name, Value: value, Expires: time.Now().Add(expiration)}
	http.SetCookie(r.writer, &cookie)
}

func (r *Request) GetCookie(name string) (*http.Cookie,bool) {
	c,err := r.Req().Cookie(name)
	if err == nil{
		return c,true
	}
	return c,false
}

func (r *Request) RemoveCookie(name string) {
	cookie := http.Cookie{Name: name, Value: "", Expires: time.Now().Add(-1000*time.Hour)}
	http.SetCookie(r.writer, &cookie)
}



func (r *Request) SetSession(name string,value interface{}) {
	r.Session.Data[name] = value
}

func (r *Request) GetSession(name string) (interface{},bool) {
	if val, ok := 	r.Session.Data[name]; ok {
		return val,true
	}
	return nil,false
}

func (r *Request) RemoveSession(name string) {
	delete(r.Session.Data,name)

}


func (r *Request) Write(bytes []byte) {
	r.writer.Write(bytes)
}

func (r *Request) WriteString(s string) {
	r.writer.Write([]byte(s))
}

type value string

func (r *Request) Req() *http.Request {
	return r.request
}

func (r *Request) Writer() *http.ResponseWriter {
	return &r.writer
}

func (v value) Int() int {
	i, _ := strconv.Atoi(string(v))
	return i
}

func (v value) Float() float64 {
	f, _ := strconv.ParseFloat(string(v), 32)
	return f
}

func (v value) String() string {
	return string(v)
}
