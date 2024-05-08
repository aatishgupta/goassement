package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// API struct defines the RESTful API for the employee data store
type API struct {
	store *Store
}

// NewAPI creates a new instance of the API
func NewAPI(store *Store) *API {
	return &API{store: store}
}

// ListEmployeesHandler handles the endpoint for listing employees with pagination
func (a *API) ListEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page, _ := strconv.Atoi(vars["page"])
	pageSize, _ := strconv.Atoi(vars["pageSize"])

	// Calculate offset based on page and pageSize
	offset := (page - 1) * pageSize

	// Retrieve paginated employee records from the store
	employees := a.store.GetPaginatedEmployees(offset, pageSize)

	// Return paginated employee records as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}

// CreateEmployeeHandler handles the endpoint for creating a new employee
func (a *API) CreateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, "Invalid request body"+err.Error(), http.StatusBadRequest)
		return
	}
	a.store.CreateEmployee(emp)
	w.WriteHeader(http.StatusCreated)
}

// GetEmployeeHandler handles the endpoint for retrieving an employee by ID
func (a *API) GetEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}
	emp, exists := a.store.GetEmployeeByID(id)
	if !exists {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emp)
}

// UpdateEmployeeHandler handles the endpoint for updating an employee
func (a *API) UpdateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}
	var emp Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	emp.ID = id
	a.store.UpdateEmployee(emp)
	w.WriteHeader(http.StatusOK)
}

// DeleteEmployeeHandler handles the endpoint for deleting an employee by ID
func (a *API) DeleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}
	a.store.DeleteEmployee(id)
	w.WriteHeader(http.StatusOK)
}
func (a *API) InsertDummyEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	for i := 1; i <= 1000; i++ {
		emp := Employee{
			ID:       i,
			Name:     fmt.Sprintf("Employee%d", i),
			Position: "Developer",
			Salary:   50000,
		}

		a.store.CreateEmployee(emp)
	}

	w.WriteHeader(http.StatusCreated)
}

// SetupRouter sets up the routes for the RESTful API
func (a *API) SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/employees/{page}/{pageSize}", a.ListEmployeesHandler).Methods("GET")
	r.HandleFunc("/employees", a.CreateEmployeeHandler).Methods("POST")
	r.HandleFunc("/employees/{id}", a.GetEmployeeHandler).Methods("GET")
	r.HandleFunc("/employees/{id}", a.UpdateEmployeeHandler).Methods("PUT")
	r.HandleFunc("/employees/{id}", a.DeleteEmployeeHandler).Methods("DELETE")
	r.HandleFunc("/insert-dummy-employees", a.InsertDummyEmployeesHandler).Methods("POST")
	return r
}

func main() {
	store := NewStore()
	api := NewAPI(store)

	r := api.SetupRouter()
	http.ListenAndServe(":8080", r)
}
