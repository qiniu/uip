package db

import "errors"

var (
	ErrUnsupportedFormat = errors.New("unsupported db format")
	ErrDatabaseError     = errors.New("database error")

	ErrFileSize = errors.New("IP Database file size error")
	ErrMetaData = errors.New("IP Database metadata error")
	ErrReadFull = errors.New("IP Database ReadFull error")

	ErrIPFormat = errors.New("query IP format error")

	ErrNoSupportLanguage = errors.New("language not support")
	ErrNoSupportIPv4     = errors.New("IPv4 not support")
	ErrNoSupportIPv6     = errors.New("IPv6 not support")

	ErrDataNotExists = errors.New("data is not exists")

	ErrIPVersionNotSupported = errors.New("IP version not supported")
	ErrInvalidFieldsLength   = errors.New("invalid fields length")
	ErrInvalidCIDR           = errors.New("invalid CIDR")
	ErrCIDRConflict          = errors.New("CIDR conflict")
	ErrMetaNotFound          = errors.New("meta not found")
	ErrDBFormatNotSupported  = errors.New("format not supported")
	ErrEmptyFile             = errors.New("empty file")
	ErrDatabaseIsInvalid     = errors.New("database is invalid")

	ErrCheckFailed = errors.New("check ip failed")
)
