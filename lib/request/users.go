package request

import (
	"fmt"
	"github.com/iesreza/foundation/lib"
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

func Login(username, password string) (*User, error) {
	if username == "test" && password == "test" {
		u := User{
			Username: username,
			Name:     username,
			Guest:    false,
		}
		return &u, nil
	}

	return &User{Guest: true}, fmt.Errorf("invalid username and password")
}

func GetUser(req Request) *User {
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

func (user *User) EstablishSession(req Request) {
	req.Session.Set("user", user)
}

func (user *User) IsGuest() bool {
	return user.Guest
}

func (user *User) SetPassword(password string) {
	user.Password = lib.GeneratePassword(password)
	//user.Save()
}

func CheckAuthentication(username, password string) bool {
	/*	Database.Find()*/
	return false
}

/*func (user *User) Save() {
	Database.Update(user)
}*/
