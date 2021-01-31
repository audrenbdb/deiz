package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
)

func getClinicianByID(ctx context.Context, db db, clinicianID int) (deiz.Clinician, error) {
	const query = `SELECT p.id, p.name, p.surname, p.email, p.phone, p.profession,
	COALESCE(a.id, 0), COALESCE(a.line, ''), a.post_code, a.city,
	COALESCE(adeli.id, 0), COALESCE(adeli.identifier, '')
	FROM person p
	LEFT JOIN address a ON a.id = p.address_id
	LEFT JOIN adeli adeli ON adeli.person_id = p.id WHERE p.id = $1`
	row := db.QueryRow(ctx, query, clinicianID)
	var c deiz.Clinician
	err := row.Scan(&c.ID, &c.Name, &c.Surname, &c.Email, &c.Phone, &c.Profession,
		&c.Address.ID, &c.Address.Line, &c.Address.PostCode, &c.Address.City,
		&c.Adeli.ID, &c.Adeli.Identifier)
	if err != nil {
		return deiz.Clinician{}, err
	}
	return c, nil
}

func insertClinicianPerson(ctx context.Context, db db, c *deiz.Clinician) error {
	p := &person{
		role:       2,
		email:      c.Email,
		profession: c.Profession,
		addressID:  c.Address.ID,
		name:       c.Name,
		surname:    c.Surname,
		phone:      c.Phone,
	}
	err := insertPerson(ctx, db, p)
	if err != nil {
		return err
	}
	c.ID = p.id
	return nil
}

func insertAdeli(ctx context.Context, db db, a *deiz.Adeli, clinicianID int) error {
	const query = `INSERT INTO adeli(person_id, identifier) VALUES($1, NULLIF($2, '')) RETURNING id`
	row := db.QueryRow(ctx, query, clinicianID, a.Identifier)
	return row.Scan(&a.ID)
}

//Update clinicianPhone
func (r *repo) EditClinicianPhone(ctx context.Context, phone string, clinicianID int) error {
	const query = `UPDATE person SET phone = $1 WHERE id = $2`
	tag, err := r.conn.Exec(ctx, query, phone, clinicianID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errNoRowsUpdated
	}
	return nil
}

//Update clinician email and clinician firebase account email if exists
func (r *repo) EditClinicianEmail(ctx context.Context, newEmail string, clinicianID int) error {
	p, err := getPersonByID(ctx, r.conn, clinicianID)
	if err != nil {
		return err
	}
	tx, err := r.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}
	//update person email
	err = updatePersonEmail(ctx, tx, newEmail, clinicianID)
	if err != nil {
		return err
	}

	//if fire base user is set, also updates his firebase account
	u, err := r.firebaseAuth.GetUserByEmail(ctx, p.email)
	if err != nil {
		if err.Error() != firebaseUserNotFound {
			return err
		}
		return tx.Commit(ctx)
	}
	err = updateFirebaseUserEmail(ctx, r.firebaseAuth, newEmail, u.UID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

//UpdateClinicianRole assign a new role to the given clinician and updates firebase role if user exists
func (r *repo) UpdateClinicianRole(ctx context.Context, role int, clinicianID int) error {
	p, err := getPersonByID(ctx, r.conn, clinicianID)
	if err != nil {
		return err
	}

	tx, err := r.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}
	err = updatePersonRole(ctx, tx, role, clinicianID)
	if err != nil {
		return err
	}

	//if fire base user is set, also updates his firebase claims
	u, err := r.firebaseAuth.GetUserByEmail(ctx, p.email)
	if err != nil {
		if err.Error() != firebaseUserNotFound {
			return err
		}
		return tx.Commit(ctx)
	}
	err = setFirebasePersonClaims(ctx, r.firebaseAuth, person{id: p.id, role: role}, u.UID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
