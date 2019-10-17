package internal

import (
	// Force a nimport of mock so that it guarantees that the mock package
	// will be included in vendor.
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate go run ./gencodes --destination=../evdevcodes.go
