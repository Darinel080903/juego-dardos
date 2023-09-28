package scenes

import (
	"image"
	"my_dart_game/models"
	"os"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	hits            int     = 0
	misses          int     = 0
	darts           int     = 5
	targetX         float64 = 400.0
	gameEnded       bool    = false
	speedLevel      int     = 3
	speeds                  = []float64{0, 5, 10, 20, 40, 60, 80}
	mu              sync.Mutex
	lastSpeedChange time.Time
)

type MainScene struct {
}

func NewMainScene() *MainScene {
	return &MainScene{}
}

func (s *MainScene) Run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Juego de dardos",
		Bounds: pixel.R(0, 0, 800, 600),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	backgroundPic, err := loadPicture("assets/background.png")
	if err != nil {
		panic(err)
	}
	backgroundSprite := pixel.NewSprite(backgroundPic, backgroundPic.Bounds())

	var wg sync.WaitGroup

	targetPic, err := loadPicture("assets/diana.png")
	if err != nil {
		panic(err)
	}
	targetSprite := pixel.NewSprite(targetPic, targetPic.Bounds())

	leftPic, err := loadPicture("assets/left.png")
	if err != nil {
		panic(err)
	}
	leftSprite := pixel.NewSprite(leftPic, leftPic.Bounds())

	rightPic, err := loadPicture("assets/right.png")
	if err != nil {
		panic(err)
	}
	rightSprite := pixel.NewSprite(rightPic, rightPic.Bounds())

	//Go routine 1
	wg.Add(1)
	go models.Render(wg, win, &hits, &misses, &targetX, &speedLevel, &gameEnded, leftSprite, rightSprite, &mu, backgroundSprite, targetSprite)

	//Go routine 2
	wg.Add(1)
	go models.Player(win, &wg, &hits, &mu, &gameEnded, &misses, &darts, &speedLevel, lastSpeedChange, &targetX)

	//Go routine 3
	wg.Add(1)
	go models.SpeedControl(&wg, &mu, &targetX, speeds, &speedLevel)

	wg.Wait()
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
