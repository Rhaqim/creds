package lib

import "github.com/pquerna/otp/totp"

func Generate() (string, string) {
	opts := totp.GenerateOpts{
		Issuer:      "Creds",
		AccountName: "Creds",
	}
	key, _ := totp.Generate(opts)
	return key.Secret(), key.URL()
}
