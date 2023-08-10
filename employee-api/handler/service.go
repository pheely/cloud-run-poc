package handler

import (
	"context"
	"io"
	"encoding/base32"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog/log"
	"github.com/pheely/employee-api/stores"
)

type sessionKey string

const sessionIDKey sessionKey = "sessionID"

func JsonHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type Service struct {
	Store        stores.Store
	SessionStore sessions.Store
}

func (s *Service) SessionHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.SessionStore.Get(r, "session")
		if session.Values["ID"] == nil {
			session.Values["ID"] = strings.TrimRight(
				base32.StdEncoding.EncodeToString(
					securecookie.GenerateRandomKey(16)), "=")
		}
		ctx := context.WithValue(r.Context(), sessionIDKey, session.Values["ID"])
		err := session.Save(r, w)
		if err != nil {
			log.Err(err).Msg("Error saving session to response")
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type employee struct {
	Id	   string  `json:"id"`
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Department string `json:"department"`
	Salary    int    `json:"salary"`
	Age      int    `json:"age"`
}

func listID(r *http.Request) string {
	return r.Context().Value(sessionIDKey).(string)
}

func addId(t stores.Employee) employee {
	tU := employee{
		Id:			t.ID,
		First_Name:     	t.First_Name,
		Last_Name: 		t.Last_Name,
		Department:     t.Department,
		Salary:       	t.Salary,
		Age: 			t.Age,
	}
	return tU
}

func (s *Service) Help(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Employee API v4 \n")
}

func (s *Service) List(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	list, err := s.Store.List(listID(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := []employee{}
	for _, t := range list {
		result = append(result, addId(t))
	}
	json.NewEncoder(w).Encode(result)
}

func (s *Service) Clear(w http.ResponseWriter, r *http.Request) {
	s.Store.Clear(listID(r))
	s.List(w, r)
}

func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := s.Store.Delete(listID(r), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Service) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	decoder := json.NewDecoder(r.Body)

	var newT stores.Employee
	err := decoder.Decode(&newT)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := s.Store.Update(listID(r), id, &newT)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if res == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(addId(*res))
}

func (s *Service) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t, err := s.Store.Get(listID(r), vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if t != nil {
		json.NewEncoder(w).Encode(addId(*t))
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t stores.Employee
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = s.Store.Create(listID(r), &t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(addId(t))
}

func (s *Service) Tokenize(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t stores.Employee
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = s.Store.Create(listID(r), &t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(addId(t))
}

func (s *Service) Detokenize(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t stores.Employee
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = s.Store.Create(listID(r), &t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(addId(t))
}
