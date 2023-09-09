// Package main launches the calculator app
//
//go:generate fyne bundle -o data.go Icon.png
package main

import "fyne.io/fyne/v2/app"

func main() {
	app := app.New()
	app.SetIcon(resourceIconPng)

	c := newCalculator()
	c.loadUI(app)
	app.Run()
}
