package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

func showCalibrationWindow(app fyne.App, onConfirm func(float64)) {
	window := app.NewWindow("Calibration - How Big Image Viewer")
	window.Resize(fyne.NewSize(600, 400))
	window.CenterOnScreen()

	lineLength := float32(100)

	line := canvas.NewLine(color.RGBA{255, 0, 0, 255})
	line.StrokeWidth = 2
	line.Position1 = fyne.NewPos(200, 200)
	line.Position2 = fyne.NewPos(200+lineLength, 200)

	lengthLabel := widget.NewLabel(fmt.Sprintf("Line Length: %.0f pixels", lineLength))

	instructions := widget.NewRichTextFromMarkdown(`
## Calibration Instructions

1. Measure something in real life (e.g., a ruler showing 1 cm)
2. Use arrow keys to adjust the red line to match that measurement on your screen:
   - **Left Arrow**: Decrease line length
   - **Right Arrow**: Increase line length
3. Click **Confirm** when the line matches your real measurement
4. Enter the real-world size this line represents
`)

	realSizeEntry := widget.NewEntry()
	realSizeEntry.SetPlaceHolder("Real size (e.g., 1.0 for 1 cm)")
	realSizeEntry.SetText("1.0")

	confirmBtn := widget.NewButton("Confirm Calibration", func() {
		realSize := 1.0
		if text := realSizeEntry.Text; text != "" {
			if parsed, err := fmt.Sscanf(text, "%f", &realSize); err == nil && parsed == 1 {
				scale := float64(lineLength) / realSize
				onConfirm(scale)
				window.Close()
			}
		}
	})

	content := container.NewVBox(
		instructions,
		container.NewHBox(
			widget.NewLabel("Real size this line represents:"),
			realSizeEntry,
			widget.NewLabel("cm"),
		),
		lengthLabel,
		container.NewWithoutLayout(line),
		confirmBtn,
	)

	window.SetContent(content)

	window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		step := float32(5)

		switch key.Name {
		case fyne.KeyLeft:
			if lineLength > step {
				lineLength -= step
			}
		case fyne.KeyRight:
			if lineLength < 500 {
				lineLength += step
			}
		}

		line.Position2 = fyne.NewPos(200+lineLength, 200)
		line.Refresh()
		lengthLabel.SetText(fmt.Sprintf("Line Length: %.0f pixels", lineLength))
	})

	window.Show()
}