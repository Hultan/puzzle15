package puzzle15

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func (m *MainForm) onDraw(da *gtk.DrawingArea, c *cairo.Context) {
	// Don't draw unless the scramble is complete
	if len(m.tiles) != numberOfTiles {
		return
	}

	for i := 0; i < numberOfTiles; i++ {
		// Don't draw the empty tile
		if m.tiles[m.scramble[i]] == nil {
			continue
		}

		x, y := getXYFromIndex(i)

		c.SetSourceSurface(m.tiles[m.scramble[i]], float64(x*tileWidth), float64(y*tileHeight))
		c.Paint()
		c.SetSourceRGB(0, 0, 0)
		c.SetLineWidth(1)
		c.Rectangle(float64(x*tileWidth), float64(y*tileHeight), float64(tileWidth), float64(tileHeight))
		c.Stroke()
	}
}

func (m *MainForm) onClick(da *gtk.EventBox, event *gdk.Event) bool {
	eventButton := gdk.EventButtonNewFromEvent(event)
	if eventButton.Button() == gdk.BUTTON_PRIMARY {
		x, y := eventButton.X(), eventButton.Y()
		m.makeMove(int(x/float64(tileWidth)), int(y/float64(tileHeight)))
	}
	return false
}
