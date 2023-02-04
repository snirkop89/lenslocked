package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
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

	err = models.Migrate(db, "migrations")
	if err != nil {
		panic(err)
	}

	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}

	usersC.Templates.New = views.Must(
		views.ParseFS(templates.FS, "signup.html.tmpl", "tailwind.html.tmpl"))
	usersC.Templates.SignIn = views.Must(
		views.ParseFS(templates.FS, "signin.html.tmpl", "tailwind.html.tmpl"))
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Get("/users/me", usersC.CurrentUser)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on 3000...")

	csrfKey := "arXLgxHVx3chuf3o2xDckWHj8GYLX4a2"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		//TODO: Fix this before deploying
		csrf.Secure(false),
	)
	http.ListenAndServe(":3000", csrfMw(r))
}
