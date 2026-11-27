package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"example.com/crud-example/internal/data"
)

type APIHandler struct {
	employeeRepository data.EmployeeRepository
}

func NewAPIHandler(employeeRepository data.EmployeeRepository) *APIHandler {
	return &APIHandler{
		employeeRepository: employeeRepository,
	}
}

func (s *APIHandler) Handler() http.Handler {
	serveMux := &http.ServeMux{}
	serveMux.HandleFunc("/employees", s.handleEmployees)
	serveMux.HandleFunc("/employees/{id}", s.handleGetEmployee)
	return serveMux
}

func (s *APIHandler) handleEmployees(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetEmployees(w, r)
	case http.MethodPost:
		s.handlePostEmployee(w, r)
	default:
		s.handleError(w, r, http.StatusMethodNotAllowed, "Method "+r.Method+" is not supported")
	}
}

func (s *APIHandler) handleGetEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := s.employeeRepository.GetEmployees()
	if err != nil {
		log.Printf("handleEmployees: error in get call: %v", err)
		s.handleError(w, r, http.StatusInternalServerError, "Could not fetch employees")
		return
	}

	respData := []employeeListItem{}
	for _, employee := range employees {
		respData = append(respData, employeeListItem{
			ID:       strconv.Itoa(employee.ID),
			Name:     employee.Name,
			Age:      employee.Age,
			Position: employee.Position,
		})
	}
	s.writeSuccessJSON(w, respData)
}

func (s *APIHandler) handlePostEmployee(w http.ResponseWriter, r *http.Request) {
	reqData := createEmployeeRequest{}
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		log.Printf("handleEmployee: error in decode: %v", err)
		s.handleError(w, r, http.StatusBadRequest, "Invalid employee data")
		return
	}

	err = reqData.validate()
	if err != nil {
		log.Printf("handleEmployee: validation error: %v", err)
		s.handleError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	exists, err := s.employeeRepository.EmployeeExists(reqData.Name)
	if err != nil {
		log.Printf("handleEmployee: error in exists check: %v", err)
		s.handleError(w, r, http.StatusInternalServerError, "Failed to save employee")
		return
	}
	if exists {
		s.handleError(w, r, http.StatusConflict, "Employee with such name already exists")
		return
	}

	employee := data.Employee{
		Name:     reqData.Name,
		Age:      reqData.Age,
		Position: reqData.Position,
	}
	employeeID, err := s.employeeRepository.SaveEmployee(employee)
	if err != nil {
		log.Printf("handleEmployee: error in save call: %v", err)
		s.handleError(w, r, http.StatusInternalServerError, "Failed to save employee")
		return
	}

	respData := employeeResponse{
		ID:       strconv.Itoa(employeeID),
		Name:     employee.Name,
		Age:      employee.Age,
		Position: employee.Position,
	}
	s.writeSuccessJSON(w, respData)
}

func (s *APIHandler) handleGetEmployee(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		s.handleError(w, r, http.StatusBadRequest, "Invalid ID")
		return
	}

	employee, err := s.employeeRepository.GetEmployeeByID(id)
	if err != nil {
		log.Printf("handleGetEmployee: error in get call: %v", err)
		s.handleError(w, r, http.StatusInternalServerError, "Could not fetch employee")
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

func (s *APIHandler) handleError(w http.ResponseWriter, r *http.Request, statusCode int, msg string) {
	log.Printf("Error in API endpoint: path=%v, statusCode=%d, msg=%v\n", r.URL.Path, statusCode, msg)

	s.writeJSON(w, statusCode, apiError{Error: msg})
}

func (s *APIHandler) writeSuccessJSON(w http.ResponseWriter, v any) error {
	return s.writeJSON(w, 200, v)
}

func (s *APIHandler) writeJSON(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Printf("Error writing JSON response: %v\n", err)
		return fmt.Errorf("writeJSON: %w", err)
	}

	return nil
}
