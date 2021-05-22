package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/jackc/pgx/v4"
)

//IsPatientTiedToClinician checks if a given patient is in clinician patient list
func (r *Repo) IsPatientTiedToClinician(ctx context.Context, p *deiz.Patient, clinicianID int) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM patient WHERE clinician_person_id = $1 AND id = $2)`
	var tied bool
	row := r.conn.QueryRow(ctx, query, clinicianID, p.ID)
	if err := row.Scan(&tied); err != nil {
		return false, err
	}
	return tied, nil
}

func (r *Repo) CreatePatient(ctx context.Context, p *deiz.Patient, clinicianID int) error {
	const query = `INSERT INTO patient(clinician_person_id, email, name, surname, phone, address_id)
	VALUES($1, NULLIF($2, ''), $3, $4, $5, NULLIF($6, 0)) RETURNING id`
	row := r.conn.QueryRow(ctx, query, clinicianID, p.Email, p.Name, p.Surname, p.Phone, p.Address.ID)
	return row.Scan(&p.ID)
}

func (r *Repo) GetPatientByEmail(ctx context.Context, email string, clinicianID int) (deiz.Patient, error) {
	const query = `SELECT id, name, surname, phone, COALESCE(email, '') FROM patient WHERE clinician_person_id = $1 AND email = $2`
	row := r.conn.QueryRow(ctx, query, clinicianID, email)
	var p deiz.Patient
	err := row.Scan(&p.ID, &p.Name, &p.Surname, &p.Phone, &p.Email)
	if err != nil && err != pgx.ErrNoRows {
		return deiz.Patient{}, err
	}
	return p, nil
}

func (r *Repo) SearchPatient(ctx context.Context, search string, clinicianID int) ([]deiz.Patient, error) {
	const query = `SELECT p.id, COALESCE(p.email, ''), p.name, p.surname, p.phone, COALESCE(p.note, ''),
		COALESCE(a.id, 0) address_id, COALESCE(a.line, '') address_line, COALESCE(a.post_code, 0) address_post_code, COALESCE(a.city, '') address_city,
		similarity(p.name, $1) AS name_sml
		FROM patient p LEFT JOIN address a ON p.address_id = a.id
		WHERE p.name % $1 AND p.clinician_person_id = $2
		ORDER BY name_sml DESC LIMIT 5`
	rows, err := r.conn.Query(ctx, query, search, clinicianID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var patients []deiz.Patient
	for rows.Next() {
		var p deiz.Patient
		var sml float64
		err := rows.Scan(&p.ID, &p.Email, &p.Name, &p.Surname, &p.Phone, &p.Note, &p.Address.ID,
			&p.Address.Line, &p.Address.PostCode, &p.Address.City, &sml)
		if err != nil {
			return nil, err
		}
		patients = append(patients, p)
	}
	return patients, nil
}

func (r *Repo) CountPatients(ctx context.Context, clinicianID int) (int, error) {
	const query = `SELECT count(*) FROM patient WHERE clinician_person_id = $1`
	var count int
	row := r.conn.QueryRow(ctx, query, clinicianID)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repo) UpdatePatient(ctx context.Context, p *deiz.Patient, clinicianID int) error {
	const query = `UPDATE patient SET name = $1, surname = $2, phone = $3, email = NULLIF($4, ''), note = NULLIF($5, '') WHERE clinician_person_id = $6 AND id = $7`
	cmdTag, err := r.conn.Exec(ctx, query, p.Name, p.Surname, p.Phone, p.Email, p.Note, clinicianID, p.ID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errNoRowsUpdated
	}
	return nil
}

func updatePatientAddress(ctx context.Context, db db, addressID int, patientID int) error {
	const query = `UPDATE patient SET address_id = $1 WHERE id = $2`
	cmdTag, err := db.Exec(ctx, query, addressID, patientID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errNoRowsUpdated
	}
	return nil
}

func (r *Repo) CreatePatientAddress(ctx context.Context, a *deiz.Address, patientID int) error {
	tx, err := r.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}
	if err := insertAddress(ctx, tx, a); err != nil {
		return err
	}
	if err := updatePatientAddress(ctx, tx, a.ID, patientID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *Repo) RemovePatient(ctx context.Context, p *deiz.Patient, clinicianID int) error {
	const query = `DELETE from patient WHERE clinician_person_id = $1 AND id = $2`
	cmdTag, err := r.conn.Exec(ctx, query, clinicianID, p.ID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errNothingDeleted
	}
	return nil
}
