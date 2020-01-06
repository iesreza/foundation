package system

import (
	"github.com/iesreza/foundation/lib/log"
	"github.com/iesreza/foundation/lib/router"
	"net/http"
)

var config = GetConfig()
var Shutdown = false

func StartWebServer() {

	if config.Server.HTTPS {
		go func() {
			if err := http.ListenAndServe(config.Server.Host+":"+config.Server.Port, http.HandlerFunc(redirectTLS)); err != nil {
				log.Error(err, "Unable to turn on web server on %s", config.Server.Host+":"+config.Server.Port)
			}
		}()
		go func() {
			Router.Middleware(func(req router.Request) bool {
				req.Req().URL.Scheme = "https"
				return true
			})
			err := http.ListenAndServeTLS(config.Server.Host+":443", config.Server.Cert, config.Server.Key, &Router)
			log.Error(err, "Unable to turn on web server on %s", config.Server.Host+":"+config.Server.Port)
		}()
	} else {
		go func() {
			Router.Middleware(func(req router.Request) bool {
				req.Req().URL.Scheme = "http"
				return true
			})
			err := http.ListenAndServe(config.Server.Host+":"+config.Server.Port, &Router)
			log.Error(err, "Unable to turn on web server on %s", config.Server.Host+":"+config.Server.Port)
		}()
	}

	Router.Middleware(func(req router.Request) bool {
		if Shutdown {

			http.Error(*req.Writer(), "Server Not Listening", http.StatusNoContent)

			return false
		}
		return true
	})

	RegisterCLI("server.stop", &struct{}{}, func(command string, data interface{}) {
		Shutdown = true
		log.Info("Server has stopped ...")
	}, "Stop HTTP server")

	RegisterCLI("server.start", &struct{}{}, func(command string, data interface{}) {
		Shutdown = false
		log.Info("Server has started ...")
	}, "Start HTTP server")
}

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}
