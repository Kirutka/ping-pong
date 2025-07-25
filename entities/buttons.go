package entities

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

type Button struct {
	x, y, width, height float32
	text                string
	color               color.RGBA
	hoverColor          color.RGBA
	isHovered           bool
}

func NewButton(x, y, width, height float32, text string, color, hoverColor color.RGBA) Button {
	return Button{
		x:          x,
		y:          y,
		width:      width,
		height:     height,
		text:       text,
		color:      color,
		hoverColor: hoverColor,
	}
}

func (b *Button) SetHovered(hovered bool) {
	b.isHovered = hovered
}

func DrawButton(screen *ebiten.Image, btn Button, font font.Face) {
	btnColor := btn.color
	if btn.isHovered {
		btnColor = btn.hoverColor
	}

	// Draw button background
	vector.DrawFilledRect(screen, btn.x, btn.y, btn.width, btn.height, btnColor, false)

	// Draw button border
	vector.StrokeRect(screen, btn.x, btn.y, btn.width, btn.height, 2, color.White, false)

	// Draw button text
	textWidth := len(btn.text) * 7
	textX := btn.x + (btn.width-float32(textWidth))/2
	textY := btn.y + btn.height/2 + 5
	text.Draw(screen, btn.text, font, int(textX), int(textY), color.White)
}

func IsInButton(x, y float32, btn Button) bool {
	return x >= btn.x && x <= btn.x+btn.width && y >= btn.y && y <= btn.y+btn.height
}