package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"image/color"
)

type CustomLightTheme struct{}

func (t CustomLightTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (t CustomLightTheme) Font(style fyne.TextStyle) fyne.Resource {
	return resourceFreeSerifTtf
}

func (t CustomLightTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t CustomLightTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

type CustomDarkTheme struct{}

func (t CustomDarkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (t CustomDarkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return resourceFreeSerifTtf
}

func (t CustomDarkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t CustomDarkTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}