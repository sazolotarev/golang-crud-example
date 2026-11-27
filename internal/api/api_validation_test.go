package api

import (
	"strings"
	"testing"
)

func newCreateEmployeeRequest() createEmployeeRequest {
	return createEmployeeRequest{
		Age:      18,
		Name:     "Name",
		Position: "Position",
	}
}

func TestValidate_createEmployeeRequest_ageBelow18(t *testing.T) {
	req := newCreateEmployeeRequest()
	req.Age = 17
	err := req.validate()
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "age") {
		t.Errorf("age %d must not pass validation", req.Age)
	}
}

func TestValidate_createEmployeeRequest_ageAbove100(t *testing.T) {
	req := newCreateEmployeeRequest()
	req.Age = 101
	err := req.validate()
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "age") {
		t.Errorf("age %d must not pass validation", req.Age)
	}
}

func TestValidate_createEmployeeRequest_nameIsEmpty(t *testing.T) {
	req := newCreateEmployeeRequest()
	req.Name = ""
	err := req.validate()
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "name is required") {
		t.Error("empty name must not pass validation")
	}
}

func TestValidate_createEmployeeRequest_nameIsTooShort(t *testing.T) {
	req := newCreateEmployeeRequest()
	req.Name = "a"
	err := req.validate()
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "name must be at least") {
		t.Error("one character name must not pass validation")
	}
}

func TestValidate_createEmployeeRequest_positionIsEmpty(t *testing.T) {
	req := newCreateEmployeeRequest()
	req.Position = ""
	err := req.validate()
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "position is required") {
		t.Error("empty position must not pass validation")
	}
}
