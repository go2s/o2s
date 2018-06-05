// authors: wangoo
// created: 2018-06-02
// template

package o2

import (
	"net/http"
	"html/template"
	"path/filepath"
	"log"
)

var loginTemplate *template.Template
var authTemplate *template.Template

func InitTemplate() {
	layout := path("layout.html")
	login := path("login.html")
	auth := path("auth.html")

	var err error
	loginTemplate, err = template.ParseFiles(layout, login)
	if err != nil {
		panic(err)
	}
	authTemplate, err = template.ParseFiles(layout, auth)
	if err != nil {
		panic(err)
	}
}

func path(name string) string {
	layouts, err := filepath.Glob(oauth2Cfg.TemplatePrefix + name)
	if err != nil {
		panic(err)
	}
	return layouts[0]
}

func execLoginTemplate(w http.ResponseWriter, r *http.Request, data interface{}) {
	execTemplate(w, r, loginTemplate, "layout", data)
}

func execAuthTemplate(w http.ResponseWriter, r *http.Request, data interface{}) {
	execTemplate(w, r, authTemplate, "layout", data)
}

func execTemplate(w http.ResponseWriter, r *http.Request, tpl *template.Template, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tpl.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Printf("The template %s exec error:%v\n", name, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
