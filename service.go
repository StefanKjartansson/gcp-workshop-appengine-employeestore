package employees

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

type EmployeeStore interface {
	getContextFromRequest(req *http.Request) context.Context
	List(context.Context) ([]*Employee, error)
	Get(context.Context, string) (*Employee, error)
	Put(context.Context, Employee) (string, error)
}

type Employee struct {
	Id       string
	Name     string
	Role     string
	HireDate time.Time
	Account  string
}

type EmployeeService struct {
	store EmployeeStore
}

func (es EmployeeService) list(w http.ResponseWriter, req *http.Request) {
	ctx := es.store.getContextFromRequest(req)
	l, err := es.store.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(l)
}

func (es EmployeeService) put(w http.ResponseWriter, req *http.Request) {
	var e Employee
	err := json.NewDecoder(req.Body).Decode(&e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx := es.store.getContextFromRequest(req)
	id, err := es.store.Put(ctx, e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	e.Id = id
	encoder := json.NewEncoder(w)
	encoder.Encode(e)
}

func (es EmployeeService) get(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	ctx := es.store.getContextFromRequest(req)
	e, err := es.store.Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(e)
}

func GetHandler(store EmployeeStore) http.Handler {
	service := EmployeeService{
		store: store,
	}
	r := mux.NewRouter()
	r.HandleFunc("/employees/", service.list).
		Methods("GET")
	r.HandleFunc("/employees/", service.put).
		Methods("POST")
	r.HandleFunc("/employees/{id}/", service.get).
		Methods("GET")
	return r
}
