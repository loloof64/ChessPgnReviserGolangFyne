package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"

	"image/color"
)

type CustomLightTheme struct{}

func (t CustomLightTheme) BackgroundColor() color.Color {
	return theme.LightTheme().BackgroundColor()
}

func (t CustomLightTheme) ButtonColor() color.Color {
	return theme.LightTheme().ButtonColor()
}

func (t CustomLightTheme) DisabledButtonColor() color.Color {
	return theme.LightTheme().DisabledButtonColor()
}

func (t CustomLightTheme) IconColor() color.Color {
	return theme.LightTheme().IconColor()
}

func (t CustomLightTheme) DisabledIconColor() color.Color {
	return theme.LightTheme().DisabledIconColor()
}

func (t CustomLightTheme) HyperlinkColor() color.Color {
	return theme.LightTheme().HyperlinkColor()
}

func (t CustomLightTheme) TextColor() color.Color {
	return theme.LightTheme().TextColor()
}

func (t CustomLightTheme) DisabledTextColor() color.Color {
	return theme.LightTheme().DisabledTextColor()
}

func (t CustomLightTheme) HoverColor() color.Color {
	return theme.LightTheme().HoverColor()
}

func (t CustomLightTheme) PlaceHolderColor() color.Color {
	return theme.LightTheme().PlaceHolderColor()
}

func (t CustomLightTheme) PrimaryColor() color.Color {
	return theme.LightTheme().PrimaryColor()
}

func (t CustomLightTheme) FocusColor() color.Color {
	return theme.LightTheme().FocusColor()
}

func (t CustomLightTheme) ScrollBarColor() color.Color {
	return theme.LightTheme().ScrollBarColor()
}

func (t CustomLightTheme) ShadowColor() color.Color {
	return theme.LightTheme().ShadowColor()
}

func (t CustomLightTheme) TextSize() int {
	return theme.LightTheme().TextSize()
}

func (t CustomLightTheme) TextFont() fyne.Resource {
	return resourceFreeSerifTtf
}

func (t CustomLightTheme) TextBoldFont() fyne.Resource {
	return theme.LightTheme().TextBoldFont()
}

func (t CustomLightTheme) TextItalicFont() fyne.Resource {
	return theme.LightTheme().TextItalicFont()
}

func (t CustomLightTheme) TextBoldItalicFont() fyne.Resource {
	return theme.LightTheme().TextBoldItalicFont()
}

func (t CustomLightTheme) TextMonospaceFont() fyne.Resource {
	return resourceFreeSerifTtf
}

func (t CustomLightTheme) Padding() int {
	return theme.LightTheme().Padding()
}

func (t CustomLightTheme) IconInlineSize() int {
	return theme.LightTheme().IconInlineSize()
}

func (t CustomLightTheme) ScrollBarSize() int {
	return theme.LightTheme().ScrollBarSize()
}

func (t CustomLightTheme) ScrollBarSmallSize() int {
	return theme.LightTheme().ScrollBarSmallSize()
}

type CustomDarkTheme struct{}

func (t CustomDarkTheme) BackgroundColor() color.Color {
	return theme.DarkTheme().BackgroundColor()
}

func (t CustomDarkTheme) ButtonColor() color.Color {
	return theme.DarkTheme().ButtonColor()
}

func (t CustomDarkTheme) DisabledButtonColor() color.Color {
	return theme.DarkTheme().DisabledButtonColor()
}

func (t CustomDarkTheme) IconColor() color.Color {
	return theme.DarkTheme().IconColor()
}

func (t CustomDarkTheme) DisabledIconColor() color.Color {
	return theme.DarkTheme().DisabledIconColor()
}

func (t CustomDarkTheme) HyperlinkColor() color.Color {
	return theme.DarkTheme().HyperlinkColor()
}

func (t CustomDarkTheme) TextColor() color.Color {
	return theme.DarkTheme().TextColor()
}

func (t CustomDarkTheme) DisabledTextColor() color.Color {
	return theme.DarkTheme().DisabledTextColor()
}

func (t CustomDarkTheme) HoverColor() color.Color {
	return theme.DarkTheme().HoverColor()
}

func (t CustomDarkTheme) PlaceHolderColor() color.Color {
	return theme.DarkTheme().PlaceHolderColor()
}

func (t CustomDarkTheme) PrimaryColor() color.Color {
	return theme.DarkTheme().PrimaryColor()
}

func (t CustomDarkTheme) FocusColor() color.Color {
	return theme.DarkTheme().FocusColor()
}

func (t CustomDarkTheme) ScrollBarColor() color.Color {
	return theme.DarkTheme().ScrollBarColor()
}

func (t CustomDarkTheme) ShadowColor() color.Color {
	return theme.DarkTheme().ShadowColor()
}

func (t CustomDarkTheme) TextSize() int {
	return theme.DarkTheme().TextSize()
}

func (t CustomDarkTheme) TextFont() fyne.Resource {
	return resourceFreeSerifTtf
}

func (t CustomDarkTheme) TextBoldFont() fyne.Resource {
	return theme.DarkTheme().TextBoldFont()
}

func (t CustomDarkTheme) TextItalicFont() fyne.Resource {
	return theme.DarkTheme().TextItalicFont()
}

func (t CustomDarkTheme) TextBoldItalicFont() fyne.Resource {
	return theme.DarkTheme().TextBoldItalicFont()
}

func (t CustomDarkTheme) TextMonospaceFont() fyne.Resource {
	return resourceFreeSerifTtf
}

func (t CustomDarkTheme) Padding() int {
	return theme.DarkTheme().Padding()
}

func (t CustomDarkTheme) IconInlineSize() int {
	return theme.DarkTheme().IconInlineSize()
}

func (t CustomDarkTheme) ScrollBarSize() int {
	return theme.DarkTheme().ScrollBarSize()
}

func (t CustomDarkTheme) ScrollBarSmallSize() int {
	return theme.DarkTheme().ScrollBarSmallSize()
}
