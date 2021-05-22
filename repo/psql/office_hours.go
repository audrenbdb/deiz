package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
)

func (r *Repo) GetClinicianOfficeHours(ctx context.Context, clinicianID int) ([]deiz.OfficeHours, error) {
	const query = `SELECT h.id, h.start_mn, h.end_mn, h.week_day, h.availability_id,
	COALESCE(a.id, 0), COALESCE(a.line, ''), COALESCE(a.post_code, 0), COALESCE(a.city, '')
	FROM office_hours h
	LEFT JOIN address a ON h.address_id = a.id
	WHERE h.person_id = $1`
	rows, err := r.conn.Query(ctx, query, clinicianID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	hours := []deiz.OfficeHours{}
	for rows.Next() {
		var h deiz.OfficeHours
		err := rows.Scan(&h.ID, &h.StartMn, &h.EndMn, &h.WeekDay, &h.AvailabilityType,
			&h.Address.ID, &h.Address.Line, &h.Address.PostCode, &h.Address.City)
		if err != nil {
			return nil, err
		}
		hours = append(hours, h)
	}
	return hours, nil
}

func (r *Repo) CreateOfficeHours(ctx context.Context, h *deiz.OfficeHours, clinicianID int) error {
	const query = `INSERT INTO office_hours(start_mn, end_mn, week_day, address_id, person_id, availability_id)
	VALUES($1, $2, $3, NULLIF($4, 0), $5, $6) RETURNING id`
	row := r.conn.QueryRow(ctx, query, h.StartMn, h.EndMn, h.WeekDay, h.Address.ID, clinicianID, h.AvailabilityType)
	return row.Scan(&h.ID)
}

func (r *Repo) DeleteOfficeHours(ctx context.Context, hoursID, clinicianID int) error {
	const query = `DELETE FROM office_hours WHERE id = $1 AND person_id = $2`
	cmdTag, err := r.conn.Exec(ctx, query, hoursID, clinicianID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errNothingDeleted
	}
	return nil
}
