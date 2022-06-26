//go:generate fyne bundle -o data.go Icon.png

package main

import (
	"log"
	"strconv"
	"strings"

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

func (c *calc) backspace() {
	if len(c.equation) == 0 {
		return
	} else if c.equation == "error" {
		c.clear()
		return
	}

	c.display(c.equation[:len(c.equation)-1])
}

func (c *calc) evaluate() {
	if strings.Contains(c.output.Text, "error") {
		c.display("error")
		return
	}

	expression, err := govaluate.NewEvaluableExpression(c.output.Text)
	if err != nil {
		log.Println("Error in calculation", err)
		c.display("error")
		return
	}

	result, err := expression.Evaluate(nil)
	if err != nil {
		log.Println("Error in calculation", err)
		c.display("error")
		return
	}

	value, ok := result.(float64)
	if !ok {
		log.Println("Invalid input:", c.output.Text)
		c.display("error")
		return
	}

	c.display(strconv.FormatFloat(value, 'f', -1, 64))
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
	} else if ev.Name == fyne.KeyBackspace {
		c.backspace()
	}
}

func (c *calc) onPasteShortcut(shortcut fyne.Shortcut) {
	content := shortcut.(*fyne.ShortcutPaste).Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err != nil {
		return
	}

	c.display(c.equation + content)
}

func (c *calc) onCopyShortcut(shortcut fyne.Shortcut) {
	shortcut.(*fyne.ShortcutCopy).Clipboard.SetContent(c.equation)
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

	canvas := c.window.Canvas()
	canvas.SetOnTypedRune(c.onTypedRune)
	canvas.SetOnTypedKey(c.onTypedKey)
	canvas.AddShortcut(&fyne.ShortcutCopy{}, c.onCopyShortcut)
	canvas.AddShortcut(&fyne.ShortcutPaste{}, c.onPasteShortcut)
	c.window.Resize(fyne.NewSize(200, 300))
	c.window.Show()
}

func newCalculator() *calc {
	return &calc{
		buttons: make(map[string]*widget.Button, 19),
	}
}
