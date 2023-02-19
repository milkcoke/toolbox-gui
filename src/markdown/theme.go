package markdown

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type CustomTheme struct{}

func (customTheme *CustomTheme) Color(themeColorName fyne.ThemeColorName, themeVariant fyne.ThemeVariant) color.Color {
	switch themeColorName {
	case theme.ColorNameBackground:
		if themeVariant == theme.VariantLight {
			return &color.NRGBA{R: 0xF0, G: 0xE9, B: 0x9B, A: 0xFF}
		}
		return &color.NRGBA{R: 0x37, G: 0x2B, B: 0x09, A: 0xFF}
	case theme.ColorNameForeground:
		if themeVariant == theme.VariantLight {
			return &color.NRGBA{R: 0x46, G: 0x3A, B: 0x11, A: 0xFF}
		}
		return &color.NRGBA{R: 0xF0, G: 0xE9, B: 0x9B, A: 0xFF}
	case theme.ColorNamePrimary:
		return &color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xAA}
	case theme.ColorNameFocus:
		return &color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x66}
	case theme.ColorNameInputBackground:
		// removes background from entry
		return color.Transparent
		//return &color.RGBA{R: 198, G: 210, B: 16, A: 75}
	}

	return theme.DefaultTheme().Color(themeColorName, themeVariant)
}

func (customTheme *CustomTheme) Font(s fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(s)
}

func (customTheme *CustomTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (customTheme *CustomTheme) Size(n fyne.ThemeSizeName) float32 {
	if n == theme.SizeNameText {
		return theme.DefaultTheme().Size(n) + 2
	}

	return theme.DefaultTheme().Size(n)
}
