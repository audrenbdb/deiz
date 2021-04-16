package psql

import (
	"context"
	"fmt"
	"github.com/audrenbdb/deiz"
)

func (r *repo) IsClinicianRegistrationComplete(ctx context.Context, email string) (bool, error) {
	_, err := r.firebaseAuth.GetUserByEmail(ctx, email)
	if err != nil {
		if err.Error() != fmt.Sprintf("cannot find user from email: \"%s\"", email) {
			return false, err
		}
		return false, nil
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
		return deiz.ClinicianAccount{}, fmt.Errorf("unable to get clinician by ID: %s", err)
	}
	acc.Business, err = getBusinessByPersonID(ctx, r.conn, clinicianID)
	if err != nil {
		return deiz.ClinicianAccount{}, fmt.Errorf("unable to get clinician business details: %s", err)
	}
	acc.OfficeAddresses, err = getOfficeAddressesByPersonID(ctx, r.conn, clinicianID)
	if err != nil {
		return deiz.ClinicianAccount{}, fmt.Errorf("unable to get office addresses: %s", err)
	}
	keys, err := getStripeKeysByPersonID(ctx, r.conn, clinicianID)
	if err != nil {
		return deiz.ClinicianAccount{}, fmt.Errorf("unable to get person stripe keys: %s", err)
	}
	acc.StripePublicKey = keys.public
	acc.OfficeHours, err = r.GetClinicianOfficeHours(ctx, clinicianID)
	if err != nil {
		return deiz.ClinicianAccount{}, fmt.Errorf("unable to get clinician office hours: %s", err)
	}
	acc.BookingMotives, err = getBookingMotivesByPersonID(ctx, r.conn, clinicianID)
	if err != nil {
		return deiz.ClinicianAccount{}, fmt.Errorf("unable to get booking motives: %s", err)
	}
	acc.CalendarSettings, err = getCalendarSettingsByPersonID(ctx, r.conn, clinicianID)
	if err != nil {
		return deiz.ClinicianAccount{}, fmt.Errorf("unable to get calendar settings: %s", err)
	}
	acc.PaymentMethods, err = r.GetPaymentMethods(ctx)
	if err != nil {
		return deiz.ClinicianAccount{}, fmt.Errorf("unable to get payment methods: %s", err)
	}
	acc.TaxExemptions, err = r.GetTaxExemptionCodes(ctx)
	if err != nil {
		return deiz.ClinicianAccount{}, fmt.Errorf("unable to get tax exemption codes: %s", err)
	}
	return acc, nil
}

//GetClinicianAccountPublicData retrieves all public available data about clinician
func (r *repo) GetClinicianAccountPublicData(ctx context.Context, clinicianID int) (deiz.ClinicianAccountPublicData, error) {
	var acc deiz.ClinicianAccountPublicData
	var err error
	acc.Clinician, err = getClinicianByID(ctx, r.conn, clinicianID)
	if err != nil {
		return deiz.ClinicianAccountPublicData{}, err
	}
	acc.Clinician.Address = deiz.Address{}
	acc.PublicMotives, err = getPublicMotives(ctx, r.conn, clinicianID)
	if err != nil {
		return deiz.ClinicianAccountPublicData{}, err
	}
	settings, err := getCalendarSettingsByPersonID(ctx, r.conn, clinicianID)
	if err != nil {
		return deiz.ClinicianAccountPublicData{}, err
	}
	acc.ClinicianTz = settings.Timezone.Name
	acc.RemoteAllowed = settings.RemoteAllowed
	keys, err := getStripeKeysByPersonID(ctx, r.conn, clinicianID)
	if err != nil {
		return deiz.ClinicianAccountPublicData{}, err
	}
	acc.StripePublicKey = keys.public
	return acc, nil
}
