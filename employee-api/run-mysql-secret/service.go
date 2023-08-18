package main

import (
	"io"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func jsonHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type Service struct {
	connection DataSource
}

type employee struct {
	Id	   string  `json:"id"`
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Department string `json:"department"`
	Salary    int    `json:"salary"`
	Age      int    `json:"age"`
}

func (s *Service) Help(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Employee API v4 \n")
}

func (s *Service) List(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	list, err := s.connection.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := []employee{}
	for _, t := range list {
		result = append(result, convertToEmployeeJson(t))
	}
	json.NewEncoder(w).Encode(result)
}

func (s *Service) Clear(w http.ResponseWriter, r *http.Request) {
	s.connection.Clear()
	s.List(w, r)
}

func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := s.connection.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Service) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	decoder := json.NewDecoder(r.Body)

	var newT Employee
	err := decoder.Decode(&newT)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := s.connection.Update(id, &newT)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if res == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(convertToEmployeeJson(*res))
}

func (s *Service) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t, err := s.connection.Get(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if t != nil {
		json.NewEncoder(w).Encode(convertToEmployeeJson(*t))
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t Employee
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = s.connection.Create(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(convertToEmployeeJson(t))
}

func convertToEmployeeJson(t Employee) employee {
	jsonT := employee {
		Id: t.ID,
		First_Name: t.First_Name,
		Last_Name: t.Last_Name,
		Department: t.Department,
		Salary: t.Salary,
		Age: t.Age,
	}

	return jsonT
}