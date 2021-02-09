package deiz_test

type (
	mockCrypter struct {
		bytes []byte
		err   error
	}
)

func (r *mockCrypter) CryptStringToBytes(str string) ([]byte, error) {
	return r.bytes, r.err
}
