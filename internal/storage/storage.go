package storage

type Storage interface {
	Save (id string, url string) error
	Get (id string) (string, bool, error)
}