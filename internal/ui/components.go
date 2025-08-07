package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

// CreateSectionHeader creates a colored section header using theme colors
func CreateSectionHeader(title string) *fyne.Container {
	background := canvas.NewRectangle(theme.Color(theme.ColorNameDisabled))
	background.SetMinSize(fyne.NewSize(0, 20))

	titleText := canvas.NewText(title, theme.Color(theme.ColorNameForeground))
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.TextSize = 14
	titleText.Alignment = fyne.TextAlignLeading

	headerContent := container.NewVBox(
		container.NewPadded(titleText),
	)
	header := container.NewStack(background, headerContent)

	// Add subtle border below header using theme color
	border := canvas.NewRectangle(theme.Color(theme.ColorNameSeparator))
	border.SetMinSize(fyne.NewSize(0, 1))
	return container.NewVBox(header, border)
}

// FixedSpacer returns a spacer that expands between min and max width/height.
func FixedSpacer(wMin, wMax, hMin, hMax float32) fyne.CanvasObject {
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(wMin, hMin))

	return container.New(&fixedSpacerLayout{
		wMin: wMin, wMax: wMax, hMin: hMin, hMax: hMax,
	}, rect)
}

type fixedSpacerLayout struct {
	wMin, wMax, hMin, hMax float32
}

func (l *fixedSpacerLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	w := size.Width
	h := size.Height
	if l.wMax > 0 && w > l.wMax {
		w = l.wMax
	}
	if w < l.wMin {
		w = l.wMin
	}
	if l.hMax > 0 && h > l.hMax {
		h = l.hMax
	}
	if h < l.hMin {
		h = l.hMin
	}
	if len(objects) > 0 {
		objects[0].Resize(fyne.NewSize(w, h))
	}
}

func (l *fixedSpacerLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(l.wMin, l.hMin)
}
