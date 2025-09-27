# How Big Image Viewer

A Go application that displays images in their real-life size on your screen.

## Features

- **Scale Calibration**: Adjust a line to match real-world measurements
- **Real-size Display**: Shows images at their actual physical dimensions
- **Multiple Formats**: Supports JPEG, PNG, GIF, BMP, TIFF
- **Persistent Config**: Saves calibration settings in config.json

## Building on Windows

### Prerequisites
1. Install Go 1.19+ from https://golang.org/dl/
2. Install Visual Studio Build Tools or Visual Studio Community
3. Install pkg-config and development libraries

### Build Steps
```bash
go mod tidy
go build -o how-big-image-viewer.exe
```

## Usage

1. **First Run**: Click "Calibrate Scale"
   - Use arrow keys to adjust the red line
   - Match it to a real measurement (e.g., 1 cm on a ruler)
   - Enter the real size and confirm

2. **View Images**: Click "Open Image"
   - Select any supported image file
   - Image displays at real-life size based on your calibration

## Configuration

The app saves settings in `config.json`:
```json
{
  "scale": 96.0,
  "sizes": ["10x15", "20x30"]
}
```

- `scale`: Pixels per centimeter (calibrated value)
- `sizes`: Predefined common sizes (future feature)

## Architecture

- `main.go`: Application entry point and main window
- `config.go`: Configuration management
- `calibration.go`: Scale calibration window
- `imageviewer.go`: Real-size image display
- `dialogs.go`: File selection dialogs

## WSL Limitations

This GUI application requires X11 forwarding or VcXsrv to run on WSL. For best results, build and run directly on Windows.