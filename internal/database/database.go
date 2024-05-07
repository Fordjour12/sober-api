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
	CreateAccountFlow(*helper.CreateAccountRequest) (int, error)
	CreateNotesFlow(*helper.CreateNotesRequest) (int, error)
	GetUserByEmail(string) (*helper.CreateAccountRequest, error)
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
)

func New() Service {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)

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

		` create table if not exists notes(
	id serial primary key,
	user_id int not null,
	content text not null,
	created_at timestamp with time zone default current_timestamp not null,
	update_at timestamp with time zone default current_timestamp not null,

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

func (s *service) CreateOnBoardingFlow(ob *helper.OnBoardingRequest) error {

	fmt.Printf("Onboarding CreateOnBoardingFlow: %+v\n", ob)

	query := `insert into onboarding(user_id, reason, sober_date,created_at) values($1, $2, $3, $4)`

	_, err := s.db.Exec(
		query,
		ob.UserId,
		ob.Sobriety.ReasonForJoining,
		ob.Sobriety.SoberDate,
		ob.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *service) CreateAccountFlow(ac *helper.CreateAccountRequest) (int, error) {

	fmt.Printf("Account CreateAccountFlow: %+v\n", ac)

	query := `insert into users(username, email, password, created_at) values ($1, $2, $3, $4) returning id`

	var id int
	err := s.db.QueryRow(
		query,
		ac.Username,
		ac.Email,
		ac.Password,
		ac.CreatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, err
}

func (s *service) CreateNotesFlow(cn *helper.CreateNotesRequest) (int, error) {
	fmt.Printf("Notes CreateNotesFlow: %+v\n", cn)

	query := `insert into notes(user_id, content, created_at, update_at) values ($1, $2, $3, $4) returning id`

	var id int
	err := s.db.QueryRow(
		query,
		cn.UserId,
		cn.Content,
		cn.CreatedAt,
		cn.UpdateAt,
	).Scan(&id)

	if err != nil {
		log.Fatalf("Error Creating Notes: %+v\n", err)
		return 0, err
	}

	return id, nil
}

func (s *service) GetUserByEmail(email string) (*helper.CreateAccountRequest, error) {
	query := `select username, email, password,created_at from users where email = $1`

	fmt.Println("email", email)

	ac := &helper.CreateAccountRequest{}
	err := s.db.QueryRow(query, email).Scan(
		&ac.Username,
		&ac.Email,
		&ac.Password,
		&ac.CreatedAt,
	)

	if err != nil {
		log.Fatalf("Error Getting User By Email: %+v\n", err)
		return nil, err
	}

	return ac, nil
}
