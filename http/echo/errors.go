package echo

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	errReadAuthClaims     Error = "unable to read auth claims"
	errBind               Error = "unable to bind resource into struct"
	errParsingBearerToken Error = "unable to parse bearer token"
	errUnauthorizedRole   Error = "unauthorized to execute this request"
	errValidating         Error = "unable to validate fields"
)
