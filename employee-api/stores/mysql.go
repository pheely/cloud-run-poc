package stores

import (
	"database/sql"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type SQLStore struct {
	DBUrl string
	DB    *sql.DB
}

func NewSQLStore(dbURL string) *SQLStore {
	s := SQLStore{dbURL, nil}
	return &s
}

func (s *SQLStore) Connect() error {
	db, err := sql.Open("mysql", s.DBUrl)

	if err != nil {
		return err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(10)

	s.DB = db

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (s SQLStore) Create(_ string, t *Employee) error {

	stmt, err := s.DB.Prepare(`INSERT INTO employees(first_name, last_name, department, salary, age) VALUES (
		?,
		?, 
		?,
		?,
		?
	)`)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(t.First_Name, t.Last_Name, t.Department, t.Salary, t.Age)
	if err != nil {
		return err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = strconv.FormatInt(lastID, 10)
	return nil
}

func (s SQLStore) Delete(_ string, id string) error {
	stmt, err := s.DB.Prepare("DELETE FROM employees WHERE id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (s SQLStore) Update(listID string, id string, newT *Employee) (*Employee, error) {

	t, err := s.Get(listID, id)
	if err != nil {
		return nil, err
	}
	if t != nil {
		if newT.First_Name != "" {
			t.First_Name = newT.First_Name
		}
		if newT.Last_Name != "" {
			t.Last_Name = newT.Last_Name
		}
		if newT.Department != "" {
			t.Department = newT.Department
		}
		if newT.Salary != 0 {
			t.Salary = newT.Salary
		}
		if newT.Age != 0 {
			t.Age = newT.Age
		}

		stmt, err := s.DB.Prepare(`UPDATE employees SET 
			first_name = ?, 
			last_name = ?, 
			department = ?, 
			salary = ?, 
			age = ?
			WHERE id=?`)

		if err != nil {
			return nil, err
		}
		_, err = stmt.Exec(t.First_Name, t.Last_Name, t.Department, t.Salary, t.Age, id)
		if err != nil {
			return nil, err
		}

		return t, nil
	}
	return nil, nil
}

func (s SQLStore) Get(_ string, id string) (*Employee, error) {
	rows, err := s.DB.Query("select id, first_name, last_name, department, salary, age from employees where id=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	t := Employee{}
	for rows.Next() {
		err := rows.Scan(&t.ID, &t.First_Name, &t.Last_Name, &t.Department, &t.Salary, &t.Age)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s SQLStore) Clear(_ string) error {
	stmt, err := s.DB.Prepare("DELETE FROM employees")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (s SQLStore) List(_ string) ([]Employee, error) {
	rows, err := s.DB.Query("select id, first_name, last_name, department, salary, age from employees")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []Employee{}
	for rows.Next() {
		t := Employee{}
		err := rows.Scan(&t.ID, &t.First_Name, &t.Last_Name, &t.Department, &t.Salary, &t.Age)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}
