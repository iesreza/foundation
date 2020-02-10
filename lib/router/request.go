package router

import (
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
	writer     http.ResponseWriter
	request    *http.Request
	Parameters map[string]value
	Query      url.Values
	Form       url.Values
	Data       map[string]interface{}
	Files      map[string][]FileUpload
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

type FileUpload struct {
	Name        string
	ContentType string
	Size        int64
	Extension   string
	Pointer     *multipart.FileHeader
}

var MaxUploadSize int64 = 32 << 20

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

	if req.Req().Header.Get("Content-Type") == "application/json" {
		return json.NewDecoder(req.Req().Body).Decode(output)
	}
	if req.Req().Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		req.Req().ParseForm()
		return schema.NewDecoder().Decode(output, req.Req().Form)
	}
	if strings.HasPrefix(req.Req().Header.Get("Content-Type"), "multipart/form-data") {
		return schema.NewDecoder().Decode(output, req.Req().MultipartForm.Value)

	}
	return schema.NewDecoder().Decode(output, req.Req().Header)
}

func (req *Request) parseForm() {
	req.Query = req.Req().URL.Query()
	if req.Req().Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		req.Req().ParseForm()
		req.Form = req.Req().Form
	} else if strings.HasPrefix(req.Req().Header.Get("Content-Type"), "multipart/form-data") {

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
