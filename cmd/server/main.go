package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/VaneZ444/forum-service/internal/handler"
	"github.com/VaneZ444/forum-service/internal/repository/postgres"
	"github.com/VaneZ444/forum-service/internal/usecase"
	ssov1 "github.com/VaneZ444/golang-forum-protos/gen/go/forum"
)

func main() {
	addr := ":50051"
	dsn := "postgres://postgres:3781@localhost:5432/forum_db?sslmode=disable"

	// Logger
	logger := slog.New(slog.NewJSONHandler(log.Writer(), nil))
	logger.Info("starting forum-service")

	// DB Connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("failed to connect to DB", slog.String("err", err.Error()))
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Error("failed to ping DB", slog.String("err", err.Error()))
		return
	}

	// Apply migrations
	if err := applyMigrations(db); err != nil {
		logger.Error("failed to apply migrations", slog.String("err", err.Error()))
		return
	}

	// Repositories
	categoryRepo := postgres.NewCategoryRepository(db)
	topicRepo := postgres.NewTopicRepository(db)
	commentRepo := postgres.NewCommentRepository(db)
	postRepo := postgres.NewPostRepository(db)
	tagRepo := postgres.NewTagRepo(db)

	// UseCases
	categoryUC := usecase.NewCategoryUseCase(categoryRepo, logger)
	topicUC := usecase.NewTopicUseCase(topicRepo, categoryRepo, logger)
	commentUC := usecase.NewCommentUseCase(commentRepo, postRepo, logger)
	postUC := usecase.NewPostUseCase(postRepo, topicRepo, tagRepo, logger)
	tagUC := usecase.NewTagUseCase(tagRepo, postRepo, logger)

	// Handlers
	forumHandler := handler.NewForumHandler(categoryUC, topicUC, postUC, commentUC, tagUC, logger)

	// gRPC server
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("failed to listen", slog.String("err", err.Error()))
		return
	}
	grpcServer := grpc.NewServer()

	ssov1.RegisterForumServiceServer(grpcServer, forumHandler)

	logger.Info("forum-service is listening", slog.String("addr", addr))
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("failed to serve", slog.String("err", err.Error()))
	}
}
func applyMigrations(db *sql.DB) error {
	driver, err := migratepg.WithInstance(db, &migratepg.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../internal/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
