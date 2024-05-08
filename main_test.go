package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
)

func TestAPI_ListEmployeesHandler(t *testing.T) {
	store := NewStore()
	api := NewAPI(store)

	// Add some employees to the store
	for i := 0; i < 10; i++ {
		emp := Employee{ID: i + 1, Name: "Employee " + strconv.Itoa(i+1), Position: "Position " + strconv.Itoa(i+1), Salary: float64((i + 1) * 10000)}
		store.CreateEmployee(emp)
	}

	tests := []struct {
		name     string
		page     int
		pageSize int
		status   int
	}{
		{"Page 1 with 5 items", 1, 5, http.StatusOK},
		{"Page 2 with 5 items", 2, 5, http.StatusOK},
		{"Invalid page", -1, 5, http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("/employees/%d/%d", tc.page, tc.pageSize), nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			api.ListEmployeesHandler(rr, req)
			if rr.Code != tc.status {
				t.Errorf("Expected status %d, got %d", tc.status, rr.Code)
			}
			if tc.status == http.StatusOK {
				var employees []Employee
				if err := json.Unmarshal(rr.Body.Bytes(), &employees); err != nil {
					t.Errorf("Error decoding JSON response: %v", err)
				}
				if len(employees) != tc.pageSize {
					t.Errorf("Expected %d employees, got %d", tc.pageSize, len(employees))
				}
			}
		})
	}
}

func TestAPI_CreateEmployeeHandler(t *testing.T) {
	store := NewStore()
	api := NewAPI(store)

	emp := Employee{ID: 1, Name: "Testuser", Position: "Developer", Salary: 50000}
	jsonEmp, _ := json.Marshal(emp)

	req, err := http.NewRequest("POST", "/employees", bytes.NewBuffer(jsonEmp))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	api.CreateEmployeeHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}
}

func TestAPI_GetEmployeeHandler(t *testing.T) {
	store := NewStore()
	api := NewAPI(store)

	// Add an employee to the store
	emp := Employee{ID: 1, Name: "test user", Position: "Developer", Salary: 50000}
	store.CreateEmployee(emp)

	tests := []struct {
		name   string
		id     int
		status int
	}{
		{"Valid employee ID", 1, http.StatusOK},
		{"Invalid employee ID", 2, http.StatusNotFound},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/employees/"+strconv.Itoa(tc.id), nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/employees/"+strconv.Itoa(tc.id), api.GetEmployeeHandler).Methods("GET")

			// Serve the HTTP request to the ResponseRecorder
			router.ServeHTTP(rr, req)
			if rr.Code != tc.status {
				t.Errorf("Expected status %d, got %d", tc.status, rr.Code)
			}
		})
	}
}

func TestAPI_UpdateEmployeeHandler(t *testing.T) {
	store := NewStore()
	api := NewAPI(store)

	// Add an employee to the store
	emp := Employee{ID: 1, Name: "TestUser", Position: "Developer", Salary: 50000}
	store.CreateEmployee(emp)

	newEmp := Employee{ID: 1, Name: "TestUser", Position: "Manager", Salary: 60000}
	jsonNewEmp, _ := json.Marshal(newEmp)

	req, err := http.NewRequest("PUT", "/employees/1", bytes.NewBuffer(jsonNewEmp))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	api.UpdateEmployeeHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestAPI_DeleteEmployeeHandler(t *testing.T) {
	store := NewStore()
	api := NewAPI(store)

	// Add an employee to the store
	emp := Employee{ID: 1, Name: "TestUser", Position: "Developer", Salary: 50000}
	store.CreateEmployee(emp)

	req, err := http.NewRequest("DELETE", "/employees/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	api.DeleteEmployeeHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}
