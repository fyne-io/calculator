//go:generate fyne bundle -o data.go Icon.png
// Package main launches the calculator app
package main

import "fyne.io/fyne/v2/app"

// see the readme for installation instructions

func main() {
	app := app.New()
	app.SetIcon(resourceIconPng)

	c := newCalculator()
	c.loadUI(app)
	app.Run()
}
