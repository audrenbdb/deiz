package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
)

//AddClinicianAccount creates a clinician and its default settings for the application to run
func (r *repo) AddClinicianAccount(ctx context.Context, c *deiz.Clinician) error {
	tx, err := r.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}
	err = insertClinicianPerson(ctx, tx, c)
	if err != nil {
		return err
	}
	err = insertAdeli(ctx, tx, &deiz.Adeli{}, c.ID)
	if err != nil {
		return err
	}
	err = insertCalendarSettings(ctx, tx, &deiz.CalendarSettings{}, c.ID)
	if err != nil {
		return err
	}

	err = insertBusiness(ctx, tx, &deiz.Business{}, c.ID)
	if err != nil {
		return err
	}

	err = insertStripeKeys(ctx, tx, &stripeKeys{}, c.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

//GetClinicianAccount gets all clinician public  required for clinicians / admin client
func (r *repo) GetClinicianAccount(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error) {
	var acc deiz.ClinicianAccount
	var err error
	acc.Clinician, err = getClinicianByID(ctx, r.conn, clinicianID)
	if err != nil {
		return deiz.ClinicianAccount{}, err
	}
	acc.Business, err = getBusinessByPersonID(ctx, r.conn, clinicianID)
	if err != nil {
		return deiz.ClinicianAccount{}, err
	}
	acc.OfficeAddresses, err = getOfficeAddressesByPersonID(ctx, r.conn, clinicianID)
	if err != nil {
		return deiz.ClinicianAccount{}, err
	}
	keys, err := getStripeKeysByPersonID(ctx, r.conn, clinicianID)
	if err != nil {
		return deiz.ClinicianAccount{}, err
	}
	acc.StripePublicKey = keys.public
	acc.OfficeHours, err = getOfficeHoursByPersonID(ctx, r.conn, clinicianID)
	if err != nil {
		return deiz.ClinicianAccount{}, err
	}
	acc.BookingMotives, err = getBookingMotivesByPersonID(ctx, r.conn, clinicianID)
	if err != nil {
		return deiz.ClinicianAccount{}, err
	}
	acc.CalendarSettings, err = getCalendarSettingsByPersonID(ctx, r.conn, clinicianID)
	if err != nil {
		return deiz.ClinicianAccount{}, err
	}
	return acc, nil
}
