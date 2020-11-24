// Package main launches the calculator app
package main

import "fyne.io/fyne/app"

func main() {
	app := app.New()
	app.SetIcon(resourceIconPng)

	c := newCalculator()
	c.loadUI(app)
	app.Run()
}
