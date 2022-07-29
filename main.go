package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func executeTemplate(w http.ResponseWriter, filepath string) {
	w.Header().Set("Content-Type", "text/html")
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was an error parsing the template", http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, nil)
	if err != nil {
		log.Printf("exeuting template: %v", err)
		http.Error(w, "Error executing the template", http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("templates", "home.html.tmpl")
	executeTemplate(w, tmplPath)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("templates", "contact.html.tmpl")
	executeTemplate(w, tmplPath)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset='utf-8")
	fmt.Fprint(w, `
<h1>FAQ</h1>
<p>
<strong>Q:</strong> Is there a free version?<br>
<strong>A:</strong> Yes! We offer a free trail for 30 days on any paid plans.
</p>

<p>
<strong>Q:</strong> What are you support hours?<br>
<strong>A:</strong> We have support staff answering emails 24/7, tough response times may be
a bit slower on weekends
</p>

<p>
<strong>Q:</strong> How do I contact support?<br>
<strong>A:</strong> Email us - support@lenslocked.com
</p>
	`)
}

func main() {
	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on 3000...")
	http.ListenAndServe(":3000", r)
}
