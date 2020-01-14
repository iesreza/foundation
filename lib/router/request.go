package router

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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
	IsAdmin    bool
	terminated bool
}
type message struct {
	Message string
	Type    string
}

func NewRequest(writer http.ResponseWriter, request *http.Request) Request {
	req := Request{
		writer:  writer,
		request: request,
	}

	return req
}

func (r *Request) Header(statusCode int) {
	if r.terminated {
		return
	}
	r.writer.WriteHeader(statusCode)
}

func (r *Request) SetCookie(name, value string, expiration time.Duration) {
	if r.terminated {
		return
	}
	cookie := http.Cookie{Name: name, Value: value, Expires: time.Now().Add(expiration)}
	http.SetCookie(r.writer, &cookie)
}

func (r *Request) GetCookie(name string) (*http.Cookie, bool) {
	c, err := r.Req().Cookie(name)
	if err == nil {
		return c, true
	}
	return c, false
}

func (r *Request) RemoveCookie(name string) {
	if r.terminated {
		return
	}
	cookie := http.Cookie{Name: name, Value: "", Expires: time.Now().Add(-1000 * time.Hour)}
	http.SetCookie(r.writer, &cookie)
}

func (r *Request) SetSession(name string, value interface{}) {
	r.Session.Data[name] = value
}

func (r *Request) GetSession(name string) (interface{}, bool) {
	if val, ok := r.Session.Data[name]; ok {
		return val, true
	}
	return nil, false
}

func (r *Request) RemoveSession(name string) {
	delete(r.Session.Data, name)

}

func (r *Request) Write(bytes []byte) {
	if r.terminated {
		return
	}
	r.writer.Write(bytes)
}

func (r *Request) WriteString(s string) {
	if r.terminated {
		return
	}
	r.writer.Write([]byte(s))
}

func (r *Request) WriteObject(obj interface{}) error {
	if r.terminated {
		return nil
	}
	data, err := json.Marshal(obj)
	if err == nil {
		r.writer.Write(data)
		return nil
	} else {
		return err
	}
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

func (req *Request) SetMessage(msg string, msgType string) {
	cookie, exist := req.GetCookie("message")
	var messages []message
	if exist {
		json.Unmarshal([]byte(cookie.Value), &messages)
	} else {
		messages = []message{{msg, msgType}}
	}
	c, err := json.Marshal(messages)
	if err == nil {
		req.SetCookie("message", string(c), 1*time.Hour)
	}
}

func (req *Request) Error(err string) {
	req.SetMessage(err, "error")
}

func (req *Request) Success(message string) {
	req.SetMessage(message, "success")
}

func (req *Request) Warning(warning string) {
	req.SetMessage(warning, "warning")
}

func (req *Request) Notice(notice string) {
	req.SetMessage(notice, "notice")
}

func (req *Request) Route(url string) string {
	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
		return url
	}
	return req.request.URL.Scheme + "://" + req.request.Host + "/" + strings.Trim(url, "/")
}

func (req *Request) Redirect(url string) {
	http.Redirect(req.writer, req.request, req.Route(url), http.StatusSeeOther)
}

func (req *Request) Fail(message string) {
	if message != "" {
		req.Error(message)
	}
	req.Redirect("error")
}

func (req *Request) Terminate() {
	req.terminated = true
}
