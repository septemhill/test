package module

const signUpRandomKey = "sign-up-key-"

func SignupKeyPrefix(key string) string {
	return signUpRandomKey + key
}
