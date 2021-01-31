package psql

type Error string

func (e Error) Error() string {
	return string(e)
}

const errNoRowsUpdated Error = "no rows updated"
const errNothingDeleted Error = "nothing deleted"
const errUnauthorized Error = "unauthorized"
const errNoRowsCreated Error = "no rows created"
