package types

// Struct tags must use the syntax `key:"value"` without spaces around the colon.
type Student struct {
	Id    int64
	Name  string `validate:"required"`
	Email string `validate:"required"`
	Age   int    `validate:"required"`
}
