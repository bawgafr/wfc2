package game

import (
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

func Test_convertConnections(t *testing.T) {
	string := "GGRG"
	got := convertConnections(string)

	want := []Connector{Grass, Grass, Road, Grass}

	if len(want) != len(got) {
		t.Errorf("got %v, want %v", got, want)
	}

	for i := range want {
		if want[i] != got[i] {
			t.Errorf("got %v, want %v", got, want)
		}
	}

}

// test with an inital L shape, and rotate it 90, 180, 270 degrees
func Test_rotateConnection(t *testing.T) {
	initial := []Connector{Road, Road, Grass, Grass}

	t.Run("rotate 90", func(t *testing.T) {
		got := rotateConnections(initial, 90)
		want := []Connector{Grass, Road, Road, Grass}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("rotate 180", func(t *testing.T) {
		want := rotateConnections(initial, 180)
		got := []Connector{Grass, Grass, Road, Road}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("rotate 270", func(t *testing.T) {
		got := rotateConnections(initial, 270)
		want := []Connector{Road, Grass, Grass, Road}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func Test_LoadRules(t *testing.T) {
	fs := getFS()
	got := LoadRules("static/rules/basicRules.json", fs)

	want := BasicRules{
		ImageSize:   32,
		BoardWidth:  20,
		BoardHeight: 10,
		BaseCards: []BaseCards{
			{Filename: "", ImageLocation: []int{0, 0, 32, 32}, Connectors: "GGGG", Rotations: []int{}},
			{Filename: "", ImageLocation: []int{0, 0, 32, 32}, Connectors: "RGRG", Rotations: []int{90}},
			{Filename: "", ImageLocation: []int{0, 0, 32, 32}, Connectors: "RRRR", Rotations: []int{}},
			{Filename: "", ImageLocation: []int{0, 0, 32, 32}, Connectors: "RRGG", Rotations: []int{90, 180, 270}},
			{Filename: "", ImageLocation: []int{0, 0, 32, 32}, Connectors: "GGGR", Rotations: []int{90, 180, 270}},
		},
		SeedTiles: []SeedTiles{
			{X: 0, Y: 0, Id: 0},
		},
	}

	if got.ImageSize != want.ImageSize || got.BoardWidth != want.BoardWidth || got.BoardHeight != want.BoardHeight {
		t.Errorf("Basic info not loaded correctly:\n\tsize: %d, %d\n\tw: %d, %d\n\theight: %d, %d",
			got.ImageSize, want.ImageSize, got.BoardWidth, want.BoardWidth, got.BoardHeight, want.BoardHeight)
	}
	if len(got.BaseCards) != len(want.BaseCards) {
		t.Errorf("BaseCards not loaded correctly len of cards: %d, %d", len(got.BaseCards), len(want.BaseCards))
	}

	for i, card := range got.BaseCards {
		if card.Filename != want.BaseCards[i].Filename {
			t.Errorf("Filename not loaded correctly: %s, %s", card.Filename, want.BaseCards[i].Filename)
		}
		if !reflect.DeepEqual(card.ImageLocation, want.BaseCards[i].ImageLocation) {
			t.Errorf("ImageLocation not loaded correctly: %v, %v", card.ImageLocation, want.BaseCards[i].ImageLocation)
		}
		if card.Connectors != want.BaseCards[i].Connectors {
			t.Errorf("Connectors not loaded correctly: %s, %s", card.Connectors, want.BaseCards[i].Connectors)
		}
		if !reflect.DeepEqual(card.Rotations, want.BaseCards[i].Rotations) {
			t.Errorf("Rotations not loaded correctly: %v, %v", card.Rotations, want.BaseCards[i].Rotations)
		}
	}

	if len(got.SeedTiles) != len(want.SeedTiles) {
		t.Errorf("SeedTiles not loaded correctly len of tiles: %d, %d", len(got.SeedTiles), len(want.SeedTiles))
	}

	for i, tile := range got.SeedTiles {
		if tile.X != want.SeedTiles[i].X || tile.Y != want.SeedTiles[i].Y || tile.Id != want.SeedTiles[i].Id {
			t.Errorf("SeedTiles not loaded correctly: %v, %v", tile, want.SeedTiles[i])
		}
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v, want %v", got, want)
	}

}

func Test_BuildCards(t *testing.T) {
	fs := getFS()
	rules := LoadRules("static/rules/basicRules.json", fs)

	got := BuildCards(rules, fs)

	want := []Card{
		{Id: 1, Image: Image{nil, 0.000000}, Connectors: []Connector{Grass, Grass, Grass, Grass}},
		{Id: 2, Image: Image{nil, 0.000000}, Connectors: []Connector{Road, Grass, Road, Grass}},
		{Id: 3, Image: Image{nil, 1.5707963267948966}, Connectors: []Connector{Grass, Road, Grass, Road}},
		{Id: 2, Image: Image{nil, 0.000000}, Connectors: []Connector{Road, Grass, Road, Grass}},
		{Id: 4, Image: Image{nil, 0.000000}, Connectors: []Connector{Road, Road, Road, Road}},
		{Id: 5, Image: Image{nil, 0.000000}, Connectors: []Connector{Road, Road, Grass, Grass}},
		{Id: 6, Image: Image{nil, 1.5707963267948966}, Connectors: []Connector{Grass, Road, Road, Grass}},
		{Id: 5, Image: Image{nil, 0.000000}, Connectors: []Connector{Road, Road, Grass, Grass}},
		{Id: 7, Image: Image{nil, 3.141592653589793}, Connectors: []Connector{Grass, Grass, Road, Road}},
		{Id: 5, Image: Image{nil, 0.000000}, Connectors: []Connector{Road, Road, Grass, Grass}},
		{Id: 8, Image: Image{nil, 4.71238898038469}, Connectors: []Connector{Road, Grass, Grass, Road}},
		{Id: 5, Image: Image{nil, 0.000000}, Connectors: []Connector{Road, Road, Grass, Grass}},
		{Id: 9, Image: Image{nil, 0.000000}, Connectors: []Connector{Grass, Grass, Grass, Road}},
		{Id: 10, Image: Image{nil, 1.5707963267948966}, Connectors: []Connector{Road, Grass, Grass, Grass}},
		{Id: 9, Image: Image{nil, 0.000000}, Connectors: []Connector{Grass, Grass, Grass, Road}},
		{Id: 11, Image: Image{nil, 3.141592653589793}, Connectors: []Connector{Grass, Road, Grass, Grass}},
		{Id: 9, Image: Image{nil, 0.000000}, Connectors: []Connector{Grass, Grass, Grass, Road}},
		{Id: 12, Image: Image{nil, 4.71238898038469}, Connectors: []Connector{Grass, Grass, Road, Grass}},
		{Id: 9, Image: Image{nil, 0.000000}, Connectors: []Connector{Grass, Grass, Grass, Road}},
	}

	if len(want) != len(got) {
		t.Errorf("loaded different number of records got %d, want %d", len(got), len(want))
	}

	for i, card := range got {
		if card.Id != want[i].Id {
			t.Errorf("interation %d, id wrong got %v, want %v", i, got[i].Id, want[i].Id)
		}
		if card.Image.rotateAngle != want[i].Image.rotateAngle {
			t.Errorf("interation %d, angle wrong got %f, want %f, (%v)", i, got[i].Image.rotateAngle, want[i].Image.rotateAngle, got[i].Image.rotateAngle)
		}
		if !reflect.DeepEqual(card.Connectors, want[i].Connectors) {
			wantCon := []Connector{want[i].Connectors[0], want[i].Connectors[1], want[i].Connectors[2], want[i].Connectors[3]}
			gotCon := []Connector{card.Connectors[0], card.Connectors[1], card.Connectors[2], card.Connectors[3]}
			t.Errorf("interation %d, connectors got [%s, %s, %s, %s], want [%s, %s, %s, %s]", i,
				gotCon[0], gotCon[1], gotCon[2], gotCon[3],
				wantCon[0], wantCon[1], wantCon[2], wantCon[3])
		}
	}

	// if !reflect.DeepEqual(want, got) {
	// 	t.Errorf("got %v, want %v", got, want)
	// }
}

///////////////////////////// Helper functions /////////////////////////////////

func getFS() fs.FS {
	fs := fstest.MapFS{
		"static/rules/basicRules.json": {
			Data: []byte(`
				{
					"imageSize": 32,
					"boardWidth": 20,
					"boardHeight": 10,
					"baseCards": [
						{"filename":"", "imageLocation":[0,0,32,32], "connectors":"GGGG", "rotations": []},
						{"filename":"", "imageLocation":[0,0,32,32], "connectors":"RGRG", "rotations": [90]},
						{"filename":"", "imageLocation":[0,0,32,32], "connectors":"RRRR", "rotations": []},
						{"filename":"", "imageLocation":[0,0,32,32], "connectors":"RRGG", "rotations": [90, 180, 270]},
						{"filename":"", "imageLocation":[0,0,32,32], "connectors":"GGGR", "rotations": [90, 180, 270]}
					],
					"seedTiles": [
						{"x": 0, "y": 0, "id":0}
					]
				}`)},
	}

	return fs
}
