package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {
	if strings.Contains(err, "email") {
		return errors.New("Email is already taken")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("Incorrect password")
	}

	if strings.Contains(err, "category_name") {
		return errors.New("Category name is already taken")
	}

	if strings.Contains(err, "bank_name") {
		return errors.New("Bank name is already taken")
	}

	return errors.New("Incorrect details")
}
