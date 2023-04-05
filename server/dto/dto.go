package dto

import "github.com/dbadoy/grinder/pkg/database"

var (
	_ database.Data = (*Contract)(nil)

	Indices = []string{Contract{}.Index()}
)
