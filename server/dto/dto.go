package dto

import "github.com/dbadoy/grinder/pkg/database"

var (
	_ = database.Data(&Contract{})

	Indices = []string{Contract{}.Index()}
)
