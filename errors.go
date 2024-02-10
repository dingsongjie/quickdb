package quickdb

import "errors"

var (
	ErrVersionMismatch    = errors.New("version mismatch")
	ErrBucketNameRequired = errors.New("bucket name required")
	ErrKeyRequired        = errors.New("key required")
	ErrKeyTooLarge        = errors.New("key too large")
)
