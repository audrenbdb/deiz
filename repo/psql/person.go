package psql

import "context"

type person struct {
	id         int
	role       int
	email      string
	profession string
	addressID  int
	name       string
	surname    string
	phone      string
}

func insertPerson(ctx context.Context, db db, p *person) error {
	const query = `INSERT INTO person(role, profession, address_id, name, surname, phone, email) VALUES($1, $2, NULLIF($3, 0), $4, $5, $6, $7) RETURNING id`
	row := db.QueryRow(ctx, query, p.role, p.profession, p.addressID, p.name, p.surname, p.phone, p.email)
	return row.Scan(&p.id)
}

func updatePersonEmail(ctx context.Context, db db, email string, personID int) error {
	const query = `UPDATE person SET email = $1 WHERE id = $2`
	cmdTag, err := db.Exec(ctx, query, email, personID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errNoRowsUpdated
	}
	return nil
}

func updatePersonRole(ctx context.Context, db db, role int, personID int) error {
	const query = `UPDATE person SET role = $1 WHERE id = $2`
	cmdTag, err := db.Exec(ctx, query, role, personID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errNoRowsUpdated
	}
	return nil
}

func getPersonByID(ctx context.Context, db db, personID int) (person, error) {
	const query = `SELECT role, email FROM person WHERE id = $1`
	row := db.QueryRow(ctx, query, personID)
	p := person{
		id: personID,
	}
	err := row.Scan(&p.role, &p.email)
	if err != nil {
		return person{}, err
	}
	return p, nil
}
