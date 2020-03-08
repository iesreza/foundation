package request

import (
	"github.com/iesreza/foundation/lib"
	"sync"
	"time"
)

var SessionAge = 3600 * time.Second
var sessionPool = map[string]*Session{}
var mu sync.Mutex

type Session struct {
	Data   map[string]interface{}
	Expire int64
	ID     string
}

func (session *Session) Set(key string, value interface{}) {
	session.Data[key] = value
}

func (session *Session) Remove(key string) {
	delete(session.Data, key)
}

func (session *Session) Get(key string) (interface{}, bool) {
	if val, ok := session.Data[key]; ok {
		return val, true
	}
	return nil, false
}

func session(req *Request) {
	mu.Lock()
	cookie, err := req.Req().Cookie("session")
	if err != nil {
		CreateNewSession(req)
	} else {
		if val, ok := sessionPool[cookie.Value]; ok {
			if val.Expire < time.Now().Unix() {
				delete(sessionPool, cookie.Value)
				CreateNewSession(req)
			} else {
				val.Expire = time.Now().Add(SessionAge).Unix()
				req.Session = *val
			}
		} else {
			CreateNewSession(req)
		}
	}
	mu.Unlock()
}

func CreateNewSession(req *Request) {
	request := req.Req()
	rand := lib.RandomString(16, lib.ALPHANUM_SIGNS)
	ua := request.UserAgent()
	ip := request.RemoteAddr
	t := time.Now().Format("Mon, 02 Jan 2006 15:04:05")
	uid := lib.GetMD5Hash(rand + ua + ip + t)
	exp := time.Now().Add(SessionAge)
	s := Session{
		Data:   map[string]interface{}{},
		Expire: exp.Unix(),
		ID:     uid,
	}
	sessionPool[uid] = &s
	req.SetCookie("session", uid, 8640*time.Hour)
	req.Session = s

}

var gc = false

func sessionGC() {
	gc = true
	go func() {
		mu.Lock()
		now := time.Now().Unix()
		for key, item := range sessionPool {
			if item.Expire < now {
				delete(sessionPool, key)
			}
		}
		mu.Unlock()
		time.Sleep(10 * time.Minute)
	}()
}
