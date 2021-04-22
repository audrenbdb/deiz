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
	addresses := []deiz.Address{}
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

func (r *Repo) UpdateAddress(ctx context.Context, a *deiz.Address) error {
	return updateAddress(ctx, r.conn, a)
}

func (r *Repo) SetClinicianHomeAddress(ctx context.Context, a *deiz.Address, clinicianID int) error {
	return updatePersonAddress(ctx, r.conn, a, clinicianID)
}

func (r *Repo) CreateClinicianHomeAddress(ctx context.Context, a *deiz.Address, clinicianID int) error {
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

func (r *Repo) CreateClinicianOfficeAddress(ctx context.Context, a *deiz.Address, clinicianID int) error {
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

func (r *Repo) DeleteAddress(ctx context.Context, addressID int) error {
	const query = `DELETE from address WHERE id = $1`
	tag, err := r.conn.Exec(ctx, query, addressID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errNoRowsUpdated
	}
	return nil
}

func (r *Repo) GetAddressByID(ctx context.Context, addressID int) (deiz.Address, error) {
	const query = `SELECT a.line, a.post_code, a.city
	FROM address WHERE id = $1`
	row := r.conn.QueryRow(ctx, query, addressID)
	address := deiz.Address{ID: addressID}
	if err := row.Scan(&address.Line, &address.PostCode, &address.City); err != nil {
		return deiz.Address{}, nil
	}
	return address, nil
}
