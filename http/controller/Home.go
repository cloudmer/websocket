package controller

import (
	"net/http"
	"os"
	"time"
	"html/template"
)

type HomeController struct {
	Basecontroller
}

func (home *HomeController) Index()  {
	tNow := time.Now()
	cookie := http.Cookie{Name: "haha", Value: "haha", Expires: tNow.AddDate(1, 0, 0)}
	http.SetCookie(home.Writer, &cookie)
	template.Must(template.ParseFiles(os.Getenv("GOPATH") + "/" + "src/websocket/index.html")).Execute(home.Writer, nil)
}