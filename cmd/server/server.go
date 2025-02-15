package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"github.com/snirkop89/lenslocked/controllers"
	"github.com/snirkop89/lenslocked/migrations"
	"github.com/snirkop89/lenslocked/models"
	"github.com/snirkop89/lenslocked/templates"
	"github.com/snirkop89/lenslocked/views"
	"golang.org/x/oauth2"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
	OAuthProviders map[string]*oauth2.Config
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	// TODO: Read PSQL values from env
	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}
	if cfg.PSQL.Host == "" && cfg.PSQL.Port == "" {
		return cfg, fmt.Errorf("no PSQL config provided")
	}

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure = os.Getenv("CSRF_SECURE") == "true"

	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")

	cfg.OAuthProviders = make(map[string]*oauth2.Config)
	dbxConfig := &oauth2.Config{
		ClientID:     os.Getenv("DROPBOX_APP_ID"),
		ClientSecret: os.Getenv("DROPBOX_APP_SECRET"),
		Scopes:       []string{"files.metadata.read", "files.content.read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.dropbox.com/oauth2/authorize",
			TokenURL: "https://api.dropboxapi.com/oauth2/token",
		},
	}
	cfg.OAuthProviders["dropbox"] = dbxConfig

	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	err = run(cfg)
	if err != nil {
		panic(err)
	}
}

func run(cfg config) error {
	// Setup the database
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		return err
	}

	// Setup services
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	pwResetService := &models.PasswordResetService{
		DB: db,
	}
	galleryService := &models.GalleryService{
		DB: db,
	}
	emailService := models.NewEmailService(cfg.SMTP)

	// Setup middleware
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"),
	)

	// Setup controllers
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}
	usersC.Templates.New = views.Must(
		views.ParseFS(templates.FS, "signup.html.tmpl", "tailwind.html.tmpl"))
	usersC.Templates.SignIn = views.Must(
		views.ParseFS(templates.FS, "signin.html.tmpl", "tailwind.html.tmpl"))
	usersC.Templates.ForgotPassword = views.Must(
		views.ParseFS(templates.FS, "forgot-pw.html.tmpl", "tailwind.html.tmpl"))
	usersC.Templates.CheckYourEmail = views.Must(
		views.ParseFS(templates.FS, "check-your-email.html.tmpl", "tailwind.html.tmpl"))
	usersC.Templates.ResetPassword = views.Must(
		views.ParseFS(templates.FS, "reset-pw.html.tmpl", "tailwind.html.tmpl"))

	galleriesC := controllers.Galleries{
		GalleryService: galleryService,
	}
	galleriesC.Templates.New = views.Must(
		views.ParseFS(templates.FS, "galleries/new.html.tmpl", "tailwind.html.tmpl"))
	galleriesC.Templates.Edit = views.Must(
		views.ParseFS(templates.FS, "galleries/edit.html.tmpl", "tailwind.html.tmpl"))
	galleriesC.Templates.Index = views.Must(
		views.ParseFS(templates.FS, "galleries/index.html.tmpl", "tailwind.html.tmpl"))
	galleriesC.Templates.Show = views.Must(
		views.ParseFS(templates.FS, "galleries/show.html.tmpl", "tailwind.html.tmpl"))

	oauthC := controllers.OAuth{
		ProviderConfigs: cfg.OAuthProviders,
	}

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
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})
	r.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", galleriesC.Show)
		r.Get("/{id}/images/{filename}", galleriesC.Image)

		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/new", galleriesC.New)
			r.Get("/", galleriesC.Index)
			r.Post("/", galleriesC.Create)
			r.Get("/{id}/edit", galleriesC.Edit)
			r.Post("/{id}", galleriesC.Update)
			r.Post("/{id}/delete", galleriesC.Delete)
			r.Post("/{id}/images", galleriesC.UploadImage)
			r.Post("/{id}/images/{filename}/delete", galleriesC.DeleteImage)
		})
	})
	r.Route("/oauth/{provider}", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/connect", oauthC.Connect)
	})

	assetsHandler := http.FileServer(http.Dir("assets"))
	r.Get("/assets/*", http.StripPrefix("/assets", assetsHandler).ServeHTTP)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Start the server
	fmt.Printf("Starting the server on %s...\n", cfg.Server.Address)
	return http.ListenAndServe(cfg.Server.Address, r)
}
