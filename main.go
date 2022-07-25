package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Welcome to my awesome site</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<h1>Contact Page</h1><p>To get in touch email me at <a href="mailto:me@here.com">mail</a></p>`)
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
