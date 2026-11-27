package api

type apiError struct {
	Error string `json:"error"`
}

type createEmployeeRequest struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Position string `json:"position"`
}

type employeeResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Position string `json:"position"`
}

type employeeListItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Position string `json:"position"`
}
