package cmd

import (
	"context"
	"it-test/adapters/psql"
	"it-test/app"
	"it-test/app/query"
	"it-test/config"
	"it-test/pkg/server"
	"it-test/ports"

	"github.com/go-chi/chi/v5"

	"net/http"

	"github.com/go-pg/pg/v10"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	port string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "starts a http server",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		loader := configLoader()
		cfg := &config.AppConfig{}
		check(loader.LoadConfig(".", cfg))
		initLogLevel(cfg)
		initGocore(cfg)
		startHTTP(cmd.Context(), cfg)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(
		&port,
		"port",
		"p",
		"8080",
		"HTTP Server port to listen on.\nExample: it-test server -p 8000",
	)
}

func startHTTP(rootCtx context.Context, cfg *config.AppConfig) {
	logrus.WithContext(rootCtx).
		WithField("port", port).
		Info("Starting a http server")

	db := createPostgresConnection(cfg)
	defer func() {
		logrus.WithContext(rootCtx).Info("Closing db connections")
		err := db.Close()
		if err != nil {
			logrus.WithContext(rootCtx).WithError(err).Info("Failed to close db connections")
		} else {
			logrus.WithContext(rootCtx).Info("Closed db connections successfully")
		}
	}()

	deps := &Dependencies{
		DB: db,
	}

	application := newApplication(rootCtx, cfg, deps)
	server.RunHTTPServer(rootCtx, cfg, port, func(router chi.Router) http.Handler {
		return ports.HandlerFromMux(ports.NewHTTPServer(application), router)
	})
}

func newApplication(ctx context.Context, cfg *config.AppConfig, deps *Dependencies) *app.Application {
	logrus.WithContext(ctx).WithField("config", cfg).Info("Creating application")

	userRepository := psql.NewUserPSQLRepository(deps.DB)

	return &app.Application{
		Commands: &app.Commands{},
		Queries: &app.Queries{
			GetUserCount: query.NewGetUserCountHandler(userRepository),
			InsertUser:   query.NewCreateUserHandler(userRepository),
			UpdateUser:   query.NewUpdateUserHandler(userRepository),
			ListUser:     query.NewListUserHandler(userRepository),
		},
		AppConfig: cfg,
	}
}

func createPostgresConnection(cfg *config.AppConfig) *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     cfg.DatabaseURL,
		User:     cfg.DatabaseUsername,
		Password: cfg.DatabasePassword,
		Database: cfg.DatabaseDB,
	})
	if cfg.DatabasePrintQueries {
		db.AddQueryHook(psqlQueryHook{})
	}
	return db
}

type psqlQueryHook struct {
}

func (p psqlQueryHook) BeforeQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	formattedQuery, _ := event.FormattedQuery()
	logrus.WithField("before query", string(formattedQuery)).Debug()
	return ctx, nil
}

func (p psqlQueryHook) AfterQuery(ctx context.Context, event *pg.QueryEvent) error {
	formattedQuery, _ := event.FormattedQuery()
	err := event.Err
	logrus.WithField("after query", string(formattedQuery)).WithError(err).Debug()
	return nil
}
