//package psql manages persistence of app data
package psql

import (
	"context"
	firebaseAuth "firebase.google.com/go/auth"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type db interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
}

type auth interface {
	GetUserByEmail(ctx context.Context, email string) (*firebaseAuth.UserRecord, error)
	CreateUser(ctx context.Context, create *firebaseAuth.UserToCreate) (*firebaseAuth.UserRecord, error)
	SetCustomUserClaims(ctx context.Context, uid string, customClaims map[string]interface{}) error
	UpdateUser(ctx context.Context, uid string, user *firebaseAuth.UserToUpdate) (*firebaseAuth.UserRecord, error)
}

//ALl functions used by this Repo has to be implemented
type Repo struct {
	conn         *pgxpool.Pool
	firebaseAuth *firebaseAuth.Client
}

func NewRepo(conn *pgxpool.Pool, authClient *firebaseAuth.Client) *Repo {
	return &Repo{
		conn:         conn,
		firebaseAuth: authClient,
	}
}
