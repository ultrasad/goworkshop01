package models

import (
	"workshop01/db/mongo"

	"github.com/globalsign/mgo/bson"
)

// Todo is todo
type Todo struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Topic string        `json:"topic" bson:"topic"`
	Done  bool          `json:"done" bson:"done"`
}

// CreateTodo is all todos
func CreateTodo(t *Todo) (*Todo, error) {
	var err error
	conn := mongo.MgoManager().Copy()
	defer conn.Close()

	err = conn.DB("document").C("todo").Insert(&t)
	return t, err
}

// UpdateTodo is all todos
func UpdateTodo(id bson.ObjectId, t *Todo) (*Todo, error) {
	var err error
	conn := mongo.MgoManager().Copy()
	defer conn.Close()

	err = conn.DB("document").C("todo").UpdateId(id, t)
	return t, err
}

// DeleteTodo is all todos
func DeleteTodo(id bson.ObjectId) error {
	var err error
	conn := mongo.MgoManager().Copy()
	defer conn.Close()

	err = conn.DB("document").C("todo").RemoveId(id)
	return err
}

// FindTodoByID is all todos
func FindTodoByID(id bson.ObjectId) (Todo, error) {

	var (
		todo Todo
		err  error
	)

	conn := mongo.MgoManager().Copy()
	defer conn.Close()

	err = conn.DB("document").C("todo").FindId(id).One(&todo)
	return todo, err
}

// FindAllTodos is all todos
func FindAllTodos() ([]Todo, error) {

	var (
		todos []Todo
		err   error
	)

	conn := mongo.MgoManager().Copy()
	defer conn.Close()

	err = conn.DB("document").C("todo").Find(nil).All(&todos)
	return todos, err
}
