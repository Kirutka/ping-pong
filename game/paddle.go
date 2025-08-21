package game

import (
    "image/color"
    
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Paddle struct {
    x, y      float64
    width     float64
    height    float64
    speed     float64
    isLeft    bool
    score     int
}

func NewPaddle(isLeft bool) *Paddle {
    y := 300.0 - 50 // Центр по вертикали
    x := 50.0
    if !isLeft {
        x = 750.0
    }
    
    return &Paddle{
        x:       x,
        y:       y,
        width:   10,
        height:  100,
        speed:   5,
        isLeft:  isLeft,
        score:   0,
    }
}

func (p *Paddle) Update() {
    // Управление для левой ракетки
    if p.isLeft {
        if ebiten.IsKeyPressed(ebiten.KeyW) && p.y > 0 {
            p.y -= p.speed
        }
        if ebiten.IsKeyPressed(ebiten.KeyS) && p.y < 600-p.height {
            p.y += p.speed
        }
    } else {
        // Управление для правой ракетки (или AI)
        if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && p.y > 0 {
            p.y -= p.speed
        }
        if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && p.y < 600-p.height {
            p.y += p.speed
        }
    }
}

func (p *Paddle) Draw(screen *ebiten.Image) {
    ebitenutil.DrawRect(screen, p.x, p.y, p.width, p.height, color.White)
}

func (p *Paddle) GetRect() (float64, float64, float64, float64) {
    return p.x, p.y, p.width, p.height
}