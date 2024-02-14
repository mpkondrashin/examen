/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

quota.go

Quota window
*/
package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/mod/semver"

	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/update"
)

type UpdateWindow struct {
	ModalWindow
	version        string
	versionLabel   *widget.Label
	progressBar    *widget.ProgressBar
	downloadButton *widget.Button
}

func NewUpdateWindow(modalWindow ModalWindow) *UpdateWindow {
	s := &UpdateWindow{
		ModalWindow:  modalWindow,
		versionLabel: widget.NewLabel("                                 "),
		progressBar:  widget.NewProgressBar(),
	}
	s.downloadButton = widget.NewButton("Download", s.Download)
	s.downloadButton.Disable()
	s.Reset()
	s.win.SetContent(container.NewPadded(container.NewVBox(s.versionLabel, s.progressBar, s.downloadButton)))
	return s
}

func (s *UpdateWindow) Download() {
	s.downloadButton.Disable()
	fileName := fmt.Sprintf("setup_%s_%s.zip", runtime.GOOS, runtime.GOARCH)
	logging.Debugf("Download: %s", fileName)
	err := update.DownloadRelease(s.version, fileName, globals.DownloadsFolder(), func(p float32) error {
		s.progressBar.SetValue(float64(p))
		return nil
	})
	if err != nil {
		dialog.ShowError(err, s.win)
		logging.LogError(err)
		return
	}
	zipFilePath := filepath.Join(globals.DownloadsFolder(), fileName)
	if err := Unzip(zipFilePath); err != nil {
		dialog.ShowError(err, s.win)
		return
	}
	folder := strings.TrimSuffix(zipFilePath, filepath.Ext(zipFilePath))
	if err := RunOpen(folder); err != nil {
		dialog.ShowError(err, s.win)
		logging.LogError(err)
		return
	}
}

func (s *UpdateWindow) Reset() {
	s.versionLabel.SetText("Checking...")
	s.version = ""
}

func (s *UpdateWindow) Update() {
	s.Reset()
	var err error
	s.version, err = update.LatestVersion(globals.Name)
	if err != nil {
		dialog.ShowError(err, s.win)
		return
	}
	cmp := semver.Compare(s.version, globals.Version)
	logging.Debugf("Compare %s vs %s: %d", s.version, globals.Version, cmp)
	switch cmp {
	case -1:
	case 0:
		s.versionLabel.SetText("You have the newest version")
	case 1:
		s.versionLabel.SetText("New version available: " + s.version)
		s.downloadButton.Enable()
	}
}

func (s *UpdateWindow) Show() {
	s.win.Show()
	go s.Update()
}

func Unzip(zipFilePath string) error {
	reader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer reader.Close()
	folder := strings.TrimSuffix(zipFilePath, filepath.Ext(zipFilePath))
	for _, f := range reader.File {
		err := unzipFile(f, folder)
		if err != nil {
			return err
		}
	}
	return nil
}

func unzipFile(f *zip.File, destination string) error {
	// Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()
	_, err = io.Copy(destinationFile, zippedFile)
	return err
}