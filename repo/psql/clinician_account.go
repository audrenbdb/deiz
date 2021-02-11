package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
)

func (r *repo) IsClinicianRegistrationComplete(ctx context.Context, email string) (bool, error) {
	_, err := r.firebaseAuth.GetUserByEmail(ctx, email)
	if err != nil {
		if err.Error() != firebaseUserNotFound {
			return false, err
		}
	}
	return true, nil
}

func (r *repo) CompleteClinicianRegistration(ctx context.Context, clinician *deiz.Clinician, password string, clinicianID int) error {
	firebaseUser, err := createFirebaseUser(ctx, r.firebaseAuth, clinician.Email, password)
	if err != nil {
		return err
	}
	return setFirebasePersonClaims(ctx, r.firebaseAuth, clinician.ID, 2, firebaseUser.UID)
}

//AddClinicianAccount creates a clinician and its default settings for the application to run
func (r *repo) CreateClinicianAccount(ctx context.Context, acc *deiz.ClinicianAccount) error {
	tx, err := r.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}
	err = insertClinicianPerson(ctx, tx, &acc.Clinician)
	if err != nil {
		return err
	}
	err = insertAdeli(ctx, tx, &acc.Clinician.Adeli, acc.Clinician.ID)
	if err != nil {
		return err
	}
	err = insertCalendarSettings(ctx, tx, &acc.CalendarSettings, acc.Clinician.ID)
	if err != nil {
		return err
	}

	err = insertBusiness(ctx, tx, &acc.Business, acc.Clinician.ID)
	if err != nil {
		return err
	}

	err = insertStripeKeys(ctx, tx, &stripeKeys{}, acc.Clinician.ID)
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
	acc.OfficeHours, err = r.GetClinicianOfficeHours(ctx, clinicianID)
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

func (r *repo) UpdateClinicianStripeKeys(ctx context.Context, pk string, sk []byte, clinicianID int) error {
	k := stripeKeys{
		public: pk,
		secret: sk,
	}
	return updatePersonStripeKeys(ctx, r.conn, k, clinicianID)
}
