package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sober-api/internal/helper"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Health() map[string]string
	CreateOnBoardingFlow(*helper.OnBoardingRequest) error
	CreateAccountFlow(ac *helper.CreateAccountRequest) error
}

type service struct {
	db *sql.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")

	database_uri = os.Getenv("DATABASE_URL")
)

func New() Service {
	// connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	connStr := fmt.Sprint(database_uri)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	s := &service{
		db: db,
	}

	if err := s.init(); err != nil {
		log.Fatalf("Error Creating Database Tables: %+v\n", err)
	}

	return s
}

// NOTE: Creation of Database tables

func (s *service) init() error {

	createTableQueries := []string{
		` create table if not exists users(
	id serial primary key,
	username varchar(150) not null,
	email varchar(200) not null,
	password varchar(200) not null,
	created_at timestamp with time zone default current_timestamp not null

		)`,

		` create table if not exists onboarding(
	id serial primary key,
	user_id int not null, 
	reason text not null,
	sober_date varchar(50) not null,
	created_at timestamp with time zone default current_timestamp not null,

		constraint fk_user
			foreign key (user_id)
			references users(id)
			on delete cascade
		)`,
	}

	for _, query := range createTableQueries {
		if _, err := s.db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) CreateOnBoardingFlow(*helper.OnBoardingRequest) error {

	return nil
}

func (s *service) CreateAccountFlow(ac *helper.CreateAccountRequest) error {
	return nil
}
