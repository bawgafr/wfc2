package game

import (
	"fmt"
	"io/fs"
	"time"

	"golang.org/x/exp/rand"
)

type Game struct {
	Fs    fs.FS
	Cards map[int]*Card
	Rules BasicRules
	Board [][]Tile
	R     Rnd
	Seed  uint64
}

type Randomiser int

const (
	Basic          Randomiser = iota // ignore the chance field -- the initial version
	SimpleWeighted                   // use the chance field to determine the weight of the card
)

type SeedTiles struct {
	X  int
	Y  int
	Id int
}

type BasicRules struct {
	ImageSize   int
	BoardWidth  int
	BoardHeight int
	BaseCards   []BaseCards
	SeedTiles   []SeedTiles
	Randomiser  Randomiser
}

type Rnd interface {
	Intn(n int) int
}

func NewSeed(seed uint64) *rand.Rand {
	// create the random number generator and seed it
	s := rand.NewSource(seed)
	return rand.New(s)
}

func NewGame(fs fs.FS, seed uint64) *Game {

	rules := LoadRules("static/rules/basicRules.json", fs)
	cards := BuildCards(rules, fs)
	tiles := NewBoard(rules, cards)

	// create the random number generator and seed it
	r := NewSeed(seed)

	g := Game{
		Fs:    fs,
		Cards: cards,
		Rules: rules,
		Board: tiles,
		Seed:  seed,
		R:     r,
	}
	return &g
}

func (g *Game) NewSeed(seed uint64) {
	g.R = NewSeed(seed)
	g.Board = NewBoard(g.Rules, g.Cards)
	g.Seed = seed
	g.CreateLandscape()
}

func (g *Game) Start() {
	buildBoard := getBuildBoard(g)

	g.evolveBoard(&buildBoard)
}

func (g *Game) CreateLandscape() {
	buildBoard := getBuildBoard(g)

	cnt := 0
	startTime := time.Now()
	for i := 0; i < g.Rules.BoardHeight*g.Rules.BoardWidth; i++ {
		cnt++
		if !g.evolveBoard(&buildBoard) {
			break
		}
	}
	endTime := time.Now()
	g.DebugPrintBoard()

	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("%d evolutions of board in %v\n", cnt, elapsedTime)

}
