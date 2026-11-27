package data

import (
	"database/sql"
	"log"
	"slices"
)

type EmployeeRepository interface {
	Init()
	EmployeeExists(name string) (bool, error)
	GetEmployees() ([]Employee, error)
	GetEmployeeByID(id int) (*Employee, error)
	SaveEmployee(employee Employee) (int, error)
}

type InMemoryEmployeeRepository struct {
	data      []Employee
	currentId int
}

func (r *InMemoryEmployeeRepository) Init() {
	r.data = make([]Employee, 0)
	r.currentId = 0
}

func (r *InMemoryEmployeeRepository) EmployeeExists(name string) (bool, error) {
	return slices.ContainsFunc(r.data, func(e Employee) bool {
		return e.Name == name
	}), nil
}

func (r *InMemoryEmployeeRepository) GetEmployees() ([]Employee, error) {
	return r.data, nil
}

func (r *InMemoryEmployeeRepository) GetEmployeeByID(id int) (*Employee, error) {
	i := slices.IndexFunc(r.data, func(e Employee) bool {
		return e.ID == id
	})
	if i >= 0 {
		return &r.data[i], nil
	} else {
		return nil, nil
	}
}

func (r *InMemoryEmployeeRepository) SaveEmployee(employee Employee) (int, error) {
	if employee.ID == 0 {
		r.data = append(r.data, employee)
		r.currentId++
		r.data[len(r.data)-1].ID = r.currentId
		return r.currentId, nil
	} else {
		i := slices.IndexFunc(r.data, func(e Employee) bool {
			return e.ID == employee.ID
		})
		r.data[i].Age = employee.Age
		r.data[i].Name = employee.Name
		r.data[i].Position = employee.Position
		return employee.ID, nil
	}
}

type SQLEmployeeRepository struct {
	db *sql.DB
}

func NewSQLEmployeeRepository(db *sql.DB) *SQLEmployeeRepository {
	return &SQLEmployeeRepository{
		db: db,
	}
}

func (r *SQLEmployeeRepository) Init() {
	log.Println("Initializing Employee DB")

	r.createTables()
	r.createTestData()
}

func (r *SQLEmployeeRepository) createTables() {
	_, err := r.db.Exec(`
CREATE TABLE IF NOT EXISTS employees (
	id SERIAL,
	name VARCHAR UNIQUE,
	age SMALLINT,
	position VARCHAR
);`)
	if err != nil {
		log.Fatal(err)
	}
}

func (r *SQLEmployeeRepository) createTestData() {
	_, err := r.db.Exec(`
INSERT INTO employees (name, age, position)
VALUES
	('Иван', 40, 'старший инженер-программист'),
	('Петр', 30, 'инженер-программист'),
	('Анастасия', 25, 'дизайнер'),
	('Николай', 35, 'менеджер'),
	('Семён', 45, 'директор')
ON CONFLICT DO NOTHING;`)
	if err != nil {
		log.Fatal(err)
	}
}

func (r *SQLEmployeeRepository) EmployeeExists(name string) (bool, error) {
	query := "SELECT EXISTS (SELECT * FROM employees WHERE name = $1)"
	exists := false

	err := r.db.QueryRow(query, name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *SQLEmployeeRepository) GetEmployees() ([]Employee, error) {
	rows, err := r.db.Query("SELECT * FROM employees")
	if err != nil {
		return nil, err
	}

	employees := []Employee{}
	for rows.Next() {
		employee := Employee{}
		err = rows.Scan(&employee.ID, &employee.Name, &employee.Age, &employee.Position)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

func (r *SQLEmployeeRepository) GetEmployeeByID(id int) (*Employee, error) {
	query := "SELECT * FROM employees WHERE id = $1"
	employee := Employee{}

	err := r.db.QueryRow(query, id).Scan(&employee.ID, &employee.Name, &employee.Age, &employee.Position)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &employee, nil
}

func (r *SQLEmployeeRepository) SaveEmployee(employee Employee) (int, error) {
	query := "INSERT INTO employees (name, age, position) VALUES ($1, $2, $3) RETURNING id"
	id := 0

	err := r.db.QueryRow(query, employee.Name, employee.Age, employee.Position).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
