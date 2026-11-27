package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/crud-example/internal/data"
)

type TestAPIClient struct {
	baseUrl string
	client  *http.Client
}

func (c *TestAPIClient) Get(url string, respData any) error {
	resp, err := c.client.Get(c.baseUrl + url)
	if err != nil {
		return err
	}
	err = json.NewDecoder(resp.Body).Decode(respData)
	if err != nil {
		return err
	}
	return nil
}

func (c *TestAPIClient) Post(url string, reqData, respData any) error {
	var reqDataBuf bytes.Buffer
	err := json.NewEncoder(&reqDataBuf).Encode(reqData)
	if err != nil {
		return err
	}
	resp, err := c.client.Post(c.baseUrl+url, "application/json", &reqDataBuf)
	if err != nil {
		return err
	}
	return json.NewDecoder(resp.Body).Decode(respData)
}

var apiHandler *APIHandler
var apiClient *TestAPIClient
var employeeRepository data.EmployeeRepository
var employee1, employee2 data.Employee

func setupEmployees() {
	employeeRepository.Init()

	employee1 = data.Employee{
		Name:     "Name 1",
		Age:      21,
		Position: "Position 1",
	}
	employeeRepository.SaveEmployee(employee1)

	employee2 = data.Employee{
		Name:     "Name 2",
		Age:      22,
		Position: "Position 2",
	}
	employeeRepository.SaveEmployee(employee2)
}

func TestMain(m *testing.M) {
	employeeRepository = &data.InMemoryEmployeeRepository{}
	apiHandler = NewAPIHandler(employeeRepository)

	server := httptest.NewUnstartedServer(apiHandler.Handler())
	go func() {
		server.Start()
	}()
	defer server.Close()

	time.Sleep(100 * time.Millisecond)

	apiClient = &TestAPIClient{
		baseUrl: server.URL,
		client:  server.Client(),
	}

	m.Run()
}

func TestGetEmployees(t *testing.T) {
	setupEmployees()

	var respData []employeeListItem
	err := apiClient.Get("/employees", &respData)
	if err != nil {
		t.Error(err)
		return
	}

	expectedLen := 2
	if len(respData) != expectedLen {
		t.Errorf("Expected %d employees, got %d", expectedLen, len(respData))
		return
	}
	if respData[0].ID != "1" ||
		respData[0].Name != "Name 1" ||
		respData[0].Age != 21 ||
		respData[0].Position != "Position 1" {
		t.Errorf("Mismatch for item 1: %v", respData[0])
		return
	}
	if respData[1].ID != "2" ||
		respData[1].Name != "Name 2" ||
		respData[1].Age != 22 ||
		respData[1].Position != "Position 2" {
		t.Errorf("Mismatch for item 2: %v", respData[1])
		return
	}
}

func TestGetEmployee1(t *testing.T) {
	setupEmployees()

	var respData employeeResponse
	err := apiClient.Get("/employees/1", &respData)
	if err != nil {
		t.Error(err)
		return
	}

	expectedRespData := employeeResponse{
		ID:       "1",
		Name:     employee1.Name,
		Age:      employee1.Age,
		Position: employee1.Position,
	}
	if respData.ID != expectedRespData.ID ||
		respData.Name != expectedRespData.Name ||
		respData.Age != expectedRespData.Age ||
		respData.Position != expectedRespData.Position {
		t.Errorf("Expected: %v, got: %v", expectedRespData, respData)
	}
}

func TestGetEmployee2(t *testing.T) {
	setupEmployees()

	var respData employeeResponse
	err := apiClient.Get("/employees/2", &respData)
	if err != nil {
		t.Error(err)
		return
	}

	expectedRespData := employeeResponse{
		ID:       "2",
		Name:     employee2.Name,
		Age:      employee2.Age,
		Position: employee2.Position,
	}
	if respData.ID != expectedRespData.ID ||
		respData.Name != expectedRespData.Name ||
		respData.Age != expectedRespData.Age ||
		respData.Position != expectedRespData.Position {
		t.Errorf("Expected: %v, got: %v", expectedRespData, respData)
	}
}

func TestPostEmployee(t *testing.T) {
	reqData := createEmployeeRequest{
		Name:     "New Name",
		Age:      30,
		Position: "New Position",
	}
	var respData employeeResponse
	err := apiClient.Post("/employees", &reqData, &respData)
	if err != nil {
		t.Error(err)
		return
	}

	expectedRespData := employeeResponse{
		ID:       "3",
		Name:     reqData.Name,
		Age:      reqData.Age,
		Position: reqData.Position,
	}
	if respData.ID != expectedRespData.ID ||
		respData.Name != reqData.Name ||
		respData.Age != reqData.Age ||
		respData.Position != reqData.Position {
		t.Errorf("Expected: %v, got: %v", expectedRespData, respData)
	}
}
