package main

import (
	"sync"
)

// Store is a simple in-memory data store for managing employee records
type Store struct {
	mu        sync.RWMutex
	employees []Employee
}

// NewStore creates a new instance of the Store
func NewStore() *Store {
	return &Store{
		employees: make([]Employee, 0),
	}
}

// CreateEmployee adds a new employee to the data store
func (s *Store) CreateEmployee(emp Employee) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.employees = append(s.employees, emp)
}

// GetEmployeeByID retrieves an employee from the data store by ID
func (s *Store) GetEmployeeByID(id int) (Employee, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, emp := range s.employees {
		if emp.ID == id {
			return emp, true
		}
	}
	return Employee{}, false
}

// UpdateEmployee updates the details of an existing employee in the data store
func (s *Store) UpdateEmployee(emp Employee) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, e := range s.employees {
		if e.ID == emp.ID {
			s.employees[i] = emp
			break
		}
	}
}

// DeleteEmployee deletes an employee from the data store by ID
func (s *Store) DeleteEmployee(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, emp := range s.employees {
		if emp.ID == id {
			s.employees = append(s.employees[:i], s.employees[i+1:]...)
			break
		}
	}
}

// GetPaginatedEmployees retrieves a slice of employee records with pagination
func (s *Store) GetPaginatedEmployees(offset, limit int) []Employee {
	s.mu.RLock()
	defer s.mu.RUnlock()

	start := offset
	end := offset + limit
	if end > len(s.employees) {
		end = len(s.employees)
	}

	return s.employees[start:end]
}
