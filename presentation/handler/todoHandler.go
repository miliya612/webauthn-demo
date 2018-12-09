package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/miliya612/webauthn-demo/domain/errUtil"
	"github.com/miliya612/webauthn-demo/domain/model"
	"github.com/miliya612/webauthn-demo/domain/service"
	"github.com/miliya612/webauthn-demo/presentation/httputil"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

//var KB int64 = 1024

type TodoHandler interface {
	TodoIndex(r *http.Request)
	TodoShow(r *http.Request)
	TodoCreate(r *http.Request)
	TodoDelete(r *http.Request)
}

type todoHandler struct {
	service service.TodoService
}

//func NewTodoHandler(s service.TodoService) TodoHandler {
//	return &todoHandler{service: s}
//}

func (h *todoHandler) TodoIndex(w http.ResponseWriter, _ *http.Request) {
	todos, err := h.service.All()
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "something went wrong", err)
	}
	httputil.Ok(w, todos)
}

func (h *todoHandler) TodoShow(w http.ResponseWriter, r *http.Request) {
	id, err := parseTodoId(r)
	if err != nil {
		httputil.Error(w, http.StatusUnprocessableEntity, "invalid parameter", err)
	}

	t, err := h.service.Find(id)
	if err != nil {
		switch err.(type) {
		case errUtil.ErrTodoNotFound:
			httputil.Error(w, http.StatusNotFound, fmt.Sprintf("failed to search todo with id: %d", id), err)
		default:
			httputil.Error(w, http.StatusInternalServerError, "something went wrong", err)
		}
	}

	httputil.Ok(w, t)
}

func (h *todoHandler) TodoCreate(w http.ResponseWriter, r *http.Request) {
	var todo model.Todo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "request body is too large", err)
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &todo); err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed marshalling json", err)
	}

	id, err := h.service.Create(todo)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "something went wrong", err)
	}
	todo.ID = id
	_ = fmt.Sprintf("http://%s/%d", r.Host, id)
	httputil.Created(w, todo)
}

func (h *todoHandler) TodoDelete(w http.ResponseWriter, r *http.Request) {
	id, err := parseTodoId(r)
	if err != nil {
		httputil.Error(w, http.StatusUnprocessableEntity, "invalid parameter", err)
	}

	if err = h.service.Remove(id); err != nil {
		httputil.Error(w, http.StatusNotFound, fmt.Sprintf("failed to delete todo for id: %d", id), err)
	}

	httputil.Empty(w, http.StatusNoContent)
}

func parseTodoId(r *http.Request) (int, error) {
	id, err := strconv.Atoi(mux.Vars(r)["todoId"])
	if err != nil {
		return -1, errors.New("todoId should be number.")
	}
	return id, nil
}
