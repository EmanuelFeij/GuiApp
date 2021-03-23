package main

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Todo struct {
	Title         string   `json:"title"`
	Todos         []string `json:"todos"`
	WidgetPointer *widget.AccordionItem
}

func (todos *Todo) NewTodo() *widget.AccordionItem {
	container := container.NewVBox()
	for _, t := range todos.Todos {
		container.Add(widget.NewCheck(t, func(a bool) {}))
	}
	todos.WidgetPointer = widget.NewAccordionItem(todos.Title, container)
	return todos.WidgetPointer
}
