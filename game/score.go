package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Score struct {
    LeftScore  int
    RightScore int
}

func NewScore() *Score {
    return &Score{0, 0}
}

func (s *Score) Draw(screen *ebiten.Image) {
    leftScoreText := fmt.Sprintf("%d", s.LeftScore)
    rightScoreText := fmt.Sprintf("%d", s.RightScore)
    
    text.Draw(screen, leftScoreText, basicfont.Face7x13, 350, 50, color.White)
    text.Draw(screen, rightScoreText, basicfont.Face7x13, 430, 50, color.White)
}