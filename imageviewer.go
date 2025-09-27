package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func showImageViewer(app fyne.App, imagePath string, scale float64) error {
	window := app.NewWindow("Image Viewer - How Big Image Viewer")

	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	bounds := img.Bounds()
	originalWidth := float64(bounds.Dx())
	originalHeight := float64(bounds.Dy())

	// Load config for predefined sizes
	config, _ := loadConfig()

	// Current display size state
	currentDisplayWidth := originalWidth
	currentDisplayHeight := originalHeight

	canvasImg := canvas.NewImageFromImage(img)
	canvasImg.FillMode = canvas.ImageFillOriginal

	// Info label that will be updated
	var infoLabel *widget.RichText

	// Function to update display
	updateDisplay := func() {
		canvasImg.Resize(fyne.NewSize(float32(currentDisplayWidth), float32(currentDisplayHeight)))

		infoText := fmt.Sprintf(
			"File: %s\nOriginal: %.0f x %.0f pixels\nCurrent Display: %.0f x %.0f pixels (%.1f x %.1f cm)\nScale: %.2f pixels/cm",
			imagePath,
			originalWidth, originalHeight,
			currentDisplayWidth, currentDisplayHeight,
			currentDisplayWidth/scale, currentDisplayHeight/scale,
			scale,
		)
		infoLabel.ParseMarkdown("```\n" + infoText + "\n```")
		canvasImg.Refresh()
	}

	infoLabel = widget.NewRichText()

	// Size controls
	sizeLabel := widget.NewLabel("Display Size:")

	// Determine image orientation
	imageAspectRatio := originalWidth / originalHeight
	var imageOrientation string
	if imageAspectRatio > 1.1 {
		imageOrientation = "horizontal"
	} else if imageAspectRatio < 0.9 {
		imageOrientation = "vertical"
	} else {
		imageOrientation = "square"
	}

	// Filter sizes based on image orientation
	sizeOptions := []string{"Original Size"}
	for _, size := range config.Sizes {
		width, _ := strconv.ParseFloat(size.Width, 64)
		height, _ := strconv.ParseFloat(size.Height, 64)
		sizeAspectRatio := width / height

		var sizeOrientation string
		if sizeAspectRatio > 1.1 {
			sizeOrientation = "horizontal"
		} else if sizeAspectRatio < 0.9 {
			sizeOrientation = "vertical"
		} else {
			sizeOrientation = "square"
		}

		// Show size based on image orientation
		shouldShow := false
		switch imageOrientation {
		case "horizontal":
			shouldShow = (sizeOrientation == "horizontal" || sizeOrientation == "square")
		case "vertical":
			shouldShow = (sizeOrientation == "vertical" || sizeOrientation == "square")
		case "square":
			shouldShow = true // Show all orientations for square images
		}

		if shouldShow {
			sizeOptions = append(sizeOptions, size.Name)
		}
	}
	sizeOptions = append(sizeOptions, "Custom...")

	sizeSelect := widget.NewSelect(sizeOptions, func(selected string) {
		if selected == "Original Size" {
			// Show at real-life size based on scale
			currentDisplayWidth = originalWidth
			currentDisplayHeight = originalHeight
			updateDisplay()
		} else if selected == "Custom..." {
			// Will be handled by custom entry
		} else {
			// Find the selected size from config
			for _, size := range config.Sizes {
				if size.Name == selected {
					if width, err1 := strconv.ParseFloat(size.Width, 64); err1 == nil {
						if height, err2 := strconv.ParseFloat(size.Height, 64); err2 == nil {
							currentDisplayWidth = width * scale
							currentDisplayHeight = height * scale
							updateDisplay()
						}
					}
					break
				}
			}
		}
	})
	sizeSelect.SetSelected("Original Size")

	// Custom size entry
	customSizeEntry := widget.NewEntry()
	customSizeEntry.SetPlaceHolder("e.g., 15x10 (width x height in cm)")
	customSizeEntry.Resize(fyne.NewSize(250, customSizeEntry.MinSize().Height))

	applyCustomBtn := widget.NewButton("Apply Custom Size", func() {
		text := customSizeEntry.Text
		parts := strings.Split(text, "x")
		if len(parts) == 2 {
			if width, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64); err1 == nil {
				if height, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64); err2 == nil {
					// Convert cm to pixels using scale
					currentDisplayWidth = width * scale   // cm * pixels/cm = pixels
					currentDisplayHeight = height * scale // cm * pixels/cm = pixels
					updateDisplay()
					sizeSelect.SetSelected("Custom...")
				}
			}
		}
	})

	scrollContainer := container.NewScroll(canvasImg)
	scrollContainer.SetMinSize(fyne.NewSize(800, 600))

	// Controls panel
	controlsPanel := container.NewVBox(
		infoLabel,
		widget.NewSeparator(),
		container.NewHBox(sizeLabel, sizeSelect),
		container.NewVBox(
			widget.NewLabel("Custom size (width x height in cm):"),
			container.NewHBox(customSizeEntry, applyCustomBtn),
		),
	)

	content := container.NewBorder(
		controlsPanel,
		nil,
		nil,
		nil,
		scrollContainer,
	)

	// Initialize display
	updateDisplay()

	window.SetContent(content)
	window.Resize(fyne.NewSize(900, 700))
	window.CenterOnScreen()
	window.Show()

	return nil
}