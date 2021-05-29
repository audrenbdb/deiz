package contact

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type psqlDB interface {
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
}

type clinician struct {
	name    string
	surname string
	email   string
}

type getClinicianByID = func(ctx context.Context, clinicianID int) (clinician, error)

func psqlQueryClinician(ctx context.Context, db psqlDB, queryConditions string, args ...interface{}) (clinician, error) {
	const query = `SELECT name, surname, email FROM person p`
	return psqlScanClinicianRow(db.QueryRow(ctx, query+` `+queryConditions, args...))
}

func psqlScanClinicianRow(row pgx.Row) (c clinician, err error) {
	return c, row.Scan(&c.name, &c.surname, &c.email)
}

func psqlGetClinicianByID(db psqlDB) getClinicianByID {
	return func(ctx context.Context, clinicianID int) (clinician, error) {
		return psqlQueryClinician(ctx, db, "WHERE id = $1", clinicianID)
	}
}
