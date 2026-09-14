package data

import (
	"database/sql"
	"log"
)

type EmployeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{
		db: db,
	}
}

func (r *EmployeeRepository) Init() {
	log.Println("Initializing Employee DB")

	r.createTables()
	r.createTestData()
}

func (r *EmployeeRepository) createTables() {
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

func (r *EmployeeRepository) createTestData() {
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

func (r *EmployeeRepository) EmployeeExists(name string) (bool, error) {
	query := "SELECT EXISTS (SELECT * FROM employees WHERE name = $1)"
	exists := false

	err := r.db.QueryRow(query, name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *EmployeeRepository) GetEmployees() ([]Employee, error) {
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

func (r *EmployeeRepository) GetEmployeeByID(id int) (*Employee, error) {
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

func (r *EmployeeRepository) SaveEmployee(employee Employee) (int, error) {
	query := "INSERT INTO employees (name, age, position) VALUES ($1, $2, $3) RETURNING id"
	id := 0

	err := r.db.QueryRow(query, employee.Name, employee.Age, employee.Position).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
