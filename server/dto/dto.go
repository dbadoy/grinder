package dto

import "github.com/dbadoy/grinder/pkg/database"

var (
	_ database.Data = (*Contract)(nil)
	_ database.Data = (*ABI)(nil)

	Indices = []string{Contract{}.Index(), ABI{}.Index()}
)
