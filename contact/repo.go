package contact

import (
	"context"
	"github.com/audrenbdb/deiz/psql"
	"github.com/jackc/pgx/v4"
)

type clinician struct {
	name    string
	surname string
	email   string
}

type getClinicianByID = func(ctx context.Context, clinicianID int) (clinician, error)

type psqlRepo struct {
	db psql.PGX
}

func (r psqlRepo) createGetClinicianByIDFunc() getClinicianByID {
	return func(ctx context.Context, clinicianID int) (clinician, error) {
		return r.queryClinician(ctx, "WHERE id = $1", clinicianID)
	}
}

func (r psqlRepo) queryClinician(ctx context.Context, queryConditions string, args ...interface{}) (clinician, error) {
	const query = `SELECT name, surname, email FROM person p`
	return r.scanClinicianRow(r.db.QueryRow(ctx, query+` `+queryConditions, args...))
}

func (r psqlRepo) scanClinicianRow(row pgx.Row) (c clinician, err error) {
	return c, row.Scan(&c.name, &c.surname, &c.email)
}
