package psql

import "context"

func (r *repo) UpdateClinicianAdeli(ctx context.Context, identifier string, clinicianID int) error {
	const query = `UPDATE adeli SET identifier = $1 WHERE person_id = $2`
	tag, err := r.conn.Exec(ctx, query, identifier, clinicianID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errNoRowsUpdated
	}
	return nil
}
