package system

import (
	"github.com/iesreza/foundation/lib"
	"github.com/iesreza/foundation/lib/router"
	"time"
)

type User struct {
	Id           int64
	Name         string
	Username     string
	LastActivity int64
	LastSeen     int64
	Guest        bool
	Password     string    `xorm:"varchar(200)"`
	Created      time.Time `xorm:"created"`
	Updated      time.Time `xorm:"updated"`
	OTP          string
	LastOTP      int64
}

func GetUser(req router.Request) *User {
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

func (user *User) EstablishSession(req router.Request) {
	req.Session.Set("user", user)
}

func (user *User) IsGuest() bool {
	return user.Guest
}

func (user *User) SetPassword(password string) {
	user.Password = lib.GeneratePassword(password)
	user.Save()
}

func CheckAuthentication(username, password string) bool {
	/*	Database.Find()*/
	return false
}

func (user *User) Save() {
	Database.Update(user)
}
