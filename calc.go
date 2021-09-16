//go:generate fyne bundle -o data.go Icon.png

package main

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/Knetic/govaluate"
)

type calc struct {
	equation string

	output  *widget.Label
	buttons map[string]*widget.Button
	window  fyne.Window
}

func (c *calc) display(newtext string) {
	c.equation = newtext
	c.output.SetText(newtext)
}

func (c *calc) character(char rune) {
	c.display(c.equation + string(char))
}

func (c *calc) digit(d int) {
	c.character(rune(d) + '0')
}

func (c *calc) clear() {
	c.display("")
}

func (c *calc) evaluate() {
	expression, err := govaluate.NewEvaluableExpression(c.output.Text)
	if err == nil {
		result, err := expression.Evaluate(nil)
		if err == nil {
			c.display(strconv.FormatFloat(result.(float64), 'f', -1, 64))
		}
	}

	if err != nil {
		log.Println("Error in calculation", err)
		c.display("error")
	}
}

func (c *calc) addButton(text string, action func()) *widget.Button {
	button := widget.NewButton(text, action)
	c.buttons[text] = button
	return button
}

func (c *calc) digitButton(number int) *widget.Button {
	str := strconv.Itoa(number)
	return c.addButton(str, func() {
		c.digit(number)
	})
}

func (c *calc) charButton(char rune) *widget.Button {
	return c.addButton(string(char), func() {
		c.character(char)
	})
}

func (c *calc) onTypedRune(r rune) {
	if r == 'c' {
		r = 'C' // The button is using a capital C.
	}

	if button, ok := c.buttons[string(r)]; ok {
		button.OnTapped()
	}
}

func (c *calc) onTypedKey(ev *fyne.KeyEvent) {
	if ev.Name == fyne.KeyReturn || ev.Name == fyne.KeyEnter {
		c.evaluate()
	} else if ev.Name == fyne.KeyBackspace && len(c.equation) > 0 {
		c.display(c.equation[:len(c.equation)-1])
	}
}

func (c *calc) loadUI(app fyne.App) {
	c.output = &widget.Label{Alignment: fyne.TextAlignTrailing}
	c.output.TextStyle.Monospace = true

	equals := c.addButton("=", c.evaluate)
	equals.Importance = widget.HighImportance

	c.window = app.NewWindow("Calc")
	c.window.SetContent(container.NewGridWithColumns(1,
		c.output,
		container.NewGridWithColumns(4,
			c.addButton("C", c.clear),
			c.charButton('('),
			c.charButton(')'),
			c.charButton('/')),
		container.NewGridWithColumns(4,
			c.digitButton(7),
			c.digitButton(8),
			c.digitButton(9),
			c.charButton('*')),
		container.NewGridWithColumns(4,
			c.digitButton(4),
			c.digitButton(5),
			c.digitButton(6),
			c.charButton('-')),
		container.NewGridWithColumns(4,
			c.digitButton(1),
			c.digitButton(2),
			c.digitButton(3),
			c.charButton('+')),
		container.NewGridWithColumns(2,
			container.NewGridWithColumns(2,
				c.digitButton(0),
				c.charButton('.')),
			equals)),
	)

	c.window.Canvas().SetOnTypedRune(c.onTypedRune)
	c.window.Canvas().SetOnTypedKey(c.onTypedKey)
	c.window.Resize(fyne.NewSize(200, 300))
	c.window.Show()
}

func newCalculator() *calc {
	return &calc{
		buttons: make(map[string]*widget.Button, 19),
	}
}
