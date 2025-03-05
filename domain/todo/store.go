package todo

type TodoStore interface {
	GetTodoAll() ([]Todo, error)
	GetTodoByTitle(title string) (Todo, error)
	CreateTodo(title string) (Todo, error)
	UpdateTodoTitle(todoToChange string, title string) (Todo, error)
	UpdateTodoStatus(todoToChange string, completed string) (Todo, error)
}
