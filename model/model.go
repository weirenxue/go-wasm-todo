package model

type Todo struct {
	Title     string
	Completed bool
}

type Todos []*Todo

// Summary is used to summarize the total number of tasks,
// the number of active tasks and the number of completed tasks.
func (t *Todos) Summary() (int, int, int) {
	total := len(*t)
	completedCount := 0
	for _, todo := range *t {
		if todo.Completed {
			completedCount++
		}
	}
	activeCount := total - completedCount
	return total, activeCount, completedCount
}

// Remove an element by specifyng the index of the element.
func (t *Todos) Remove(index int) {
	*t = append((*t)[:index], (*t)[index+1:]...)
}

// Add an element to the last position.
func (t *Todos) Add(title string) {
	*t = append(*t, &Todo{Title: title})
}
