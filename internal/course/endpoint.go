package course

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nOm-1986/goApi_web/pkg/meta"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		Create Controller
		GetAll Controller
		Get    Controller
		Update Controller
		Delete Controller
	}

	CreateReq struct {
		Name      string `json:"course_name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	UpdateReq struct {
		Name      *string `json:"course_name"`
		StartDate *string `json:"start_date"`
		EndDate   *string `json:"end_date"`
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
		GetAll: makeGetAllEndpoint(s),
		Get:    makeGetEndpoint(s),
		Update: makeUpdateEndpoint(s),
		//Delete: makeDeletEndpoint(s),
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
		if req.Name == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "First name is required!!"})
			return
		}
		if req.StartDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Last name is required!!"})
			return
		}
		if req.EndDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Email is required!!"})
			return
		}

		course, err := s.Create(req.Name, req.StartDate, req.EndDate)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: course})
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		urlParams := r.URL.Query()
		filters := Filters{
			Name: urlParams.Get("name_course"),
			/*StartDate: urlParams.Get("start_date"),
			EndDate:   urlParams.Get("end_date"),*/
		}

		limit, _ := strconv.Atoi(urlParams.Get("limit"))
		page, _ := strconv.Atoi(urlParams.Get("page"))

		count, err := s.Count(filters)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		meta, err := meta.NewMeta(page, limit, count)

		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		courses, err := s.GetAll(filters, meta.Offset(), meta.Limit())
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: courses, Meta: meta})
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		course, err := s.Get(id)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Course does not exist"})
		}
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: course})
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateReq

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Formato de solicitud inválido!"})
			return
		}
		//Validaciones
		if req.Name != nil && *req.Name == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "El nombre del curso es requerido!"})
			return
		}
		if req.StartDate != nil && *req.StartDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "La fecha de incio del curso es requerida!"})
			return
		}
		//Ak puedo validar que el formato sea de tipo fecha.
		if req.EndDate != nil && *req.EndDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "La fecha de finalización del curso es requerida!"})
			return
		}

		path := mux.Vars(r)
		id := path["id"]

		if err := s.Update(id, req.Name, req.StartDate, req.EndDate); err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Course does not exits!"})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: "Updated successfully"})
	}
}
