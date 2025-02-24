package main

import (
	circle2 "collision/circle"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/png"
	"log"
)

const (
	//screenWidth  = 300
	//screenHeight = 300
	//screenWidth  = 640
	//screenHeight = 480
	//screenWidth  = 1280
	//screenHeight = 720
	screenWidth  = 1900
	screenHeight = 1000
)

func init() {
	// Decode an image from the image file's byte slice.
}

//type sprite struct {
//	count    int
//	maxCount int
//	dir      float64
//
//	img   *ebiten.Image
//	scale float64
//	angle float64
//	alpha float32
//}

type Game struct {
	circles *circle2.CircleContainer
}

func NewGame(circlesNum int) *Game {
	return &Game{
		circles: circle2.NewCircleContainer(circlesNum, screenWidth, screenHeight),
	}
}

func (g *Game) Update() error {
	g.circles.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//screen.Fill(color.RGBA{0x99, 0xcc, 0xff, 0xff})
	g.circles.Draw(screen)
	//for e := g.sprites.Front(); e != nil; e = e.Next() {
	//	s := e.Value.(*sprite)
	//	s.draw(screen)
	//}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nEnergy: %f\n", ebiten.ActualTPS(),
		ebiten.ActualFPS(), g.circles.KineticEnergy()))
	fmt.Println(g.circles.KineticEnergy())
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nSprites: %d", ebiten.ActualTPS(), g.sprites.Len()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Particles (Ebitengine Demo)")
	//temp := NewGame(0)
	//temp_v2 := circle2.New_(32, screenWidth, screenHeight, 400., 128., -0., 5.)
	//temp.circles.Add(&temp_v2)
	//temp_v2 = circle2.New_(16, screenWidth, screenHeight, 100., 128., 5., 0.)
	//temp.circles.Add(&temp_v2)
	//if err := ebiten.RunGame(temp); err != nil {
	//	log.Fatal(err)
	//}
	if err := ebiten.RunGame(NewGame(800)); err != nil {
		log.Fatal(err)
	}
}
