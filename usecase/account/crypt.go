package account

type (
	StringToBytesCrypter interface {
		CryptStringToBytes(str string) ([]byte, error)
	}
)
