package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
)

func getOfficeHoursByPersonID(ctx context.Context, db db, personID int) ([]deiz.OfficeHours, error) {
	const query = `SELECT h.id, h.start_mn, h.end_mn, h.week_day,
	a.id, a.line, a.post_code, a.city
	FROM office_hours h
	INNER JOIN address a ON h.address_id = a.id
	WHERE h.person_id = $1`
	rows, err := db.Query(ctx, query, personID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var hours []deiz.OfficeHours
	for rows.Next() {
		var h deiz.OfficeHours
		err := rows.Scan(&h.ID, &h.StartMn, &h.EndMn, &h.WeekDay, &h.Address.ID, &h.Address.Line, &h.Address.PostCode, &h.Address.City)
		if err != nil {
			return nil, err
		}
		hours = append(hours, h)
	}
	return hours, nil
}

func (r *repo) GetOfficeHours(ctx context.Context, clinicianID int) ([]deiz.OfficeHours, error) {
	return getOfficeHoursByPersonID(ctx, r.conn, clinicianID)
}

func (r *repo) AddOfficeHours(ctx context.Context, h *deiz.OfficeHours, clinicianID int) error {
	const query = `INSERT INTO office_hours(start_mn, end_mn, week_day, address_id, person_id)
	VALUES($1, $2, $3, NULLIF($4, 0), $5) RETURNING id`
	row := r.conn.QueryRow(ctx, query, h.StartMn, h.EndMn, h.WeekDay, h.Address.ID, clinicianID)
	return row.Scan(&h.ID)
}

func (r *repo) RemoveOfficeHours(ctx context.Context, h *deiz.OfficeHours, clinicianID int) error {
	const query = `DELETE FROM office_hours WHERE id = $1 AND person_id = $2`
	cmdTag, err := r.conn.Exec(ctx, query, h.ID, clinicianID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errNothingDeleted
	}
	return nil
}
