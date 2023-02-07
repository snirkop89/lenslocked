package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/snirkop89/lenslocked/controllers"
	"github.com/snirkop89/lenslocked/migrations"
	"github.com/snirkop89/lenslocked/models"
	"github.com/snirkop89/lenslocked/templates"
	"github.com/snirkop89/lenslocked/views"
)

func main() {
	// Setup the database
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Setup services
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}

	// Setup middleware
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	csrfKey := "arXLgxHVx3chuf3o2xDckWHj8GYLX4a2"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		//TODO: Fix this before deploying
		csrf.Secure(false),
	)

	// Setup controllers
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	usersC.Templates.New = views.Must(
		views.ParseFS(templates.FS, "signup.html.tmpl", "tailwind.html.tmpl"))
	usersC.Templates.SignIn = views.Must(
		views.ParseFS(templates.FS, "signin.html.tmpl", "tailwind.html.tmpl"))

	// Setup router and routes
	r := chi.NewRouter()
	r.Use(csrfMw)
	r.Use(umw.SetUser)
	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.html.tmpl", "tailwind.html.tmpl"))))

	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.html.tmpl", "tailwind.html.tmpl"))))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.html.tmpl", "tailwind.html.tmpl"))))
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})
	// r.Get("/users/me", usersC.CurrentUser)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Start the server
	fmt.Println("Starting the server on 3000...")
	http.ListenAndServe(":3000", r)
}
