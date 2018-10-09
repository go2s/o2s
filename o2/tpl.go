// authors: wangoo
// created: 2018-06-02
// template

package o2

import (
	"html/template"
	"net/http"

	"github.com/go2s/o2s/tpl"
	"github.com/golang/glog"
)

var loginTemplate *template.Template
var authTemplate *template.Template
var indexTemplate *template.Template

func parse(name, content string) *template.Template {
	t, err := template.New(name).Parse(content)
	if err != nil || t == nil {
		glog.Fatalf("can't load template %v", name)
	}
	return t

}
func templateParse(name, content string) *template.Template {
	layout := parse("layout", tpl.Files["layout.html"])
	t, err := layout.New(name).Parse(content)
	if err != nil || t == nil {
		glog.Fatalf("can't load template %v", name)
	}
	return t
}

//InitTemplate initial tempalte
func InitTemplate() {
	loginTemplate = templateParse("login", tpl.Files["login.html"])
	authTemplate = templateParse("auth", tpl.Files["auth.html"])
	indexTemplate = templateParse("index", tpl.Files["index.html"])

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
		glog.Infof("The template %s exec error:%v", name, err)
		ErrorResponse(w, err, http.StatusInternalServerError)
	}
}
