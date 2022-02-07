package stringHelper

import "net/mail"

// returns true if the given string is an email
func CheckEmail(value string) bool {
	_, err := mail.ParseAddress(value)
	return err == nil
}
