package main

import (
	"fmt"
	"math/rand"
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
	leftContainer   fyne.CanvasObject
	rightContainer  fyne.CanvasObject
	bottomContainer fyne.CanvasObject
}

var layout TodoLayout

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

func loadUIbottom() fyne.CanvasObject {
	quote := strings.Split(getQuote(), "--")

	widgetQuote := widget.NewLabelWithStyle(quote[0], fyne.TextAlignCenter, fyne.TextStyle{Italic: true})

	widgetAutor := widget.NewLabelWithStyle(quote[1], fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	return container.NewVBox(widgetQuote, widgetAutor)

}

// Everything Related to the rightContainer
func loadUIright() fyne.CanvasObject {
	return widget.NewLabel("Hello")
}

// Everything Related to the leftContainer

func loadUIleft() fyne.CanvasObject {
	leftContainer := container.NewVBox()
	//toolbarLeft
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentRemoveIcon(), func() {}),
	)

	leftContainer.Add(toolbar)

	return leftContainer
}

func loadUI() fyne.CanvasObject {
	layout = TodoLayout{
		leftContainer:   loadUIleft(),
		rightContainer:  loadUIright(),
		bottomContainer: loadUIbottom(),
	}
	mainContainer := container.NewVSplit(container.NewHSplit(layout.leftContainer, layout.rightContainer), layout.bottomContainer)
	mainContainer.Offset = 0.75
	return mainContainer
}

func main() {

	myApp := app.New()
	myWindow := myApp.NewWindow("ToDoS")
	content := loadUI()

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(640, 480))
	myWindow.ShowAndRun()
}
