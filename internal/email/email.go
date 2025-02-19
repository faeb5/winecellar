package email

import "net/mail"

func ValidateEmail(email string) (string, error) {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", err
	}
	return addr.Address, nil
}
