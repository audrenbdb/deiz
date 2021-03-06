package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
)

func insertBusiness(ctx context.Context, db db, b *deiz.Business, personID int) error {
	const query = `INSERT INTO business(person_id, identifier, name, tax_exemption_id) VALUES($1, NULLIF($2, ''), NULLIF($3, ''), NULLIF($4, 0)) RETURNING id`
	row := db.QueryRow(ctx, query, personID, b.Identifier, b.Name, b.TaxExemption.ID)
	err := row.Scan(&b.ID)
	if err != nil {
		return err
	}
	return nil
}

func getBusinessByPersonID(ctx context.Context, db db, personID int) (deiz.Business, error) {
	const query = `SELECT b.id, COALESCE(b.name, ''), COALESCE(b.identifier, ''),
	COALESCE(t.id, 0), COALESCE(t.code, ''),
    COALESCE(a.id, 0), COALESCE(a.line, ''), COALESCE(a.post_code, 0), COALESCE(a.city, '')
	FROM business b
	LEFT JOIN tax_exemption t ON b.tax_exemption_id = t.id
    LEFT JOIN address a ON b.address_id = a.id
	WHERE b.person_id = $1`
	row := db.QueryRow(ctx, query, personID)
	b := deiz.Business{}
	err := row.Scan(&b.ID, &b.Name, &b.Identifier, &b.TaxExemption.ID, &b.TaxExemption.Code,
		&b.Address.ID, &b.Address.Line, &b.Address.PostCode, &b.Address.City)
	if err != nil {
		return deiz.Business{}, err
	}
	return b, nil
}

func (r *Repo) GetClinicianBusiness(ctx context.Context, clinicianID int) (deiz.Business, error) {
	return getBusinessByPersonID(ctx, r.conn, clinicianID)
}

func (r *Repo) UpdateClinicianBusiness(ctx context.Context, b *deiz.Business, clinicianID int) error {
	const query = `UPDATE business SET name = $1, identifier = $2, tax_exemption_id = NULLIF($3, 0) WHERE person_id = $4`
	cmdTag, err := r.conn.Exec(ctx, query, b.Name, b.Identifier, b.TaxExemption.ID, clinicianID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errNoRowsUpdated
	}
	return nil
}
