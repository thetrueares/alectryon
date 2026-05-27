package entities

type NoEntityFound struct {
	Message string
}

func (err NoEntityFound) Error() string {
	return err.Message
}
