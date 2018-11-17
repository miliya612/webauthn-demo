package service

import (
	"github.com/miliya612/webauthn-demo/domain/errUtil"
	"github.com/miliya612/webauthn-demo/domain/model"
	"github.com/miliya612/webauthn-demo/domain/repo"
)

type TodoService interface {
	All() (model.Todos, error)
	Find(int) (model.Todo, error)
	Create(todo model.Todo) (int, error)
	Remove(id int) (error)
}

type todoService struct {
	repo repo.TodoRepo
}

func NewTodoService(repository repo.TodoRepo) TodoService {
	return &todoService{repo: repository}
}

func (s todoService) All() (model.Todos, error){
	return s.repo.GetAll()
}

func (s todoService) Find(id int) (t model.Todo, err error){
	t, err = s.repo.GetByID(id)
	if err != nil {
		err = errUtil.ErrTodoNotFound{}
		return
	}
	return
}

func (s todoService) Create(todo model.Todo) (id int, err error){
	id, err = s.repo.Create(todo)
	return
}

func (s todoService) Remove(id int) (err error){
	count, err := s.repo.Remove(id)
	if err != nil {
		return
	}
	if count == 0 {
		err = errUtil.ErrTodoNotFound{}
	}
	return
}