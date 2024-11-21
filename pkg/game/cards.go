package game

import (
	"encoding/json"
	"fmt"
	_ "image/png"
	"io/fs"
	"math"
	"wfc2/pkg/boiler"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/exp/rand"
)

type Connector int

const ConnectorCount = 2

const (
	Grass Connector = 1
	Road  Connector = 2
)

func (c Connector) String() string {
	switch c {
	case Grass:
		return "Grass"
	case Road:
		return "Road"
	}
	return fmt.Sprintf("%d", c)
}

type Tile struct {
	Card *Card
	X    int
	Y    int
}

// card represents a card that can be place in world
// connectors are always north, east, south, west

type Image struct {
	img         *ebiten.Image
	rotateAngle float64
}

func (i Image) String() string {
	return fmt.Sprintf("Img rot: %f", i.rotateAngle)
}

type Card struct {
	Id         int
	Image      *Image
	Connectors []Connector
}

func (c Card) String() string {
	str := fmt.Sprintf("{Id: %d, Image: Image{nil, %f}, Connectors: []Connector{%s, %s, %s, %s}},\n", c.Id, c.Image.rotateAngle, c.Connectors[0], c.Connectors[1], c.Connectors[2], c.Connectors[3])
	return str
}

type BaseCards struct {
	Filename      string
	ImageLocation []int
	Connectors    string
	Rotations     []int
}

func BuildCards(rules BasicRules, fs fs.FS) map[int]*Card {

	cards := make(map[int]*Card)

	id := 1
	var img *ebiten.Image
	var err error
	for _, baseCard := range rules.BaseCards {
		if baseCard.Filename != "" {
			img, _, err = ebitenutil.NewImageFromFileSystem(fs, baseCard.Filename)
			if err != nil {
				fmt.Printf("error reading image file %s: %v", baseCard.Filename, err)
				panic(err)
			}
		}
		image := Image{
			img: img,
		}
		card := Card{
			Id:         id,
			Image:      &image,
			Connectors: convertConnections(baseCard.Connectors),
		}
		cards[card.Id] = &card
		id++
		for _, rotation := range baseCard.Rotations {
			rotCard := rotateCard(card, rotation, id)
			cards[rotCard.Id] = &rotCard
			id++
			cards[card.Id] = &card
		}
	}

	return cards
}

func convertConnections(connections string) []Connector {
	var connectors []Connector

	for _, c := range connections {
		switch c {
		case 'G':
			connectors = append(connectors, Grass)
		case 'R':
			connectors = append(connectors, Road)
		}
	}

	return connectors
}

func LoadRules(filename string, fs fs.FS) BasicRules {
	string, err := boiler.ReadJsonFromDisk(fs, filename)
	if err != nil {
		fmt.Println("error reading rules file:", err)
		panic(err)
	}

	var rules BasicRules

	err = json.Unmarshal([]byte(string), &rules)
	if err != nil {
		fmt.Println("error unmarshalling rules file:", err)
		panic(err)
	}

	return rules
}

func rotateCard(card Card, rotation, id int) Card {

	rad := float64(rotation) * math.Pi / 180.00
	rotImage := &Image{
		img:         card.Image.img,
		rotateAngle: rad,
	}

	rotCard := Card{
		Id:         id,
		Image:      rotImage,
		Connectors: rotateConnections(card.Connectors, rotation),
	}
	return rotCard
}

func rotateConnections(connectors []Connector, rotation int) []Connector {
	rotated := make([]Connector, len(connectors))
	rot := rotation / 90
	for i, c := range connectors {
		rotated[(i+rot)%len(connectors)] = c
	}
	return rotated

}

func NewBoard(rules BasicRules, cards map[int]*Card) [][]Tile {
	tiles := make([][]Tile, rules.BoardWidth)
	rows := make([]Tile, rules.BoardHeight*rules.BoardWidth)
	for i, startRow := 0, 0; i < rules.BoardWidth; i, startRow = i+1, startRow+rules.BoardHeight {
		endRow := startRow + rules.BoardHeight
		tiles[i] = rows[startRow:endRow:endRow]
	}

	// add in the seed tiles
	for _, seedTile := range rules.SeedTiles {

		card := cards[seedTile.Id]
		tiles[seedTile.X][seedTile.Y] = Tile{
			X:    seedTile.X,
			Y:    seedTile.Y,
			Card: card,
		}
	}

	return tiles
}

func NewBasicRandom(seed uint64) Rnd {

	s := rand.NewSource(seed)
	return rand.New(s)
}

func (t Tile) Draw(screen *ebiten.Image) {
	if t.Card != nil {
		op := &ebiten.DrawImageOptions{}
		xPos := (t.Y * 32) + 48
		yPos := (t.X * 32) + 48
		op.GeoM.Translate(-16.0, -16.0)
		op.GeoM.Rotate(t.Card.Image.rotateAngle)
		op.GeoM.Translate(float64(xPos), float64(yPos))
		screen.DrawImage(t.Card.Image.img, op)
	}

}
