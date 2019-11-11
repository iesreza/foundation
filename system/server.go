package system

import (
	"github.com/iesreza/foundation/lib/router"
	"net/http"
)

var config = GetConfig()
func StartWebServer()  {
	if config.Server.HTTPS{
		go func() {
			if err := http.ListenAndServe(config.Server.Host+":"+config.Server.Port, http.HandlerFunc(redirectTLS)); err != nil {
				Error(err,"Unable to turn on web server on %s",config.Server.Host+":"+config.Server.Port)
			}
		}()
		go func() {
			Router.Middleware(func(req router.Request) bool {
				req.Req().URL.Scheme = "https"
				return true
			})
			err := http.ListenAndServeTLS(config.Server.Host+":443", config.Server.Cert, config.Server.Key, &Router)
			Error(err,"Unable to turn on web server on %s",config.Server.Host+":"+config.Server.Port)
		}()
	}else{
		go func() {
			Router.Middleware(func(req router.Request) bool {
				req.Req().URL.Scheme = "http"
				return true
			})
			err := http.ListenAndServe(config.Server.Host+":"+config.Server.Port, &Router)
			Error(err,"Unable to turn on web server on %s",config.Server.Host+":"+config.Server.Port)
		}()
	}
}



func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}
