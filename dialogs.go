package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

func showImageFileDialog(window fyne.Window, callback func(string)) {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		defer reader.Close()

		uri := reader.URI()
		if uri == nil {
			return
		}

		path := uri.Path()
		ext := strings.ToLower(path[strings.LastIndex(path, ".")+1:])

		supportedFormats := map[string]bool{
			"jpg":  true,
			"jpeg": true,
			"png":  true,
			"gif":  true,
			"bmp":  true,
			"tiff": true,
			"tif":  true,
		}

		if !supportedFormats[ext] {
			dialog.ShowError(
				fmt.Errorf("unsupported image format: %s", ext),
				window,
			)
			return
		}

		callback(path)
	}, window)
}

func createImageFileFilter() storage.FileFilter {
	return storage.NewExtensionFileFilter([]string{
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif",
	})
}