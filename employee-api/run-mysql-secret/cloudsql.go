package main

import (
        "context"
        "database/sql"
        "fmt"
        "log"
        "net"
        "strconv"
        "time"

        "cloud.google.com/go/cloudsqlconn"
        "github.com/go-sql-driver/mysql"
)

type DataSource struct {
        Parameters ConnectionParams
        DB *sql.DB
}

func createDataSource(params ConnectionParams) DataSource {
        ds := DataSource{params, nil}
        ds.Connect()

        return ds
}

func (ds *DataSource) Connect() error {
        d, err := cloudsqlconn.NewDialer(context.Background())
        if err != nil {
                return fmt.Errorf("cloudsqlconn.NewDialer: %w", err)
        }
        var opts []cloudsqlconn.DialOption
        if ds.Parameters.PrivateIP != "" {
             opts = append(opts, cloudsqlconn.WithPrivateIP())
        }
        mysql.RegisterDialContext("cloudsqlconn",
                func(ctx context.Context, addr string) (net.Conn, error) {
                        return d.Dial(ctx, ds.Parameters.InstanceConnectionName, opts...)
                })

        dbURI := fmt.Sprintf("%s:%s@cloudsqlconn(localhost:3306)/%s?parseTime=true",
                ds.Parameters.DBUser, ds.Parameters.DBPwd, ds.Parameters.DBName)

        log.Printf("dbURI: %s", dbURI)

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

        log.Println("sql.DB created")

        ds.DB = dbPool
        return nil
}

func (ds DataSource) Create(em *Employee) error {
        stmt, err := ds.DB.Prepare(`INSERT INTO employees(first_name, last_name, department, salary, age) VALUES (
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

func (ds DataSource) Delete(id string) error {
        stmt, err := ds.DB.Prepare("DELETE FROM employees WHERE id=?")
        if err != nil {
                return err
        }
        _, err = stmt.Exec(id)
        if err != nil {
                return err
        }
        return nil
}

func (ds DataSource) Update(id string, newT *Employee) (*Employee, error) {

        t, err := ds.Get(id)
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

                stmt, err := ds.DB.Prepare(`UPDATE employees SET 
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

func (ds DataSource) Get(id string) (*Employee, error) {
        rows, err := ds.DB.Query("select id, first_name, last_name, department, salary, age from employees where id=?", id)
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

func (ds DataSource) Clear() error {
        stmt, err := ds.DB.Prepare("DELETE FROM employees")
        if err != nil {
                return err
        }
        _, err = stmt.Exec()
        if err != nil {
                return err
        }
        return nil
}

func (ds DataSource) List() ([]Employee, error) {
        rows, err := ds.DB.Query("select id, first_name, last_name, department, salary, age from employees")
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