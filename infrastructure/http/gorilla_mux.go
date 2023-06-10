package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jeferagudeloc/gorilla-sessions-auth/application/adapter/api/action"
	"github.com/jeferagudeloc/gorilla-sessions-auth/application/adapter/api/middleware"
	"github.com/jeferagudeloc/gorilla-sessions-auth/application/adapter/logger"
	"github.com/jeferagudeloc/gorilla-sessions-auth/application/adapter/repository"
	"github.com/jeferagudeloc/gorilla-sessions-auth/application/usecase"
	"github.com/urfave/negroni"
)

type gorillaMux struct {
	router     *mux.Router
	middleware *negroni.Negroni
	log        logger.Logger
	port       Port
	sqlDB      repository.SQL
	ctxTimeout time.Duration
}

func newGorillaMux(
	log logger.Logger,
	port Port,
	t time.Duration,
	sqlDB repository.SQL,
) *gorillaMux {
	return &gorillaMux{
		router:     mux.NewRouter(),
		middleware: negroni.New(),
		log:        log,
		port:       port,
		ctxTimeout: t,
		sqlDB:      sqlDB,
	}
}

var store = sessions.NewCookieStore([]byte("super-secret-password"))

func (g gorillaMux) Listen() {
	g.setAppHandlers(g.router)
	g.middleware.UseHandler(g.router)

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		Addr:         fmt.Sprintf(":%d", g.port),
		Handler:      g.middleware,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		g.log.WithFields(logger.Fields{"port": g.port}).Infof("Starting HTTP Server On %s", g.port)
		if err := server.ListenAndServe(); err != nil {
			g.log.WithError(err).Fatalln("Error starting HTTP server")
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		g.log.WithError(err).Fatalln("Server Shutdown Failed")
	}

	g.log.Infof("Service down")
}

func (g gorillaMux) setAppHandlers(router *mux.Router) {
	api := router
	enableCORS(api)
	api.HandleFunc("/health", action.HealthCheck).Methods(http.MethodGet)
	api.Handle("/authentication", g.buildAuthAction()).Methods(http.MethodPost)
	api.Handle("/logout", g.buildLogoutAction()).Methods(http.MethodPost)
	api.Handle("/users", g.buildCreateAction()).Methods(http.MethodPost)
}

func (g gorillaMux) buildAuthAction() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			uc = usecase.NewAuthInteractor(
				repository.NewPostgresSQL(g.sqlDB),
			)
			act = action.NewAuthAction(uc, g.log, store)
		)

		act.Execute(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (g gorillaMux) buildLogoutAction() *negroni.Negroni {
	var handler http.HandlerFunc = func(w http.ResponseWriter, req *http.Request) {
		session, err := store.Get(req, "session-name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Options.MaxAge = -1
		err = session.Save(req, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (g gorillaMux) buildCreateAction() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			uc = usecase.NewCreateUserInteractor(
				repository.NewPostgresSQL(g.sqlDB),
			)
			act = action.NewCreateUserAction(uc, g.log, store)
		)

		act.Execute(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}
