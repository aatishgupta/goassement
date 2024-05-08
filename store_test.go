// store_test.go
package main

import (
	"reflect"
	"testing"
)

func TestStore_CreateEmployee(t *testing.T) {
	store := NewStore()
	emp := Employee{ID: 1, Name: "John Doe", Position: "Developer", Salary: 50000}
	store.CreateEmployee(emp)
	if len(store.employees) != 1 {
		t.Errorf("Expected 1 employee, got %d", len(store.employees))
	}
}

func TestStore_GetEmployeeByID(t *testing.T) {
	store := NewStore()
	emp := Employee{ID: 1, Name: "John Doe", Position: "Developer", Salary: 50000}
	store.employees = append(store.employees, emp)

	tests := []struct {
		id         int
		wantEmp    Employee
		wantExists bool
	}{
		{1, emp, true},
		{2, Employee{}, false},
	}

	for _, tc := range tests {
		gotEmp, exists := store.GetEmployeeByID(tc.id)
		if exists != tc.wantExists {
			t.Errorf("For ID %d, expected exists=%t, got exists=%t", tc.id, tc.wantExists, exists)
		}
		if !reflect.DeepEqual(gotEmp, tc.wantEmp) {
			t.Errorf("For ID %d, expected employee %v, got %v", tc.id, tc.wantEmp, gotEmp)
		}
	}
}

func TestStore_UpdateEmployee(t *testing.T) {
	store := NewStore()
	emp := Employee{ID: 1, Name: "John Doe", Position: "Developer", Salary: 50000}
	store.employees = append(store.employees, emp)

	newEmp := Employee{ID: 1, Name: "Jane Doe", Position: "Manager", Salary: 60000}
	store.UpdateEmployee(newEmp)

	gotEmp, _ := store.GetEmployeeByID(1)
	if !reflect.DeepEqual(gotEmp, newEmp) {
		t.Errorf("Expected updated employee %v, got %v", newEmp, gotEmp)
	}
}

func TestStore_DeleteEmployee(t *testing.T) {
	store := NewStore()
	emp := Employee{ID: 1, Name: "John Doe", Position: "Developer", Salary: 50000}
	store.employees = append(store.employees, emp)

	store.DeleteEmployee(1)
	if len(store.employees) != 0 {
		t.Errorf("Expected 0 employees after deletion, got %d", len(store.employees))
	}
}
