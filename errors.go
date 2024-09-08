package snowflake

import "fmt"

// Custom error type for wrapping.
type SnowflakeError struct {
	Code    int
	Message string
}

var (
	ErrorInvalid         = SnowflakeError{0x0, "invalid id"}
	ErrorInvalidByte     = SnowflakeError{0x1, "invalid byte detected"}
	ErrorInvalidJson     = SnowflakeError{0x2, "invalid json format"}
	ErrorEncodeMapLength = SnowflakeError{0x100, "encode map is not long enough"}
)

func (e *SnowflakeError) Error() string {
	return fmt.Sprintf("snowflake ERROR %d: %s", e.Code, e.Message)
}
