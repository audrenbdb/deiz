package deiz

type Role int32

const (
	PublicRole Role = iota
	PatientRole
	ClinicianRole
	AdminRole
)

type Credentials struct {
	UserID int
	Role   Role
}

func (c *Credentials) IsPatient() bool {
	return c.Role == PatientRole
}
