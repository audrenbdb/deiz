package address

import (
	"context"
	"github.com/audrenbdb/deiz"
)

func isAddressToClinician(ctx context.Context, addressID int, clinicianID int, getter accountGetter) (bool, error) {
	acc, err := getter.GetClinicianAccount(ctx, clinicianID)
	if err != nil {
		return false, err
	}

	return isAddressIDInAddresses(addressID, append(acc.OfficeAddresses, acc.Clinician.Address)), nil
}

func isAddressIDInAddresses(addressID int, addresses []deiz.Address) bool {
	for _, a := range addresses {
		if a.ID == addressID {
			return true
		}
	}
	return false
}
