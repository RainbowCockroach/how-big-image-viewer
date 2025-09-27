package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("com.imageviewer.howbig")

	// Check if an image file was passed as argument (drag & drop)
	if len(os.Args) > 1 {
		imagePath := os.Args[1]
		if isImageFile(imagePath) {
			config, err := loadConfig()
			if err != nil || config.Scale <= 0 {
				// Show calibration first if not calibrated
				showMainWindow(a)
			} else {
				// Open image directly
				err := showImageViewer(a, imagePath, config.Scale)
				if err != nil {
					showMainWindow(a)
				}
			}
		} else {
			showMainWindow(a)
		}
	} else {
		showMainWindow(a)
	}

	a.Run()
}

func isImageFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	supportedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".tiff": true,
		".tif":  true,
	}
	return supportedExts[ext]
}

func showMainWindow(app fyne.App) {
	window := app.NewWindow("How Big Image Viewer")
	window.Resize(fyne.NewSize(400, 300))
	window.CenterOnScreen()

	config, err := loadConfig()
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	scaleLabel := widget.NewLabel("Current scale: Not calibrated")
	if config.Scale > 0 {
		scaleLabel.SetText(fmt.Sprintf("Current scale: %.2f pixels/cm", config.Scale))
	}

	calibrateBtn := widget.NewButton("Calibrate Scale", func() {
		showCalibrationWindow(app, func(scale float64) {
			config.Scale = scale
			config.save()
			scaleLabel.SetText(fmt.Sprintf("Current scale: %.2f pixels/cm", scale))
		})
	})

	openImageBtn := widget.NewButton("Open Image", func() {
		if config.Scale <= 0 {
			dialog.ShowInformation("Calibration Required",
				"Please calibrate the scale first before opening images.", window)
			return
		}

		showImageFileDialog(window, func(imagePath string) {
			err := showImageViewer(app, imagePath, config.Scale)
			if err != nil {
				dialog.ShowError(err, window)
			}
		})
	})

	instructions := widget.NewRichTextFromMarkdown(`
# How Big Image Viewer

This application displays images in their real-life size on your screen.

## Usage:
1. **Calibrate** - Set the scale by matching a line to a real measurement
2. **Open Image** - View any image in its actual physical size

**Supported formats:** JPEG, PNG, GIF, BMP, TIFF
`)

	content := container.NewVBox(
		instructions,
		scaleLabel,
		calibrateBtn,
		openImageBtn,
	)

	window.SetContent(content)
	window.Show()
}