package game

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"math"
	"wfc2/pkg/boiler"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Connector int

const connectorCount = 2

const (
	Grass Connector = 1
	Road  Connector = 2
)

type Tile struct {
	Card Card
	X    int
	Y    int
}

// card represents a card that can be place in world
// connectors are always north, east, south, west

type Image struct {
	img         *ebiten.Image
	rotateAngle float64
}

type Card struct {
	Id         int
	Image      Image
	Connectors []Connector
}

type BaseCards struct {
	Filename      string
	ImageLocation []int
	Connectors    string
	Rotations     []int
}

func BuildCards(rules BasicRules, fs fs.FS) []Card {

	cards := make([]Card, 0)

	id := 1
	for _, baseCard := range rules.BaseCards {
		img, _, err := ebitenutil.NewImageFromFileSystem(fs, baseCard.Filename)
		if err != nil {
			fmt.Printf("error reading image file %s: %v", baseCard.Filename, err)
			panic(err)
		}
		image := Image{
			img: img,
		}
		card := Card{
			Id:         id,
			Image:      image,
			Connectors: convertConnections(baseCard.Connectors),
		}
		cards = append(cards, card)
		id++
		for _, rotation := range baseCard.Rotations {
			rotCard := rotateCard(card, rotation, id)
			cards = append(cards, rotCard)
			id++
			cards = append(cards, card)
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
	rotImage := Image{
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
	fmt.Println("rot:", rot)
	for i, c := range connectors {
		rotated[(i+rot)%len(connectors)] = c
	}
	return rotated

}
