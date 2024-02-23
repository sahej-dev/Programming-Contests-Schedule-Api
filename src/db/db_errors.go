package db

type DbDoesNotExist struct{}

func (e DbDoesNotExist) Error() string {
	return "Database file does not exists"
}
