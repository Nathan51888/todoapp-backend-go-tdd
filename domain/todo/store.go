package todo

type TodoStore interface {
	GetTodoByTitle(title string) (Todo, error)
	CreateTodo(title string) (Todo, error)
}
