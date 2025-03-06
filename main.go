package main

import (
	circle2 "collision/circle"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	_ "image/png"
	"log"
)

const (
	//screenWidth  = 1700
	//screenHeight = 900
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
}

type Runnable interface {
	Update() error
	Draw(screen *ebiten.Image)
}

type Game struct {
	context Runnable
}

func NewGame(circlesNum int) *Game {
	return &Game{
		context: circle2.NewContextGame(circlesNum, screenWidth, screenHeight),
	}
}

//func NewGame(circlesNum int) *Game {
//	return &Game{
//		circles: circle2.NewCircleContainer(circlesNum, screenWidth, screenHeight),
//	}
//}

func (g *Game) Update() error {
	//cont := g.context.GetCurrentIteration()
	//g.context.BeginNewIteration()
	//cont.UpdateConcurrent()
	g.context.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//screen.Fill(color.RGBA{0x99, 0xcc, 0xff, 0xff})
	//cont := g.context.GetCurrentIteration()
	//cont.Draw(screen)
	//for e := g.sprites.Front(); e != nil; e = e.Next() {
	//	s := e.Value.(*sprite)
	//	s.draw(screen)
	//}
	img := ebiten.NewImage(screenWidth, screenHeight)
	img.Fill(color.RGBA{R: 128, G: 222, B: 234, A: 0})
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(img, op)
	g.context.Draw(screen)
	//fmt.Println(g.circles.KineticEnergy())
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nSprites: %d", ebiten.ActualTPS(), g.sprites.Len()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

//func test() {
//	temp := &Game{
//		context: circle2.NewContextEmpty(screenWidth, screenHeight),
//	}
//	var temp_ *circle2.BasicCircle
//	temp_ = circle2.GetInstance().GetCircle(16, screenWidth, screenHeight, 200., 124., -5., 0.)
//	temp.context.State[temp_.GetId()] = temp_
//	temp_ = circle2.GetInstance().GetCircle(16, screenWidth, screenHeight, 300., 124., 4., 0.)
//	temp.context.State[temp_.GetId()] = temp_
//	temp_ = circle2.GetInstance().GetCircle(16, screenWidth, screenHeight, 200., 184., -2., 0.)
//	temp.context.State[temp_.GetId()] = temp_
//	//temp_ = circle2.GetInstance().GetCircle(16, screenWidth, screenHeight, 200., 184., 7., 0.)
//	//temp.context.State[temp_.GetId()] = temp_
//	if err := ebiten.RunGame(temp); err != nil {
//		log.Fatal(err)
//	}
//}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Particles (Ebitengine Demo)")
	//temp := NewGame(0)
	//temp_v2 := circle2.New_(8, screenWidth, screenHeight, 600., 124., -3., 0.)
	//temp.circles.Add(&temp_v2)
	//temp_v2 = circle2.New_(9, screenWidth, screenHeight, 100., 126., 5., 0.)
	//temp.circles.Add(&temp_v2)
	//if err := ebiten.RunGame(temp); err != nil {
	//	log.Fatal(err)
	//}

	if err := ebiten.RunGame(NewGame(50)); err != nil {
		log.Fatal(err)
	}
	//test()
}
