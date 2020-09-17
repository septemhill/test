package module

const sessionKey = "session-token-key-"
const signUpRandomKey = "sign-up-key-"
const forgetPasswdKey = "forget-passwd-key-"
const resetPasswdKey = "reset-passwd-key-"

func SessionTokenPrefix(key string) string {
	return sessionKey + key
}

func SignupKeyPrefix(key string) string {
	return signUpRandomKey + key
}

func ForgetPasswordKeyPrefix(key string) string {
	return forgetPasswdKey + key
}

func ResetPaswordKeyPrefix(key string) string {
	return resetPasswdKey + key
}
