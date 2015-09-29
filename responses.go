package statuscake

type autheticationErrorResponse struct {
	ErrNo int
	Error string
}

type updateResponse struct {
	Issues   map[string]string `json:"Issues"`
	Success  bool              `json:"Success"`
	Message  string            `json:"Message"`
	InsertID int               `json:"InsertID"`
}

type deleteResponse struct {
	Success bool   `json:"Success"`
	Error   string `json:"Error"`
}
