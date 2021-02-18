package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
)

func getOfficeAddressesByPersonID(ctx context.Context, db db, personID int) ([]deiz.Address, error) {
	const query = `SELECT a.id, a.line, a.post_code, a.city
	FROM office_address o INNER JOIN address a ON o.address_id = a.id
	WHERE o.person_id = $1`
	rows, err := db.Query(ctx, query, personID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var addresses []deiz.Address
	for rows.Next() {
		var a deiz.Address
		err := rows.Scan(&a.ID, &a.Line, &a.PostCode, &a.City)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, a)
	}
	return addresses, nil
}

func insertAddress(ctx context.Context, db db, a *deiz.Address) error {
	const query = `INSERT INTO address(line, post_code, city) VALUES($1, $2, $3) RETURNING id`
	row := db.QueryRow(ctx, query, a.Line, a.PostCode, a.City)
	err := row.Scan(&a.ID)
	if err != nil {
		return err
	}
	return nil
}

func updatePersonAddress(ctx context.Context, db db, a *deiz.Address, personID int) error {
	const query = `UPDATE person SET address_id = $1 WHERE id = $2`
	tag, err := db.Exec(ctx, query, a.ID, personID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errNoRowsUpdated
	}
	return nil
}

func insertOfficeAddress(ctx context.Context, db db, a *deiz.Address, personID int) error {
	const query = `INSERT INTO office_address(person_id, address_id) VALUES($1, $2)`
	cmdTag, err := db.Exec(ctx, query, personID, a.ID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errNoRowsCreated
	}
	return nil
}

func updateAddress(ctx context.Context, db db, a *deiz.Address) error {
	const query = `UPDATE address SET line = $1, post_code = $2, city = $3 WHERE id = $4`
	tag, err := db.Exec(ctx, query, a.Line, a.PostCode, a.City, a.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errNoRowsUpdated
	}
	return nil
}

func (r *repo) UpdateAddress(ctx context.Context, a *deiz.Address) error {
	return updateAddress(ctx, r.conn, a)
}

func (r *repo) SetClinicianHomeAddress(ctx context.Context, a *deiz.Address, clinicianID int) error {
	return updatePersonAddress(ctx, r.conn, a, clinicianID)
}

func (r *repo) CreateClinicianHomeAddress(ctx context.Context, a *deiz.Address, clinicianID int) error {
	tx, err := r.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}
	err = insertAddress(ctx, tx, a)
	if err != nil {
		return err
	}
	err = updatePersonAddress(ctx, tx, a, clinicianID)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) IsAddressToClinician(ctx context.Context, a *deiz.Address, clinicianID int) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM address a
			LEFT JOIN person p ON p.address_id = a.id LEFT JOIN office_address o ON o.address_id = a.id 
			WHERE (p.id = $1 OR o.person_id = $1) AND a.id = $2)`
	var owns bool
	row := r.conn.QueryRow(ctx, query, clinicianID, a.ID)
	if err := row.Scan(&owns); err != nil {
		return false, err
	}
	return owns, nil
}

func (r *repo) CreateClinicianOfficeAddress(ctx context.Context, a *deiz.Address, clinicianID int) error {
	tx, err := r.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}
	err = insertAddress(ctx, tx, a)
	if err != nil {
		return err
	}
	err = insertOfficeAddress(ctx, tx, a, clinicianID)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
