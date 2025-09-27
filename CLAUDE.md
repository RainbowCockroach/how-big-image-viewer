# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
# Build the application
go build -o how-big-image-viewer.exe

# Install dependencies
go mod tidy

# Development build (with debug info)
go build -o how-big-image-viewer.exe
```

## Application Architecture

This is a Go desktop application using Fyne v2 for GUI that displays images at their real-life physical size on screen.

### Core Components

- **main.go**: Entry point and main window controller. Handles command-line arguments for drag & drop functionality
- **config.go**: Configuration management with JSON persistence. Defines `Config` and `Size` structs
- **calibration.go**: Interactive calibration window where users adjust a red line to match real measurements
- **imageviewer.go**: Core image display with real-size calculations and dynamic size controls
- **dialogs.go**: File selection dialogs with image format validation

### Key Data Flow

1. **Calibration**: User adjusts line → saves pixels/cm scale to config.json
2. **Image Opening**: Load image → calculate real size using scale factor → display at correct pixel dimensions
3. **Size Selection**: Smart filtering shows only relevant paper sizes based on image orientation

### Configuration Structure

The `config.json` uses this structure:
```json
{
  "scale": 45.0,
  "sizes": [
    {
      "name": "A4 x 2▭",
      "width": "21",
      "height": "15"
    }
  ]
}
```

- `scale`: Pixels per centimeter (from calibration)
- `sizes`: Named paper sizes with width/height in cm

### Smart Size Filtering

The image viewer implements orientation-aware size filtering:
- Horizontal images (aspect ratio > 1.1): Shows horizontal + square sizes
- Vertical images (aspect ratio < 0.9): Shows vertical + square sizes
- Square images (0.9 ≤ aspect ratio ≤ 1.1): Shows all sizes

### Drag & Drop Support

The application supports drag & drop by checking `os.Args[1]` for image file paths. If a valid image is dragged onto the executable and scale is calibrated, it opens directly in the viewer.

## Platform Notes

- **Target Platform**: Windows (primary)
- **WSL Limitations**: Requires X11 forwarding for GUI
- **Dependencies**: Fyne v2, requires Go 1.18+
- **Build Requirements**: pkg-config and development libraries on Windows