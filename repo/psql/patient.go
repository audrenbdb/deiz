package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
)

//IsPatientTiedToClinician checks if a given patient is in clinician patient list
func (r *repo) IsPatientTiedToClinician(ctx context.Context, p *deiz.Patient, clinicianID int) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM patient WHERE clinician_person_id = $1 AND id = $2)`
	var tied bool
	row := r.conn.QueryRow(ctx, query, clinicianID, p.ID)
	if err := row.Scan(&tied); err != nil {
		return false, err
	}
	return tied, nil
}

func (r *repo) GetPatientNotes(ctx context.Context, patientID int) ([]deiz.PatientNote, error) {
	const query = `SELECT id, content FROM patient_note WHERE patient_id = $1`
	rows, err := r.conn.Query(ctx, query, patientID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var notes []deiz.PatientNote
	for rows.Next() {
		var n deiz.PatientNote
		err := rows.Scan(&n.ID, &n.Content)
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

func (r *repo) CreatePatientNote(ctx context.Context, n *deiz.PatientNote, patientID int) error {
	const query = `INSERT INTO patient_note(person_id, content) VALUES($1, $2) RETURNING id`
	row := r.conn.QueryRow(ctx, query, patientID, n.Content)
	return row.Scan(&n.ID)
}

func (r *repo) DeletePatientNote(ctx context.Context, noteID int, patientID int) error {
	const query = `DELETE FROM patient_note WHERE id = $1 AND patient_id = $2`
	cmdTag, err := r.conn.Exec(ctx, query, noteID, patientID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errNothingDeleted
	}
	return nil
}

func (r *repo) CreatePatient(ctx context.Context, p *deiz.Patient, clinicianID int) error {
	const query = `INSERT INTO patient(clinician_person_id, email, name, surname, phone, address_id)
	VALUES($1, $2, $3, $4, $5, NULLIF($6, 0)) RETURNING id`
	row := r.conn.QueryRow(ctx, query, clinicianID, p.Email, p.Name, p.Surname, p.Phone, p.Address.ID)
	return row.Scan(&p.ID)
}

func (r *repo) SearchPatient(ctx context.Context, search string, clinicianID int) ([]deiz.Patient, error) {
	const query = `SELECT p.id, p.email, p.name, p.surname, p.phone,
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
		err := rows.Scan(&p.ID, &p.Email, &p.Name, &p.Surname, &p.Phone, &p.Address.ID,
			&p.Address.Line, &p.Address.PostCode, &p.Address.City, &sml)
		if err != nil {
			return nil, err
		}
		patients = append(patients, p)
	}
	return patients, nil
}

func (r *repo) CountPatients(ctx context.Context, clinicianID int) (int, error) {
	const query = `SELECT count(*) FROM patient WHERE clinician_person_id = $1`
	var count int
	row := r.conn.QueryRow(ctx, query, clinicianID)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repo) UpdatePatient(ctx context.Context, p *deiz.Patient, clinicianID int) error {
	const query = `UPDATE patient SET name = $1, surname = $2, phone = $3, email = $4 WHERE clinician_person_id = $5 AND id = $6`
	cmdTag, err := r.conn.Exec(ctx, query, p.Name, p.Surname, p.Phone, p.Email, clinicianID, p.ID)
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

func (r *repo) CreatePatientAddress(ctx context.Context, a *deiz.Address, patientID int) error {
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

func (r *repo) RemovePatient(ctx context.Context, p *deiz.Patient, clinicianID int) error {
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
