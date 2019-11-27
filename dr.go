package dr

type Storage interface {
	Add(dr *Dr) AddResult
	List(category string) ListResult
	Request(category string, id string) RequestResult
}

type Dr struct {
	Category    string
	Description string
	Id          string
	Resource    string
	Reusable    bool
	ValidUntil  int64
}

type AddResult struct {
	Success bool
	Error   string
}

type ListResult struct {
	Success   bool
	Error     string
	Resources []Dr
}

type RequestResult struct {
	Success  bool
	Error    string
	Resource Dr
}
