package dr

type Storage interface {
	Add(dr Dr) error
	List(category string) (error, []Dr)
	Request(category string, id string) (error, Dr)
	HealthCheck() error
}

type Dr struct {
	Category    string
	Description string
	ID          string
	Resource    string
	Reusable    bool
	TTL         int64
}
