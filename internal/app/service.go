package app

import (
	"context"
	"net/http"
	"sync"

	"github.com/strick-j/scimplistic/internal/config"
	"github.com/strick-j/scimplistic/internal/web"
	"github.com/strick-j/scimplistic/internal/web/handler"
	"github.com/strick-j/scimplistic/internal/web/views"
	"go.uber.org/zap"
)

type Service struct {
	server *web.Server
	logger *zap.Logger
}

func NewService(baseCtx context.Context, logger *zap.Logger, conn *Connectors, cfg *config.Config) *Service {
	srv := web.NewServer(cfg.Server.ListenParams())

	hWrapper := web.NewWrapper(logger.Named("http"))
	//requireAuth := hWrapper.MiddlewareFunc(middleware.NewAuthMiddleware(authSvc))

	// Serve files for use, omit static from URL
	srv.Router.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	// Handler for initial Index
	srv.Router.HandleFunc("/", views.IndexReq)

	srv.Router.Methods(http.MethodGet).
		Path("/ping").
		HandlerFunc(hWrapper.WrapResourceHandler(handler.Ping))

	// Handler for Settings functions
	settingsRouter := srv.Router.PathPrefix("/settings").Subrouter()
	settingsRouter.HandleFunc("/", views.SettingsHandler)
	settingsRouter.HandleFunc("/general", views.GeneralSettingsHandler)

	// Handlers for user functions
	pamUserRouter := srv.Router.PathPrefix("/users").Subrouter()
	pamUserRouter.HandleFunc("/", views.UsersHandler).Methods("GET")
	pamUserRouter.HandleFunc("/{id}", views.UserHandler).Methods("GET")
	pamUserRouter.HandleFunc("/add", views.UserAddHandler).Methods("POST")
	pamUserRouter.HandleFunc("/update/{id}", views.UserUpdateHandler).Methods("POST")
	pamUserRouter.HandleFunc("/del/{id}", views.UserDelHandler).Methods("POST")

	// Handlers for group functions
	pamGroupRouter := srv.Router.PathPrefix("/groups").Subrouter()
	pamGroupRouter.HandleFunc("/", views.GroupsHandler).Methods("GET")
	pamGroupRouter.HandleFunc("/{id}", views.GroupHandler).Methods("GET")
	pamGroupRouter.HandleFunc("/add", views.GroupAddHandler).Methods("POST")
	pamGroupRouter.HandleFunc("/update/{id}", views.GroupUpdateHandler).Methods("POST")
	pamGroupRouter.HandleFunc("/del/{id}", views.GroupDelHandler).Methods("POST")

	// Handlers for safe functions
	pamSafeRouter := srv.Router.PathPrefix("/safes").Subrouter()
	pamSafeRouter.HandleFunc("/", views.SafesHandler).Methods("GET")
	pamSafeRouter.HandleFunc("/{id}", views.SafeHandler).Methods("GET")
	pamSafeRouter.HandleFunc("/add", views.SafeAddHandler).Methods("POST")
	pamSafeRouter.HandleFunc("/update/{id}", views.SafeUpdateHandler).Methods("POST")
	pamSafeRouter.HandleFunc("/del/{id}", views.SafeDelHandler).Methods("POST")

	// Handlers for safe functions
	pamAccountsRouter := srv.Router.PathPrefix("/accounts").Subrouter()
	pamAccountsRouter.HandleFunc("/", views.AccountsHandler).Methods("GET")
	pamAccountsRouter.HandleFunc("/{id}", views.AccountHandler).Methods("GET")
	pamAccountsRouter.HandleFunc("/add", views.AccountAddHandler).Methods("POST")
	pamAccountsRouter.HandleFunc("/update/{id}", views.AccountUpdateHandler).Methods("POST")
	pamAccountsRouter.HandleFunc("/del/{id}", views.AccountDelHandler).Methods("POST")

	// Auth
	/*authHandler := handler.NewAuthHandler(userSvc, authSvc)
	srv.Router.Methods(http.MethodPost).
		Path("/auth").
		HandlerFunc(hWrapper.WrapResourceHandler(authHandler.Login))
	srv.Router.Methods(http.MethodPost).
		Path("/auth/register").
		HandlerFunc(hWrapper.WrapResourceHandler(authHandler.Register))

	// Session
	sessionRouter := srv.Router.Path("/auth/session").Subrouter()
	sessionRouter.Use(requireAuth)
	sessionRouter.Methods(http.MethodGet).
		HandlerFunc(hWrapper.WrapResourceHandler(authHandler.GetSession))
	sessionRouter.Methods(http.MethodDelete).
		HandlerFunc(hWrapper.WrapHandler(authHandler.Logout))

	// Users
	usrHandler := handler.NewUserHandler(userSvc)
	usrRouter := srv.Router.NewRoute().Subrouter()
	usrRouter.Use(requireAuth)
	usrRouter.Path("/users").Methods(http.MethodGet).
		HandlerFunc(hWrapper.WrapResourceHandler(usrHandler.GetUsersList))
	usrRouter.Path("/users/self").Methods(http.MethodGet).
		HandlerFunc(hWrapper.WrapResourceHandler(usrHandler.GetCurrentUser))
	//usrRouter.Path("/users/self/balance").Methods(http.MethodGet).
	//	HandlerFunc(hWrapper.WrapResourceHandler(usrHandler.GetBalance))
	usrRouter.Path("/users/{userId}").Methods(http.MethodGet).
		HandlerFunc(hWrapper.WrapResourceHandler(usrHandler.GetByID)) */

	return &Service{
		server: srv,
		logger: logger,
	}
}

// Start starts the service
func (s Service) Start(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.logger.Info("starting http service", zap.String("addr", s.server.Addr))
		if err := s.server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				return
			}
			s.logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	go func() {
		<-ctx.Done()
		if err := s.server.Shutdown(ctx); err != nil {
			if err == context.Canceled {
				return
			}
			s.logger.Error("failed to shutdown server", zap.Error(err))
		}
	}()

	wg.Wait()
	s.logger.Info("goodbye")
}
