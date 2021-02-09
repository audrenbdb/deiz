package deiz

type CryptService struct {
	StringToBytesCrypter StringToBytesCrypter
}

type (
	StringToBytesCrypter interface {
		CryptStringToBytes(str string) ([]byte, error)
	}
)
