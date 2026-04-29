package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"example.com/crud-example/internal/dao"
	_ "github.com/lib/pq"
)

type APIServer struct {
	dao *dao.DAO
}

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

func NewAPIServer(dao *dao.DAO) APIServer {
	return APIServer{
		dao: dao,
	}
}

func (s *APIServer) Run() {
	const serverAddress = ":8080"
	log.Println("Starting server on " + serverAddress)

	http.HandleFunc("/employees", s.handleEmployees)
	http.HandleFunc("/employees/{id}", s.handleGetEmployee)

	err := http.ListenAndServe(serverAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *APIServer) handleEmployees(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetEmployees(w, r)
	case http.MethodPost:
		s.handlePostEmployee(w, r)
	default:
		s.handleError(w, r, http.StatusMethodNotAllowed, "Method "+r.Method+" is not supported")
	}
}

func (s *APIServer) handleGetEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := s.dao.GetEmployees()
	if err != nil {
		log.Printf("handleEmployees: error in get call: %v", err)
		s.handleError(w, r, http.StatusInternalServerError, "Could not fetch employees")
		return
	}

	resData := []employeeListItem{}
	for _, employee := range employees {
		resData = append(resData, employeeListItem{
			ID:       strconv.Itoa(employee.ID),
			Name:     employee.Name,
			Age:      employee.Age,
			Position: employee.Position,
		})
	}
	s.writeSuccessJSON(w, resData)
}

func (s *APIServer) handlePostEmployee(w http.ResponseWriter, r *http.Request) {
	reqData := createEmployeeRequest{}
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		log.Printf("handleEmployee: error in decode: %v", err)
		s.handleError(w, r, http.StatusInternalServerError, "Failed to save employee")
		return
	}

	exists, err := s.dao.EmployeeExists(reqData.Name)
	if err != nil {
		log.Printf("handleEmployee: error in exists check: %v", err)
		s.handleError(w, r, http.StatusInternalServerError, "Failed to save employee")
		return
	}
	if exists {
		s.handleError(w, r, http.StatusConflict, "Employee with this name already exists")
		return
	}

	employee := dao.Employee{
		Name:     reqData.Name,
		Age:      reqData.Age,
		Position: reqData.Position,
	}
	employeeID, err := s.dao.SaveEmployee(employee)
	if err != nil {
		log.Printf("handleEmployee: error in save call: %v", err)
		s.handleError(w, r, http.StatusInternalServerError, "Failed to save employee")
		return
	}

	resData := employeeResponse{
		ID:       strconv.Itoa(employeeID),
		Name:     employee.Name,
		Age:      employee.Age,
		Position: employee.Position,
	}
	s.writeSuccessJSON(w, resData)
}

func (s *APIServer) handleGetEmployee(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		s.handleError(w, r, http.StatusBadRequest, "Invalid ID")
		return
	}

	employee, err := s.dao.GetEmployeeByID(id)
	if err != nil {
		log.Printf("handleGetEmployee: error in get call: %v", err)
		s.handleError(w, r, http.StatusInternalServerError, "Could not fetch this employee")
		return
	}
	if employee == nil {
		s.handleError(w, r, http.StatusNotFound, "Employee not found")
		return
	}

	resData := employeeListItem{
		ID:       strconv.Itoa(employee.ID),
		Name:     employee.Name,
		Age:      employee.Age,
		Position: employee.Position,
	}
	s.writeSuccessJSON(w, resData)
}

func (s *APIServer) handleError(w http.ResponseWriter, r *http.Request, statusCode int, msg string) {
	log.Printf("Error in API endpoint: path=%v, statusCode=%d, msg=%v\n", r.URL.Path, statusCode, msg)

	s.writeJSON(w, statusCode, apiError{Error: msg})
}

func (s *APIServer) writeSuccessJSON(w http.ResponseWriter, v any) error {
	return s.writeJSON(w, 200, v)
}

func (s *APIServer) writeJSON(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Printf("Error writing JSON response: %v\n", err)
		return fmt.Errorf("writeJSON: %w", err)
	}

	return nil
}
