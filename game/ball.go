package game

import (
    "image/color"
    "math/rand"
    
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Ball struct {
    x, y          float64
    width         float64
    height        float64
    velocityX     float64
    velocityY     float64
    speed         float64
}

func NewBall() *Ball {
    return &Ball{
        x:         400,
        y:         300,
        width:     10,
        height:    10,
        velocityX: 1,
        velocityY: 1,
        speed:     4,
    }
}

func (b *Ball) Update(leftPaddle, rightPaddle *Paddle, score *Score) {

    b.x += b.velocityX * b.speed
    b.y += b.velocityY * b.speed
    
    if b.y <= 0 || b.y >= 590 {
        b.velocityY *= -1
    }
    
    if b.checkPaddleCollision(leftPaddle) || b.checkPaddleCollision(rightPaddle) {
        b.velocityX *= -1
        b.speed += 0.1 
    }
    
    if b.x <= 0 {
        score.RightScore++
        b.reset()
    }
    if b.x >= 800 {
        score.LeftScore++
        b.reset()
    }
}

func (b *Ball) checkPaddleCollision(paddle *Paddle) bool {
    px, py, pwidth, pheight := paddle.GetRect()
    return b.x >= px && b.x <= px+pwidth &&
           b.y >= py && b.y <= py+pheight
}

func (b *Ball) reset() {
    b.x = 400
    b.y = 300
    b.speed = 4
    
    if rand.Float64() > 0.5 {
        b.velocityX = 1
    } else {
        b.velocityX = -1
    }
    b.velocityY = rand.Float64()*2 - 1 
}

func (b *Ball) Draw(screen *ebiten.Image) {
    ebitenutil.DrawRect(screen, b.x, b.y, b.width, b.height, color.White)
}