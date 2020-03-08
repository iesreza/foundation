package request

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/schema"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	writer      http.ResponseWriter
	request     *http.Request
	Parameters  Parameter
	Query       url.Values
	Form        url.Values
	Path        string
	ContentType ContentType
	Data        Data
	Files       File
	Session     Session
	Cookie      Cookie
	Matched     bool
	IsAdmin     bool
	terminated  bool
}
type message struct {
	Message string
	Type    string
}

type Cookie map[string]http.Cookie
type Data map[string]interface{}
type Parameter map[string]value
type File map[string][]FileUpload
type ContentType string

const (
	REQ_JSON             = "application/json"
	REQ_FORM_URL_ENCODED = "application/x-www-form-urlencoded"
	REQ_FORM_MULTI_PART  = "multipart/form-data"
)

type FileUpload struct {
	Name        string
	ContentType string
	Size        int64
	Extension   string
	Pointer     *multipart.FileHeader
}

var MaxUploadSize int64 = 32 << 20

func (c Cookie) Get(key string) *http.Cookie {
	if val, ok := c[key]; ok {
		return &val
	}
	return nil
}
func (c Data) Get(key string) *interface{} {
	if val, ok := c[key]; ok {
		return &val
	}
	return nil
}
func (c File) Get(key string) *[]FileUpload {
	if val, ok := c[key]; ok {
		return &val
	}
	return nil
}
func (c Parameter) Get(key string) *value {
	if val, ok := c[key]; ok {
		return &val
	}
	return nil
}

func (c *Cookie) Set(key string, v http.Cookie) {
	(*c)[key] = v
}
func (c *Data) Set(key string, v interface{}) {
	(*c)[key] = v
}
func (c *File) Set(key string, v []FileUpload) {
	(*c)[key] = v
}
func (c *Parameter) Set(key string, v value) {
	(*c)[key] = v
}

func NewRequest(writer http.ResponseWriter, request *http.Request) Request {
	req := Request{
		writer:  writer,
		request: request,
		Data:    map[string]interface{}{},
	}
	req.parseForm()
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
	r.Terminate()
}

func (r *Request) WriteString(s string) {
	if r.terminated {
		return
	}
	r.writer.Write([]byte(s))
	r.Terminate()
}

func (r *Request) WriteObject(obj interface{}) error {
	if r.terminated {
		return nil
	}
	data, err := json.Marshal(obj)
	if err == nil {
		r.writer.Header().Set("Content-Type", "application/json")
		r.writer.Write(data)
		r.Terminate()
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

func (req *Request) Unmarshal(output interface{}) error {
	//guess data type

	if req.ContentType == REQ_JSON {

		return json.NewDecoder(req.Req().Body).Decode(output)
	}
	if req.ContentType == REQ_FORM_URL_ENCODED {
		req.Req().ParseForm()
		return schema.NewDecoder().Decode(output, req.Req().Form)
	}
	if strings.HasPrefix(string(req.ContentType), REQ_FORM_MULTI_PART) {
		return schema.NewDecoder().Decode(output, req.Req().MultipartForm.Value)

	}
	return schema.NewDecoder().Decode(output, req.Req().Header)
}

func (req *Request) parseForm() {
	req.Query = req.Req().URL.Query()
	req.ContentType = ContentType(req.Req().Header.Get("Content-Type"))
	req.Path = strings.TrimRight(req.request.URL.Path, "/")
	if req.ContentType == REQ_FORM_URL_ENCODED {
		req.Req().ParseForm()
		req.Form = req.Req().Form
	} else if strings.HasPrefix(string(req.ContentType), REQ_FORM_MULTI_PART) {

		req.Req().ParseMultipartForm(MaxUploadSize)
		req.Form = req.Req().MultipartForm.Value
		if len(req.Req().MultipartForm.File) > 0 {
			for key, item := range req.Req().MultipartForm.File {
				files := req.Req().MultipartForm.File[key]
				req.Files[key] = make([]FileUpload, len(item))
				for i, file := range files {
					chunks := strings.Split(file.Filename, ".")
					req.Files[key][i] = FileUpload{
						Name:        file.Filename,
						Extension:   chunks[len(chunks)-1],
						ContentType: file.Header.Get("Content-Type"),
						Size:        file.Size,
						Pointer:     file,
					}

				}

			}

		}

	}
}

func (file FileUpload) Move(dst string) error {
	src, err := file.Pointer.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	// This is path which we want to store the file
	f, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(f, src)

	return nil
}

func (req Request) GetUser() *User {
	user, exist := req.Session.Get("user")
	if exist {
		user.(*User).LastActivity = time.Now().Unix()
		return user.(*User)
	}
	res := &User{
		Guest: true,
	}
	req.Session.Set("user", res)
	return res
}

func (req Request) GetRedirect() string {
	redirect := ""
	if req.Query.Get("redirect") != "" {
		b, err := base64.StdEncoding.DecodeString(req.Query.Get("redirect"))
		if err == nil {
			return string(b)
		}
	}
	redirect = req.request.URL.Path
	if req.request.URL.RawQuery != "" {
		redirect += "?" + req.request.URL.RawQuery
	}

	return redirect
}
