package main

import (
	"excalidraw-complete/handlers/api/documents"
	"excalidraw-complete/handlers/api/firebase"
	"excalidraw-complete/handlers/api/kv"
	"excalidraw-complete/handlers/api/openai"
	"excalidraw-complete/handlers/auth"
	authMiddleware "excalidraw-complete/middleware"
	"excalidraw-complete/stores"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func setupRouter(store stores.Store) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Content-Length", "X-CSRF-Token", "Token", "session", "Origin", "Host", "Connection", "Accept-Encoding", "Accept-Language", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Route("/v1/projects/{project_id}/databases/{database_id}", func(r chi.Router) {
		r.Post("/documents:commit", firebase.HandleBatchCommit())
		r.Post("/documents:batchGet", firebase.HandleBatchGet())
	})

	r.Route("/api/v2", func(r chi.Router) {
		// Route for canvases, protected by JWT auth
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.AuthJWT)
			r.Route("/kv", func(r chi.Router) {
				r.Get("/", kv.HandleListCanvases(store))
				r.Route("/{key}", func(r chi.Router) {
					r.Get("/", kv.HandleGetCanvas(store))
					r.Put("/", kv.HandleSaveCanvas(store))
					r.Delete("/", kv.HandleDeleteCanvas(store))
				})
			})
			r.Route("/chat", func(r chi.Router) {
				r.Post("/completions", openai.HandleChatCompletion())
			})
		})

		// Old routes for anonymous document sharing
		r.Post("/post/", documents.HandleCreate(store))
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", documents.HandleGet(store))
		})
	})

	r.Route("/auth", func(r chi.Router) {
		r.Get("/login", auth.HandleLogin)
		r.Get("/callback", auth.HandleCallback)
	})

	return r
}

func waitForShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	logrus.Info("shutting down")
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		logrus.Info("No .env file found")
	}

	listenAddress := flag.String("listen", ":3002", "The address to listen on.")
	logLevel := flag.String("loglevel", "info", "The log level (debug, info, warn, error).")
	flag.Parse()

	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		logrus.Fatalf("Invalid log level: %v", err)
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	auth.InitAuth()
	openai.Init()
	store := stores.GetStore()

	r := setupRouter(store)

	logrus.WithField("addr", *listenAddress).Info("starting backend server")
	if err := http.ListenAndServe(*listenAddress, r); err != nil {
	logrus.WithField("event", "start server").Fatal(err)
	}
}
