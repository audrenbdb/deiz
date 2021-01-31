package psql

import (
	"context"
	firebaseAuth "firebase.google.com/go/auth"
)

const firebaseUserNotFound = "auth/person-not-found"

func setFirebasePersonClaims(ctx context.Context, auth auth, p person, uid string) error {
	claims := map[string]interface{}{
		"personId": p.id,
		"role":     p.role,
	}
	return auth.SetCustomUserClaims(ctx, uid, claims)
}

func updateFirebaseUserEmail(ctx context.Context, auth auth, email, firebaseUserID string) error {
	userToUpdate := (&firebaseAuth.UserToUpdate{}).Email(email)
	_, err := auth.UpdateUser(ctx, firebaseUserID, userToUpdate)
	return err
}

func createFirebaseUser(ctx context.Context, auth auth, email, password string) (*firebaseAuth.UserRecord, error) {
	userToCreate := (&firebaseAuth.UserToCreate{}).Email(email).EmailVerified(false).Password(password).Disabled(false)
	return auth.CreateUser(ctx, userToCreate)
}

//login logs the person
//two scenarios :
//1°) person never logged before. His account is not tied to a firebase account.
//    his firebase account will then be created with new claims (role and ID).
//2°) person has already firebase authentication credentials
//	  claims are then refreshed
func (r *repo) Login(ctx context.Context, email, password string) error {
	u, err := getPersonByEmail(ctx, r.conn, email)
	if err != nil {
		return err
	}
	firebaseUser, err := r.firebaseAuth.GetUserByEmail(ctx, email)
	if err != nil {
		if err.Error() != firebaseUserNotFound {
			return err
		}
		firebaseUser, err = createFirebaseUser(ctx, r.firebaseAuth, email, password)
		if err != nil {
			return err
		}
		return setFirebasePersonClaims(ctx, r.firebaseAuth, u, firebaseUser.UID)
	}
	return nil
}
