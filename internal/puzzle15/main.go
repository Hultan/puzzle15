package puzzle15

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/softteam/framework"
)

const applicationTitle = "puzzle15"
const applicationVersion = "v 0.01"
const applicationCopyRight = "©SoftTeam AB, 2020"

var numberOfRows = 3
var numberOfTiles = numberOfRows * numberOfRows
var cardinal = [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
var menuLow, menuMedium, menuHigh *gtk.RadioMenuItem

type MainForm struct {
	window      *gtk.ApplicationWindow
	builder     *framework.GtkBuilder
	aboutDialog *gtk.AboutDialog
	drawingArea *gtk.DrawingArea
	tiles       []*cairo.Surface
	scramble    map[int]int
}

var tileWidth, tileHeight int

// NewMainForm : Creates a new MainForm object
func NewMainForm() *MainForm {
	mainForm := new(MainForm)
	return mainForm
}

// OpenMainForm : Opens the MainForm window
func (m *MainForm) OpenMainForm(app *gtk.Application) {
	// Initialize gtk
	gtk.Init(&os.Args)

	// Create a new softBuilder
	fw := framework.NewFramework()
	builder, err := fw.Gtk.CreateBuilder("main.glade")
	if err != nil {
		panic(err)
	}
	m.builder = builder

	// Get the main window from the glade file
	m.window = m.builder.GetObject("main_window").(*gtk.ApplicationWindow)

	// Set up main window
	m.window.SetApplication(app)
	m.window.SetTitle("Puzzle-15")

	// Hook up the destroy event
	m.window.Connect("destroy", m.window.Close)

	// Quit button
	button := m.builder.GetObject("main_window_quit_button").(*gtk.ToolButton)
	button.Connect("clicked", m.window.Close)

	// Status bar
	statusBar := m.builder.GetObject("main_window_status_bar").(*gtk.Statusbar)
	statusBar.Push(statusBar.GetContextId("gtk-startup"), "Puzzle-15 : version 0.1.0")

	// Drawing area
	m.drawingArea = m.builder.GetObject("drawingArea").(*gtk.DrawingArea)
	m.drawingArea.Connect("draw", m.onDraw)
	event := m.builder.GetObject("drawingAreaEvent").(*gtk.EventBox)
	event.Connect("button-press-event", m.onClick)

	// Menu
	m.setupMenu(fw)

	// Set random seed
	rand.Seed(time.Now().UnixNano())

	// Show the main window
	m.window.ShowAll()
}

func (m *MainForm) setupMenu(fw *framework.Framework) {
	menuNewGame := m.builder.GetObject("mnuFileNewGame").(*gtk.MenuItem)
	menuNewGame.Connect("activate", m.NewGame)

	menuQuit := m.builder.GetObject("mnuFileQuit").(*gtk.MenuItem)
	menuQuit.Connect("activate", m.window.Close)

	menuLow = m.builder.GetObject("mnuDifficultyLow").(*gtk.RadioMenuItem)
	menuLow.Connect("activate", m.SetDifficultyLevel)
	menuMedium = m.builder.GetObject("mnuDifficultyMedium").(*gtk.RadioMenuItem)
	menuMedium.Connect("activate", m.SetDifficultyLevel)
	menuHigh = m.builder.GetObject("mnuDifficultyHigh").(*gtk.RadioMenuItem)
	menuHigh.Connect("activate", m.SetDifficultyLevel)
}

func (m *MainForm) NewGame() {
	fileChooserDlg, err := gtk.FileChooserNativeDialogNew("Please select a puzzle image...", m.window, gtk.FILE_CHOOSER_ACTION_OPEN, "_Open", "_Cancel")
	filter, err := gtk.FileFilterNew()
	if err != nil {
		log.Fatal(err)
	}
	filter.AddPattern("*.png")
	fileChooserDlg.AddFilter(filter)
	fileChooserDlg.SetCurrentFolder("/home/per/code/puzzle15/assets")
	if err != nil {
		log.Fatal("Unable to create fileChooserDlg:", err)
	}
	response := fileChooserDlg.NativeDialog.Run()

	if gtk.ResponseType(response) == gtk.RESPONSE_ACCEPT {
		fileChooser := fileChooserDlg
		filename := fileChooser.GetFilename()

		m.setupNewGame(filename)
	}
}

func (m *MainForm) setupNewGame(filename string) {
	// Load image into surface
	surface, err := cairo.NewSurfaceFromPNG(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Resize form to match image size
	m.drawingArea.SetSizeRequest(surface.GetWidth(), surface.GetHeight())

	// Calculate tile width and height
	tileWidth = surface.GetWidth() / numberOfRows
	tileHeight = surface.GetHeight() / numberOfRows

	// Create tiles
	m.tiles = nil
	for i := 0; i < numberOfTiles-1; i++ {
		x, y := getXYFromIndex(i)
		tileSurface := surface.CreateForRectangle(
			float64(x*tileWidth), float64(y*tileHeight), float64(tileWidth), float64(tileHeight))
		m.tiles = append(m.tiles, tileSurface)
	}
	m.tiles = append(m.tiles, nil)

	m.resetScramble()
	m.Scramble(1000)
}

func (m *MainForm) Scramble(n int) {
	var x, y int
	for i := 0; i < n; i++ {
		ei := m.getEmptyTileIndex()
		ex, ey := getXYFromIndex(ei)
	loop:
		for {
			r := rand.Intn(4)
			c := cardinal[r]
			cx, cy := c[0], c[1]
			if m.isValidMove(ex+cx, ey+cy) {
				x = ex + cx
				y = ey + cy
				break loop
			}
		}
		i1 := getIndexFromXY(ex, ey)
		i2 := getIndexFromXY(x, y)
		m.scramble[i1], m.scramble[i2] = m.scramble[i2], m.scramble[i1]
	}
}

func (m *MainForm) resetScramble() {
	m.scramble = make(map[int]int, numberOfTiles)
	for i := 0; i < numberOfTiles; i++ {
		m.scramble[i] = i
	}
}

func (m *MainForm) getEmptyTileIndex() int {
	for i := 0; i < numberOfTiles; i++ {
		if m.tiles[m.scramble[i]] == nil {
			return i
		}
	}
	panic("ERROR!")
}

func (m *MainForm) makeMove(x int, y int) {
	if !m.isValidMove(x, y) {
		return
	}
	i := getIndexFromXY(x, y)
	ei := m.getEmptyTileIndex()
	m.scramble[i], m.scramble[ei] = m.scramble[ei], m.scramble[i]
	m.drawingArea.QueueDraw()

	if m.isGameWon() {
		noPatternDlg := gtk.MessageDialogNew(m.window, gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, "%s", "You won the game!")
		noPatternDlg.Run()
		noPatternDlg.Destroy()
	}
}

func (m *MainForm) isValidMove(x int, y int) bool {
	if x < 0 || x > numberOfRows-1 || y < 0 || y > numberOfRows-1 {
		return false
	}
	ei := m.getEmptyTileIndex()
	ex, ey := getXYFromIndex(ei)
	return abs(ex-x) == 1 && abs(ey-y) == 0 || abs(ex-x) == 0 && abs(ey-y) == 1
}

func (m *MainForm) isGameWon() bool {
	for i := 0; i < numberOfTiles; i++ {
		if m.scramble[i] != i {
			return false
		}
	}
	return true
}

func (m *MainForm) SetDifficultyLevel() {
	switch {
	case menuLow.GetActive():
		numberOfRows = 3
	case menuMedium.GetActive():
		numberOfRows = 4
	case menuHigh.GetActive():
		numberOfRows = 5
	}

	numberOfTiles = numberOfRows * numberOfRows
}
