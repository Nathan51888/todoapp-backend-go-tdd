package todo

type TodoStore interface {
	GetTodoAll() ([]Todo, error)
	GetTodoByTitle(title string) (Todo, error)
	GetTodoById(todoId int) (Todo, error)
	CreateTodo(title string) (Todo, error)
	UpdateTodoTitle(todoId int, title string) (Todo, error)
	UpdateTodoStatus(todoId int, completed bool) (Todo, error)
	UpdateTodoById(todoId int, changedTodo Todo) (Todo, error)
	DeleteTodoById(todoId int) (Todo, error)
}
