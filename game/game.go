package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Game struct {
    leftPaddle  *Paddle
    rightPaddle *Paddle
    ball        *Ball
    score       *Score
    paused      bool
}

func NewGame() *Game {
    return &Game{
        leftPaddle:  NewPaddle(true),
        rightPaddle: NewPaddle(false),
        ball:        NewBall(),
        score:       NewScore(),
        paused:      false,
    }
}

func (g *Game) Update() error {
    // Обработка паузы
    if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
        g.paused = !g.paused
    }
    
    if g.paused {
        return nil
    }

    // Обновление состояния игры
    g.leftPaddle.Update()
    g.rightPaddle.Update()
    g.ball.Update(g.leftPaddle, g.rightPaddle, g.score)
    
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    // Очистка экрана
    screen.Fill(color.Black)
    
    // Отрисовка элементов
    g.leftPaddle.Draw(screen)
    g.rightPaddle.Draw(screen)
    g.ball.Draw(screen)
    g.score.Draw(screen)
    
    // Отрисовка разделительной линии
    for y := 0; y < 600; y += 20 {
        ebitenutil.DrawRect(screen, 400-1, float64(y), 2, 10, color.White)
    }
    
    // Отображение паузы
    if g.paused {
        text.Draw(screen, "PAUSED", basicfont.Face7x13, 350, 300, color.White)
    }
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return 800, 600
}