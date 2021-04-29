package deiz

type Role int32

const (
	PUBLIC Role = iota
	PATIENT
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
