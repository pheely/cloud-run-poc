package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/go-sql-driver/mysql"
)

type ConnectionPool struct {
	Parameters ConnectionInfo
	DB *sql.DB
}

type ConnectionInfo struct {
	DBUser                 string
	DBPwd                  string
	DBName                 string
	// InstanceConnectionName string
	PrivateIP              string
}

func CreateConnectionPool(params ConnectionInfo) *ConnectionPool {
	p := ConnectionPool{params, nil}
	return &p
}

func (pool *ConnectionPool) Connect() error {
	// d, err := cloudsqlconn.NewDialer(context.Background())
	// if err != nil {
	// 	return nil, fmt.Errorf("cloudsqlconn.NewDialer: %w", err)
	// }
	// var opts []cloudsqlconn.DialOption
	// if pool.Parameters.PrivateIP {
	// 	opts = append(opts, cloudsqlconn.WithPrivateIP())
	// }
	// mysql.RegisterDialContext("cloudsqlconn",
	// 	func(ctx context.Context, addr string) (net.Conn, error) {
	// 		return d.Dial(ctx, pool.Parameters.InstanceConnectionName, opts...)
	// 	})

	dbURI := fmt.Sprintf("%s:%s@cloudsqlconn(localhost:3306)/%s?parseTime=true",
		pool.Parameters.DBUser, pool.Parameters.DBPwd, pool.Parameters.DBName)

	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}

	dbPool.SetConnMaxLifetime(time.Minute * 3)
	dbPool.SetMaxOpenConns(2)
	dbPool.SetMaxIdleConns(10)

	err = dbPool.Ping()
	if err != nil {
		return fmt.Errorf("sql.Ping: %w", err)
	}

	pool.DB = dbPool
	return nil
}

func (pool ConnectionPool) Create(em *Employee) error {
	stmt, err := pool.DB.Prepare(`INSERT INTO employees(first_name, last_name, department, salary, age) VALUES (
		?,
		?, 
		?,
		?,
		?
	)`)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(em.First_Name, em.Last_Name, em.Department, em.Salary, em.Age)
	if err != nil {
		return err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	em.ID = strconv.FormatInt(lastID, 10)
	return nil
}

func (pool ConnectionPool) Delete(id string) error {
	stmt, err := pool.DB.Prepare("DELETE FROM employees WHERE id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (pool ConnectionPool) Update(id string, newT *Employee) (*Employee, error) {

	t, err := s.Get(id)
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

		stmt, err := pool.DB.Prepare(`UPDATE employees SET 
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

func (pool ConnectionPool) Get(id string) (*Employee, error) {
	rows, err := pool.DB.Query("select id, first_name, last_name, department, salary, age from employees where id=?", id)
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

func (pool ConnectionPool) Clear() error {
	stmt, err := pool.DB.Prepare("DELETE FROM employees")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (pool ConnectionPool) List() ([]Employee, error) {
	rows, err := pool.DB.Query("select id, first_name, last_name, department, salary, age from employees")
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