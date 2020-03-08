package admin

import (
	"encoding/base64"
	"github.com/iesreza/foundation/lib/request"
)

type Controller struct {
}

type AuthController struct{}

func (c AuthController) login(req request.Request) {
	username := req.Form.Get("username")
	password := req.Form.Get("password")
	if req.ContentType == request.REQ_JSON {
		data := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}

		req.Unmarshal(&data)
		username = data.Username
		password = data.Password
	}

	res, err := request.Login(username, password)
	if err == nil {
		res.EstablishSession(req)
		redirect := "/admin/dashboard"
		if req.Query.Get("redirect") != "" {
			b, err := base64.StdEncoding.DecodeString(req.Query.Get("redirect"))
			if err == nil {
				redirect = string(b)
			}
		}
		req.Response(true, "You have logged in", nil, redirect)
	} else {
		req.Response(false, err.Error(), nil, req.GetRedirect())
	}
}
