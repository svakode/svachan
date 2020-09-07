package constant

import "github.com/lib/pq"

var (
	DBDuplicateError = pq.ErrorCode("23505")
)
