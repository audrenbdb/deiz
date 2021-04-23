package deiz

type Role int32

const (
	PATIENT Role = iota
	CLINICIAN
	ADMIN
)

type Credentials struct {
	UserID int
	Role   Role
}

func (c *Credentials) IsPatient() bool {
	return c.Role == PATIENT
}
