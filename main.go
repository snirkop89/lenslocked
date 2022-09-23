package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/snirkop89/lenslocked/controllers"
	"github.com/snirkop89/lenslocked/models"
	"github.com/snirkop89/lenslocked/templates"
	"github.com/snirkop89/lenslocked/views"
)

func main() {
	r := chi.NewRouter()
	// parse the templates
	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.html.tmpl", "tailwind.html.tmpl"))))

	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.html.tmpl", "tailwind.html.tmpl"))))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.html.tmpl", "tailwind.html.tmpl"))))

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userService := models.UserService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService: &userService,
	}

	usersC.Templates.New = views.Must(
		views.ParseFS(templates.FS, "signup.html.tmpl", "tailwind.html.tmpl"))
	usersC.Templates.SignIn = views.Must(
		views.ParseFS(templates.FS, "signin.html.tmpl", "tailwind.html.tmpl"))
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Get("/users/me", usersC.CurrentUser)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on 3000...")
	http.ListenAndServe(":3000", r)
}
