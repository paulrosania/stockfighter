package api

type ErrorResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"error"`
}

func (r ErrorResponse) Error() string {
	return r.Message
}
