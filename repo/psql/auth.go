package psql

import (
	"context"
	firebaseAuth "firebase.google.com/go/auth"
)

const firebaseUserNotFound = "auth/person-not-found"

func setFirebasePersonClaims(ctx context.Context, auth auth, id, role int, uid string) error {
	claims := map[string]interface{}{
		"personId": id,
		"role":     role,
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
