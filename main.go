package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gocolly/colly"
)

type TodoLayout struct {
	leftContainer   *fyne.Container
	rightContainer  *fyne.Container
	bottomContainer *fyne.Container
}

var (
	layout       TodoLayout
	myApp        fyne.App
	currentTasks []string
)

// Everything Related to the Bottom Container
func getRandomNumber(max int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(max)

}
func getQuote() string {

	liS := make([]string, 0)

	c := colly.NewCollector()

	c.OnHTML(`div[class=gate-check]`, func(e *colly.HTMLElement) {

		e.ForEach("li", func(i int, e *colly.HTMLElement) {
			liS = append(liS, string(e.Text))
		})

	})

	err := c.Visit("https://www.entrepreneur.com/article/247213/")
	if err != nil {
		fmt.Println(err)
	}

	return liS[getRandomNumber(len(liS))]
}

func loadUIbottom() *fyne.Container {
	newQuote := getQuote()
	quoteToDisplay := strings.Split(newQuote, "--")
	var widgetAutor *widget.Label
	if len(quoteToDisplay) == 1 {
		quoteToDisplay = strings.Split(newQuote, "—")
		widgetAutor = widget.NewLabelWithStyle(quoteToDisplay[1], fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	} else {
		widgetAutor = widget.NewLabelWithStyle(quoteToDisplay[1], fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	}

	widgetQuote := widget.NewLabelWithStyle(quoteToDisplay[0], fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
	widgetQuote.Wrapping = fyne.TextWrapBreak
	return container.NewVBox(widgetQuote, widgetAutor)

}

// Everything Related to the leftContainer

func loadUIleft() *fyne.Container {

	titleContainer := widget.NewCard("ToDos", "Here to Help you Achieve Everything", nil)

	//titleContainer.Title = TextAlignCenter
	// main leftContainer
	leftContainer := container.NewVBox(titleContainer)

	//title left

	//toolbarLeft

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {

			newWindow := myApp.NewWindow("Add a Todo")
			newWindow.Resize(fyne.NewSize(300, 300))

			var textArea *fyne.Container
			entry := &widget.Entry{PlaceHolder: "Add a Day"}
			textArea = container.NewVBox(

				widget.NewToolbar(
					widget.NewToolbarAction(theme.ContentAddIcon(), func() {
						widg := &widget.Entry{
							PlaceHolder: "Add Todo!",
						}

						textArea.Add(widg)

					},
					),
				),
			)

			var title string
			var todo []string

			form := &widget.Form{
				Items: []*widget.FormItem{
					{Widget: entry},
					{Widget: textArea},
				},

				OnSubmit: func() {
					title = fmt.Sprint(entry.Text)
					for _, v := range textArea.Objects {
						value, ok := v.(*widget.Entry)
						if ok {
							todo = append(todo, value.Text)
						}

					}

					newWindow.Close()

					t := NewTodo(title, todo)
					if t != nil {
						leftContainer.Add(t)
					}

				},
			}

			newWindow.SetContent(form)
			newWindow.Show()

		}),
		widget.NewToolbarAction(theme.ContentRemoveIcon(), func() {

			newWindow := myApp.NewWindow("Remove a Todo")
			newWindow.Resize(fyne.NewSize(300, 300))

			removeContainer := container.NewVBox()

			selectTodoToRemove := widget.NewSelect(
				currentTasks,
				func(changed string) {
					for _, item := range layout.leftContainer.Objects {
						acc, ok := item.(*widget.Accordion)
						if ok {
							for _, accItem := range acc.Items {
								if accItem.Title == changed {
									layout.leftContainer.Remove(item)
								}
							}
						}
					}
					newWindow.Close()
				},
			)
			removeContainer.Add(selectTodoToRemove)
			newWindow.SetContent(removeContainer)
			newWindow.Show()

		}),
	)

	leftContainer.Add(toolbar)

	return leftContainer
}

func removeEmptyAccordions() {
	for _, lObject := range layout.leftContainer.Objects {
		acc, ok := lObject.(*widget.Accordion)
		if ok {

			for _, accItem := range acc.Items {
				isContainer, ok := accItem.Detail.(*fyne.Container)
				if ok && len(isContainer.Objects) == 0 {
					layout.leftContainer.Remove(acc)
				}
			}
		}
	}
}

func NewTodo(title string, todos []string) *widget.Accordion {
	if title == "" {
		timeNow := time.Now()
		title = fmt.Sprintf("%v/%v/%v", timeNow.Day(), timeNow.Month(), timeNow.Year())
	}
	currentTasks = append(currentTasks, title)

	checkContainer := container.NewVBox()
	for _, t := range todos {
		checkContainer.Add(widget.NewCheck(t, func(a bool) {
			for _, t := range checkContainer.Objects {
				ch, ok := t.(*widget.Check)

				if ok && ch.Checked {

					flag := false
					for _, righObject := range layout.rightContainer.Objects {

						accItem, ok := righObject.(*widget.Accordion)
						if ok {
							for _, child := range accItem.Items {
								if child.Title == title {
									flag = true
									cont, ok := child.Detail.(*fyne.Container)
									if ok {
										cont.Add(widget.NewLabel(ch.Text))
									}
								}

							}
						}

					}
					if !flag {
						layout.rightContainer.Add(widget.NewAccordion(widget.NewAccordionItem(title, container.NewVBox(widget.NewLabel(ch.Text)))))
					}
					checkContainer.Remove(t)
					removeEmptyAccordions()
				}
			}
		}))
	}
	if len(checkContainer.Objects) == 0 {
		return nil
	}
	return widget.NewAccordion(widget.NewAccordionItem(title, checkContainer))
}
func loadUIright() *fyne.Container {

	return container.NewVBox(widget.NewCard("ToDosDone", "Great Work", nil))
}

func loadUI() fyne.CanvasObject {
	layout = TodoLayout{
		leftContainer:   loadUIleft(),
		rightContainer:  loadUIright(),
		bottomContainer: loadUIbottom(),
	}
	mainContainer := container.NewVSplit(container.NewHSplit(container.NewVScroll(layout.leftContainer), container.NewVScroll(layout.rightContainer)), layout.bottomContainer)
	mainContainer.Offset = 0.75
	return mainContainer
}

func newIcon() *fyne.StaticResource {
	icon, err := os.OpenFile("img/icon.jpg", os.O_RDONLY, 0644)
	if err != nil {
		fmt.Print(err)
	}
	iconByteSlice := make([]byte, 60000)
	_, err = icon.Read(iconByteSlice)
	if err != nil {
		fmt.Println(err)
	}
	return fyne.NewStaticResource("icon.jpg", iconByteSlice)
}

func main() {

	myApp = app.New()
	currentTasks = make([]string, 0)

	myApp.SetIcon(newIcon())

	myWindow := myApp.NewWindow("ToDoS")
	content := loadUI()

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(640, 480))
	myWindow.ShowAndRun()
}
