package etype

type Error string

// actually a type conversion for string to error that allows us to write
// rejected := etype.Error("rejected")
// and use rejected for error returns
func (e Error) Error() string { return string(e) }
