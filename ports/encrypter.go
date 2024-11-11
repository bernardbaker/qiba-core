package ports

type Encrypter interface {
	EncryptGameData(data interface{}) (string, string, error)
}
