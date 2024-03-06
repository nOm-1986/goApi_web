package user

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nOm-1986/goApi_web/pkg/meta"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	UpdateReq struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"` //omit empty, Si viene vacío lo omita, interface es para poner lo que sea.
		Err    string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}
)

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeletEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Invalid request format"})
			return
		}
		if req.FirstName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "First name is required!!"})
			return
		}
		if req.LastName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Last name is required!!"})
			return
		}
		if req.Email == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Email is required!!"})
			return
		}

		user, err := s.Create(req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: user})
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		//Mapa con los parametros de la ruta.
		path := mux.Vars(r)
		id := path["id"]
		user, err := s.Get(id)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "User does not exist"})
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: user})
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		//URL.Query devuelve el objeto de los parametros que vienen por URL
		urlParams := r.URL.Query()

		filters := Filters{
			FirstName: urlParams.Get("first_name"),
			LastName:  urlParams.Get("last_name"),
		}

		count, err := s.Count(filters)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		meta, err := meta.NewMeta(count)

		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		users, err := s.GetAll(filters)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: users, Meta: meta})
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateReq

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Invalid request format"})
			return
		}
		//Validations
		if req.FirstName != nil && *req.FirstName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "First name is required"})
			return
		}
		if req.LastName != nil && *req.LastName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Last name is required"})
			return
		}
		if req.Email != nil && *req.Email == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Email is required"})
			return
		}

		path := mux.Vars(r)
		id := path["id"]

		if err := s.Update(id, req.FirstName, req.LastName, req.Email, req.Phone); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "User does not exits"})
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"data": "Updated successfully"})
	}
}

func makeDeletEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		if err := s.Delete(id); err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "User does not exist"})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"data": "User deleted succesfully"})
	}
}
