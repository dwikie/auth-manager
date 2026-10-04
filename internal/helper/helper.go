package helper

type ISecretProvider interface {
	GetAccessTokenSecret() string
	GetRefreshTokenSecret() string
}

type Helper struct {
	secretProvider ISecretProvider
}

func NewHelper(secretProvider ISecretProvider) *Helper {
	return &Helper{
		secretProvider: secretProvider,
	}
}
