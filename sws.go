package main

import (
	"os"
	"github.com/jessevdk/go-flags"
	"github.com/bestmethod/go-logger"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"fmt"
	"crypto/subtle"
	"strings"
)

type config struct {
	BindAddress string `short:"b" long:"bind-address" default:":8080" description:"address to bind to. E.g. ':80' or '127.0.0.1:8000'"`
	PrintAccessLog bool `short:"l" long:"print-accesslog" description:"print each file access log as it's requested"`
	AuthUser string `short:"u" long:"auth-user" description:"if set, will enable basic HTTP auth, requiring this user"`
	AuthPass string `short:"p" long:"auth-pass" description:"password to set for the basic HTTP auth user"`
	log *Logger.Logger
	pwd string
	fileServer http.Handler
}

func main() {
	Conf := config{}
	Conf.log = new(Logger.Logger)
	Conf.log.Init("","SWS",Logger.LEVEL_DEBUG|Logger.LEVEL_WARN|Logger.LEVEL_INFO,Logger.LEVEL_CRITICAL|Logger.LEVEL_ERROR,Logger.LEVEL_NONE)
	p := flags.NewParser(&Conf, flags.Default^flags.PrintErrors)
	tail, err := p.ParseArgs(os.Args)
	if err != nil {
		Conf.log.Fatalf(1,"%s",strings.Replace(err.Error(),"[OPTIONS]","[OPTIONS] [directoryToServe]",-1))
	}
	if len(tail) > 2 {
		Conf.log.Fatalf(1,"switch error")
	}
	if len(tail) == 2 {
		var dirstat os.FileInfo
		if dirstat, err = os.Stat(tail[1]); err != nil {
			Conf.log.Fatalf(1,"directory does not exist or is not accessible: %s, err: %s",tail[1],err)
		}
		if dirstat.IsDir() == false {
			Conf.log.Fatalf(1,"not a directory: %s",tail[1])
		}
		os.Chdir(tail[1])
		Conf.pwd, _ = os.Getwd()
		Conf.log.Info("Serving files from: %s",Conf.pwd)
	}
	Conf.ListenAndServe()
}

func (c *config) ListenAndServe() {
	router := httprouter.New()
	router.GET("/*path", c.ServeHTTP)
	router.POST("/*path", c.ServeHTTP)
	c.fileServer = http.FileServer(http.Dir("."))
	c.log.Info("Starting webserver on %s",c.BindAddress)
	http.ListenAndServe(c.BindAddress,router)
}

func (c *config) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if c.AuthUser != "" {
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(c.AuthUser)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(c.AuthPass)) != 1 {
			if c.PrintAccessLog == true {
				c.log.Info(fmt.Sprintf("path=%s client=%s x-forwarded-for=%s method=%s AuthRequired", r.URL, r.RemoteAddr, r.Header.Get("X-Forwarded-For"), r.Method))
			}
			w.Header().Set("WWW-Authenticate", `Basic realm="SWS: Auth required"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}
	}
	if c.PrintAccessLog == true {
		c.log.Info(fmt.Sprintf("path=%s client=%s x-forwarded-for=%s method=%s", r.URL, r.RemoteAddr, r.Header.Get("X-Forwarded-For"), r.Method))
	}
	w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	w.Header().Add("Pragma", "no-cache") // HTTP 1.0.
	w.Header().Add("Expires", "0") // Proxies.
	c.fileServer.ServeHTTP(w,r)
}
