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
var indexTemplate *template.Template

func InitTemplate() {
	layout, err := path("layout.html")
	if err != nil || layout == "" {
		panic("cant load template")
		return
	}

	loginTemplate = initPageTemplate(layout, "login.html")
	authTemplate = initPageTemplate(layout, "auth.html")
	indexTemplate = initPageTemplate(layout, "index.html")
}

func initPageTemplate(layout string, filename string) *template.Template {
	page, err := path(filename)
	if err != nil || page == "" {
		panic("cant load template")
		return nil
	}
	t, err := template.ParseFiles(layout, page)
	if err != nil {
		panic(err)
		return nil
	}
	log.Printf("load template:%v\n", t)
	return t
}

func path(name string) (path string, err error) {
	layouts, err := filepath.Glob(oauth2Cfg.TemplatePrefix + name)
	if err != nil {
		panic(err)
		return
	}
	if len(layouts) > 0 {
		path = layouts[0]
		log.Printf("read path:%v\n", path)
		return
	}
	return
}

func execLoginTemplate(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	execTemplate(w, r, loginTemplate, "layout", data)
}

func execAuthTemplate(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	execTemplate(w, r, authTemplate, "layout", data)
}

func execIndexTemplate(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	execTemplate(w, r, indexTemplate, "layout", data)
}

func execTemplate(w http.ResponseWriter, r *http.Request, tpl *template.Template, name string, data map[string]interface{}) {
	data["cfg"] = oauth2Cfg

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tpl.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Printf("The template %s exec error:%v\n", name, err)
		errorResponse(w, err, http.StatusInternalServerError)
	}
}
