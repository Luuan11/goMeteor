package game

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"go-meteor/src/pkg"
)

const (
	shootCooldown     = time.Millisecond * 500
	bulletSpawnOffset = 50.0
)

type Player struct {
	game *Game

	position Vector
	sprite   *ebiten.Image

	shootCooldown *Timer
}

func NewPlayer(game *Game) *Player {
	sprite := assets.PlayerSprite

	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2

	pos := Vector{
		X: (screenWidth / 2) - halfW,
		Y: (screenHeight) - 170,
	}

	return &Player{
		game:          game,
		position:      pos,
		sprite:        sprite,
		shootCooldown: NewTimer(shootCooldown),
	}
}

func (p *Player) MoveLeft() {
    speed := 6.0
    p.position.X -= speed
    if p.position.X < 0 {
        p.position.X = 0 
    }
}

func (p *Player) MoveRight() {
    speed := 6.0
    p.position.X += speed
    bounds := p.sprite.Bounds()
    maxX := float64(screenWidth) - float64(bounds.Dx())
    if p.position.X > maxX {
        p.position.X = maxX 
    }
}

func (p *Player) MoveUp() {
    speed := 6.0
    p.position.Y -= speed
    if p.position.Y < 0 {
        p.position.Y = 0 
    }
}

func (p *Player) MoveDown() {
    speed := 6.0
    p.position.Y += speed
    bounds := p.sprite.Bounds()
    maxY := float64(screenHeight) - float64(bounds.Dy())
    if p.position.Y > maxY {
        p.position.Y = maxY 
    }
}

func (p *Player) Shoot() {
    if !p.shootCooldown.IsReady() {
        return
    }

    p.shootCooldown.Reset()

    bounds := p.sprite.Bounds()
    halfW := float64(bounds.Dx()) / 2
    halfH := float64(bounds.Dy()) / 2

    spawnPos := Vector{
        p.position.X + halfW,
        p.position.Y - halfH/2,
    }

    bullet := NewLaser(p.game, spawnPos)
    p.game.AddLaser(bullet)

    if p.game.superPowerActive {
        spawnLeftPos := Vector{
            p.position.X - halfW,
            p.position.Y,
        }
        spawnRightPos := Vector{
            p.position.X + halfW*3,
            p.position.Y,
        }

        bulletLeft := NewLaser(p.game, spawnLeftPos)
        bulletRight := NewLaser(p.game, spawnRightPos)
        p.game.AddLaser(bulletLeft)
        p.game.AddLaser(bulletRight)
    }
}

func (p *Player) Update() {
    if ebiten.IsKeyPressed(ebiten.KeyLeft) {
        p.MoveLeft()
    }
    if ebiten.IsKeyPressed(ebiten.KeyRight) {
        p.MoveRight()
    }
    if ebiten.IsKeyPressed(ebiten.KeyUp) {
        p.MoveUp()
    }
    if ebiten.IsKeyPressed(ebiten.KeyDown) {
        p.MoveDown()
    }
    if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
        p.Shoot()
    }

    p.shootCooldown.Update()
}

func (p *Player) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.position.X, p.position.Y)
	screen.DrawImage(p.sprite, op)
}

func (p *Player) Collider() Rect {
	bounds := p.sprite.Bounds()

	return NewRect(
		p.position.X,
		p.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}
