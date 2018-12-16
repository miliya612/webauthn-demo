package pg

import (
	"bytes"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/miliya612/webauthn-demo/domain/model"
	"github.com/miliya612/webauthn-demo/domain/repo"
)

type credentialRepo struct {
	db *sql.DB
}

var credentials []*model.Credential

func NewCredentialRepo(db *sql.DB) repo.CredentialRepo {
	return credentialRepo{db: db}
}

func (repo credentialRepo) GetByCredentialID(id []byte) (*model.Credential, error) {
	for _, c := range credentials {
		if bytes.Equal(id, c.CredentialID) {
			c.SignCount ++
			return c, nil
		}
	}
	//return nil, errors.New("credential not found")
	return nil, nil
}
func (repo credentialRepo) Create(credential model.Credential) (*model.Credential, error) {
	credentials = append(credentials, &credential)
	return &credential, nil
}
func (repo credentialRepo) Update(credential model.Credential) (*model.Credential, error) {
	return nil, nil
}
func (repo credentialRepo) Delete(id []byte) ([]byte, error) {
	return nil, nil
}

//
//
//func (repo credentialRepo) GetAll() (todos model.Todos, err error) {
//	rows, err := repo.db.Query("select id, name, completed, due from todos")
//	if err != nil {
//		return
//	}
//	defer rows.Close()
//	for rows.Next() {
//		todo := model.Todo{}
//		err = rows.Scan(&todo.ID, &todo.Name, &todo.Completed, &todo.Due)
//		if err != nil {
//			return
//		}
//		todos = append(todos, todo)
//	}
//	return
//}
//
//func (repo credentialRepo) GetByID(id int) (todo model.Todo, err error) {
//	todo = model.Todo{}
//	err = repo.db.QueryRow("select id, name, completed, due from todos where id = $1", id).Scan(&todo.ID, &todo.Name, &todo.Completed, &todo.Due)
//	if err == sql.ErrNoRows {
//		err = errUtil.ErrTodoNotFound{}
//	}
//	return
//}
//
//func (repo credentialRepo) Create(todo model.Todo) (int, error) {
//	stmt, err := repo.db.Prepare("insert into todos (name, due) VALUES ($1, $2) returning id")
//	if err != nil {
//		return -1, err
//	}
//	defer stmt.Close()
//	err = stmt.QueryRow(todo.Name, todo.Due).Scan(&todo.ID)
//	id := todo.ID
//	return id, err
//}
//
//func (repo credentialRepo) Update(todo model.Todo) (model.Todo, error) {
//	_, err := repo.db.Exec("update todos set name = $2, completed = $3, due = $4 where id = $1", todo.ID, todo.Name, todo.Completed, todo.Due)
//	return todo, err
//}
//
//func (repo credentialRepo) Remove(id int) (int, error) {
//	result, err := repo.db.Exec("delete from todos where id = $1", id)
//	count, err := result.RowsAffected()
//	if count == 0 {
//		err = errUtil.ErrTodoNotFound{}
//	}
//	return id, err
//}
