package main

import (
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/puzzle15/internal/puzzle15"
)

const (
	ApplicationId    = "se.softteam.puzzle15"
	ApplicationFlags = glib.APPLICATION_FLAGS_NONE
)

func main() {
	// Create a new application
	application, err := gtk.ApplicationNew(ApplicationId, ApplicationFlags)
	if err != nil {
		panic("Failed to create GTK Application : " + err.Error())
	}

	mainForm := puzzle15.NewMainForm()
	// Hook up the activate event handler
	application.Connect("activate", mainForm.OpenMainForm)
	if err != nil {
		panic("Failed to connect Application.Activate event : " + err.Error())
	}

	// Start the application (and exit when it is done)
	os.Exit(application.Run(nil))
}
