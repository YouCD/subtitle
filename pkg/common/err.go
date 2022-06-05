package common

import "github.com/pkg/errors"

var (
	InvalidCharacter = errors.New("invalid character 'ÿ' looking for beginning of value")
)
