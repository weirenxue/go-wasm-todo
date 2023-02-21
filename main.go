package main

import (
	"fmt"
	"go-wasm-todo/model"
	"syscall/js"
)

type DisplayStatus int

const (
	All DisplayStatus = iota
	Active
	Completed
)

var (
	// Get the document object in the browser environment.
	// js.Global() will get either window (in the browser) or global
	// (in Node). Thus, js.Global().Get("document") is equivalent to
	// window.document in browser JavaScript.
	document = js.Global().Get("document")

	// Get the div with id root in DOM.
	// This is equivalent to window.document.getElementById("root").
	root = document.Call("getElementById", "root")

	// A channel used to notify that something has changed.
	changed = make(chan struct{}, 0)

	// What should we display? 'All' for all tasks, 'Active' for
	// unfinished tasks and 'Completed' for completed tasks.
	displayStatus = All
)

var todos = model.Todos{}

// newElementAndAppend creates a new element and appends it to the
// container. Then return the element.
func newElementAndAppend(container js.Value, name string) js.Value {
	// Create an element.
	el := document.Call("createElement", name)
	// Append the new element to the container.
	container.Call("appendChild", el)
	return el
}

// RenderSummary render task summary.
// <div>There are {total} tasks. ({active} Active / {completed} Completed)</div>
func RenderSummary(container js.Value) {
	summaryDiv := newElementAndAppend(container, "div")
	total, active, completed := todos.Summary()
	if total <= 1 {
		summaryDiv.Set("textContent", fmt.Sprintf("There is %d task. (%d Active / %d Completed)", total, active, completed))
	} else {
		summaryDiv.Set("textContent", fmt.Sprintf("There are %d tasks. (%d Active / %d Completed)", total, active, completed))
	}
}

// RenderTodoList render todo list.
// <ol>
//
//	  <li>
//		<input type="checkbot">
//		<button>DELETE</button>
//		<button>EDIT</button>
//		<span>Buy milk</span>
//	  </li>
//
// </ol>
func RenderTodoList(container js.Value) {
	// isShow is a utility function that checks
	// if the task should be displayed.
	isShow := func(todo *model.Todo) bool {
		switch displayStatus {
		case Active:
			return !todo.Completed
		case Completed:
			return todo.Completed
		default:
			return true
		}
	}

	ol := newElementAndAppend(container, "ol")
	for i := range todos {
		index := i
		todo := todos[index]

		if !isShow(todo) {
			continue
		}

		li := newElementAndAppend(ol, "li")

		// Create a checked checkbox element.
		// <input type="checkbox">
		checkBox := newElementAndAppend(li, "input")
		checkBox.Set("type", "checkbox")
		checkBox.Set("checked", todo.Completed)
		// Add a listener to handle value change events.
		checkBox.Call("addEventListener", "change", js.FuncOf(func(this js.Value, args []js.Value) any {
			// Get the event object.
			event := args[0]
			// Set the task Completed property to the value of the checkbox.
			todo.Completed = event.Get("currentTarget").Get("checked").Bool()
			// Inform the app that the screen should be re-rendered.
			changed <- struct{}{}
			return nil
		}))

		// Create a DELETE button.
		// <button>DELETE</button>
		deleteBtn := newElementAndAppend(li, "button")
		deleteBtn.Set("textContent", "DELETE")
		// Add a listener to handle click events.
		deleteBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
			// Remove the task from the task list.
			todos.Remove(index)
			// Inform the app that the screen should be re-rendered.
			changed <- struct{}{}
			return nil
		}))

		// Create an EDIT button.
		// <button>EDIT</button>
		editBtn := newElementAndAppend(li, "button")
		editBtn.Set("textContent", "EDIT")
		// Add a listener to handle click events.
		editBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
			// Displays a prompt asking the user to enter something.
			todoTitle := js.Global().Call("prompt", fmt.Sprintf(`Change '%s' to`, todo.Title), todo.Title).String()
			// Check if something has been entered and is not equal to the origin task title.
			if len(todoTitle) != 0 && todoTitle != todo.Title {
				// Change the task title.
				todo.Title = todoTitle
				// Inform the app that the screen should be re-rendered.
				changed <- struct{}{}
			}
			return nil
		}))

		// Create a text span.
		// <span>{todo.Title}</span>
		span := newElementAndAppend(li, "span")
		span.Set("textContent", todo.Title)
	}
}

// RenderInput render the input form.
// <input> <button>ADD</button>
func RenderInput(container js.Value) {
	// Create an input element for the user to enter a new task title.
	// <input>
	input := newElementAndAppend(container, "input")

	// Create an ADD button.
	// <button>ADD</button>
	btn := newElementAndAppend(container, "button")
	btn.Set("textContent", "ADD")
	// Add a listener to handle click events.
	btn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
		// Get the value of the input element.
		todoTitle := input.Get("value").String()
		// Check if something has been entered.
		if len(todoTitle) > 0 {
			// Create a new task in the task list.
			todos.Add(todoTitle)
			// Clear the input element.
			input.Set("value", "")
			// Inform the app that the screen should be re-rendered.
			changed <- struct{}{}
		}
		return nil
	}))
}

// RenderHeader render the header.
// <h1>Golang Web Assembly for TODO :)</h1>
func RenderHeader(container js.Value) {
	header := newElementAndAppend(container, "h1")
	header.Set("textContent", "Golang Web Assembly for TODO :)")
}

// RenderToggles render the toggles.
// <button class="active">Show All</button>
// <button>Show Active Tasks</button>
// <button>Show Completed Tasks</button>
func RenderToggles(container js.Value) {
	allBtn := newElementAndAppend(container, "button")
	activeBtn := newElementAndAppend(container, "button")
	completedBtn := newElementAndAppend(container, "button")

	// The callback function for the three button clicks.
	// Clear className for each button and add a className "active"
	// for the clicked button.
	commonClickCallback := func(this js.Value) {
		allBtn.Set("className", "")
		activeBtn.Set("className", "")
		completedBtn.Set("className", "")
		this.Set("className", "active")
	}

	// Create a Show All button that has an active class by default.
	// <button class="active">Show All</button>
	allBtn.Set("textContent", "Show All")
	allBtn.Set("className", "active")
	// Add a listener to handle click events.
	allBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
		commonClickCallback(this)
		// Change the display status to "All".
		displayStatus = All
		// Inform the app that the screen should be re-rendered.
		changed <- struct{}{}
		return nil
	}))

	// Create a Show Active Tasks button that has an active class by default.
	// <button class="active">Show Active Tasks</button>
	activeBtn.Set("textContent", "Show Active Tasks")
	// Add a listener to handle click events.
	activeBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
		commonClickCallback(this)
		// Change the display status to "Active".
		displayStatus = Active
		// Inform the app that the screen should be re-rendered.
		changed <- struct{}{}
		return nil
	}))

	// Create a Show Completed Tasks button that has an active class by default.
	// <button class="active">Show Completed Tasks</button>
	completedBtn.Set("textContent", "Show Completed Tasks")
	// Add a listener to handle click events.
	completedBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
		commonClickCallback(this)
		// Change the display status to "Completed".
		displayStatus = Completed
		// Inform the app that the screen should be re-rendered.
		changed <- struct{}{}
		return nil
	}))
}

func main() {
	RenderHeader(newElementAndAppend(root, "header"))
	RenderInput(newElementAndAppend(root, "div"))
	RenderToggles(newElementAndAppend(root, "div"))

	go func() {
		content := newElementAndAppend(root, "div")
		RenderSummary(content)
		RenderTodoList(content)
		for range changed {
			content.Set("innerHTML", "")
			RenderSummary(content)
			RenderTodoList(content)
		}
	}()

	done := make(chan struct{}, 0)
	<-done
}
