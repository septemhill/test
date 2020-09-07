package module

const signUpRandomKey = "sign-up-key-"
const forgetPasswdKey = "forget-passwd-key-"
const resetPasswdKey = "reset-passwd-key-"

func SignupKeyPrefix(key string) string {
	return signUpRandomKey + key
}

func ForgetPasswordKeyPrefix(key string) string {
	return forgetPasswdKey + key
}

func ResetPaswordKeyPrefix(key string) string {
	return resetPasswdKey + key
}
