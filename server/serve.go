package server

import (
	"fmt"
	"gowiki/model"
	"log"
	"net/http"
	"text/template"
)

var RouteList = []RouteMap{
	{Path: "/edit/", Handler: EditHandler},
	{Path: "/view/", Handler: ViewHandler},
	{Path: "/save/", Handler: SaveHandler},
}

func StartServer() {
	for _, route := range RouteList {
		http.HandleFunc(route.Path, route.Handler)
	}
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := model.LoadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := model.LoadPage(title)
	if err != nil {
		p = &model.Page{Title: title}
	}
	fmt.Fprintf(w, "<h1>Editing %s</h1>"+
		"<form action=\"/save/%s\" method=\"POST\">"+
		"<textarea name=\"body\">%s</textarea><br>"+
		"<input type=\"submit\" value=\"Save\">"+
		"</form>",
		p.Title, p.Title, p.Body)
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &model.Page{Title: title, Body: []byte(body)}
	err := p.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *model.Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}
