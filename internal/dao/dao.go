package dao

import (
	"database/sql"
	"log"
)

type DAO struct {
	db *sql.DB
}

type Employee struct {
	ID       int
	Name     string
	Age      int
	Position string
}

const initQueries = `
CREATE TABLE IF NOT EXISTS employees (
	id SERIAL,
	name VARCHAR UNIQUE,
	age SMALLINT,
	position VARCHAR
);
INSERT INTO employees (name, age, position)
VALUES
	('Иван', 40, 'старший инженер-программист'),
	('Петр', 30, 'инженер-программист'),
	('Анастасия', 25, 'дизайнер'),
	('Николай', 35, 'менеджер'),
	('Семён', 45, 'директор')
ON CONFLICT DO NOTHING;
`

func (dao *DAO) Init() {
	log.Println("Initializing DB")

	db, err := sql.Open("postgres", "host=postgres user=crud_user password=crud_password dbname=crud_example sslmode=disable connect_timeout=5")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(initQueries)
	if err != nil {
		log.Fatal(err)
	}

	dao.db = db
}

func (dao *DAO) Deinit() {
	dao.db.Close()
}

func (dao *DAO) EmployeeExists(name string) (bool, error) {
	query := "SELECT EXISTS (SELECT * FROM employees WHERE name = $1)"
	exists := false

	err := dao.db.QueryRow(query, name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (dao *DAO) GetEmployees() ([]Employee, error) {
	rows, err := dao.db.Query("SELECT * FROM employees")
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

func (dao *DAO) GetEmployeeByID(id int) (*Employee, error) {
	query := "SELECT * FROM employees WHERE id = $1"
	employee := Employee{}

	err := dao.db.QueryRow(query, id).Scan(&employee.ID, &employee.Name, &employee.Age, &employee.Position)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &employee, nil
}

func (dao *DAO) SaveEmployee(employee Employee) (int, error) {
	query := "INSERT INTO employees (name, age, position) VALUES ($1, $2, $3) RETURNING id"
	id := 0

	err := dao.db.QueryRow(query, employee.Name, employee.Age, employee.Position).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
