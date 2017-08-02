package utorrent

import "fmt"

//ClientError ...
type ClientError struct {
	code    int
	message string
}

func (ce ClientError) Error() string {
	return fmt.Sprintf("%d: %s", ce.code, ce.message)
}
