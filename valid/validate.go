/*
Package valid provides simple validation for primitive types
*/
package valid

import "regexp"

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func Email(email string) bool {
	return len(email) >= 3 && len(email) <= 254 && emailRegex.MatchString(email)
}

func Phone(phone string) bool {
	return len(phone) >= 10 && len(phone) <= 20
}
