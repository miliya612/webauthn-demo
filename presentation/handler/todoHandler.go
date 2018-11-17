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

var KB int64 = 1024

type TodoHandler interface {
	TodoIndex(r *http.Request) httputil.Responder
	TodoShow(r *http.Request) httputil.Responder
	TodoCreate(r *http.Request) httputil.Responder
	TodoDelete(r *http.Request) httputil.Responder
}

type todoHandler struct {
	service service.TodoService
}

func NewTodoHandler(s service.TodoService) TodoHandler {
	return &todoHandler{service: s}
}

func (h *todoHandler) TodoIndex(_ *http.Request) httputil.Responder {
	todos, err := h.service.All()
	if err != nil {
		return httputil.Error(http.StatusInternalServerError, "something went wrong", err)
	}
	return httputil.Ok(todos)
}

func (h *todoHandler) TodoShow(r *http.Request) httputil.Responder {
	id, err := parseTodoId(r)
	if err != nil {
		return httputil.Error(http.StatusUnprocessableEntity, "invalid parameter", err)
	}

	t, err := h.service.Find(id)
	if err != nil {
		switch err.(type) {
		case errUtil.ErrTodoNotFound:
			return httputil.Error(http.StatusNotFound, fmt.Sprintf("failed to search todo with id: %d", id), err)
		default:
			return httputil.Error(http.StatusInternalServerError, "something went wrong", err)
		}
	}

	return httputil.Ok(t)
}

func (h *todoHandler) TodoCreate(r *http.Request) httputil.Responder {
	var todo model.Todo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*KB)) // 1MB
	if err != nil {
		return httputil.Error(http.StatusInternalServerError, "request body is too large", err)
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &todo); err != nil {
		return httputil.Error(http.StatusInternalServerError, "failed marshalling json", err)
	}

	id, err := h.service.Create(todo)
	if err != nil {
		return httputil.Error(http.StatusInternalServerError, "something went wrong", err)
	}
	todo.ID = id
	location := fmt.Sprintf("http://%s/%d", r.Host, id)
	return httputil.Created(todo, location)
}

func (h *todoHandler) TodoDelete(r *http.Request) httputil.Responder {
	id, err := parseTodoId(r)
	if err != nil {
		return httputil.Error(http.StatusUnprocessableEntity, "invalid parameter", err)
	}

	if err = h.service.Remove(id); err != nil {
		return httputil.Error(http.StatusNotFound, fmt.Sprintf("failed to delete todo for id: %d", id), err)
	}

	return httputil.Empty(http.StatusNoContent)
}

func parseTodoId(r *http.Request) (int, error) {
	id, err := strconv.Atoi(mux.Vars(r)["todoId"])
	if err != nil {
		return -1, errors.New("todoId should be number.")
	}
	return id, nil
}
