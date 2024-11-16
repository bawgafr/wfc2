package game

import (
	"fmt"
	"testing"
)

func Test_getInitialBuildBoard(t *testing.T) {

	// initial tests will be done on a 3x3 board
	// rather than the much larger ones that will be used in reality
	// as the rules should be the same, but its far easier to test

	t.Run("test initial build cell", func(t *testing.T) {
		want := buildCell{connectors: []Connector{Grass + Road, Grass + Road, Grass + Road, Grass + Road}}
		got := initialBuildCell

		if want.placed != false {
			t.Errorf("initial build cell not set correctly: got %v, want %v", got, want)
		}

		for i, wantConnector := range want.connectors {
			gotConnector := got.connectors[i]

			if wantConnector != gotConnector {
				t.Errorf("initial build cell not set correctly: got %v, want %v", got, want)
			}
		}
	})

	// test the build board without any seed tiles
	t.Run("test initial build board", func(t *testing.T) {
		g := getTestGame()
		g.Rules.SeedTiles = []SeedTiles{}
		got := getInitialBuildBoard(g)

		want := make([][]buildCell, g.Rules.BoardWidth)
		for i := range want {
			want[i] = make([]buildCell, g.Rules.BoardHeight)
			for j := range want[i] {
				want[i][j] = initialBuildCell
			}
		}

		if len(want) != len(got) {
			t.Errorf("initial board not what in correct state size-wise got %v, want %v", got, want)
		}

		// deepequals and slices.equal have failed me...
		for i, wantrow := range want {
			gotrow := got[i]
			if len(wantrow) != len(gotrow) {
				t.Errorf("initial board not what in correct state row size-wise got %v, want %v", gotrow, wantrow)
			}

			for j, wantcell := range wantrow {
				gotcell := gotrow[j]
				for k, wantConnector := range wantcell.connectors {
					gotConnector := gotcell.connectors[k]

					if wantConnector != gotConnector {
						t.Errorf("initial board not what in correct state got %v, want %v", gotcell, wantcell)
					}
				}

			}
		}

	})

	// test the build board without any seed tiles
	t.Run("test initial build board with centre cross seed", func(t *testing.T) {
		g := getTestGame()
		got := getInitialBuildBoard(g)

		want := make([][]buildCell, g.Rules.BoardWidth)
		for i := range want {
			want[i] = make([]buildCell, g.Rules.BoardHeight)
			for j := range want[i] {
				want[i][j] = initialBuildCell
			}
		}
		want[1][1] = buildCell{placed: true, connectors: []Connector{Road, Road, Road, Road}}

		if len(want) != len(got) {
			t.Errorf("initial board not what in correct state size-wise got %v, want %v", got, want)
		}

		// deepequals and slices.equal have failed me...
		for i, wantrow := range want {
			gotrow := got[i]
			if len(wantrow) != len(gotrow) {
				t.Errorf("initial board not what in correct state row size-wise got %v, want %v", gotrow, wantrow)
			}

			for j, wantcell := range wantrow {
				gotcell := gotrow[j]
				for k, wantConnector := range wantcell.connectors {
					gotConnector := gotcell.connectors[k]

					if wantConnector != gotConnector {
						t.Errorf("initial board not what in correct state got %v, want %v (%d, %d, %d)", gotcell, wantcell, i, j, k)
					}
				}

			}
		}

	})

	t.Run("test blank entropy board", func(t *testing.T) {

		// test in the initial state, and all cards should be possible

		g := getTestGame()
		g.Rules.SeedTiles = []SeedTiles{}

		want := make([][][]int, g.Rules.BoardWidth)
		for i := range want {
			want[i] = make([][]int, g.Rules.BoardHeight)
			for j := range want[i] {
				want[i][j] = make([]int, 0)

				for _, card := range g.Cards {
					want[i][j] = append(want[i][j], card.Id)
				}
			}
		}

		got := getEntropyBoard(getInitialBuildBoard(g), g)

		for i, wantRow := range want {
			gotRow := got[i]
			for j, wantCell := range wantRow {
				gotCell := gotRow[j]
				for k, wantEntropy := range wantCell {
					gotEntropy := gotCell[k]

					if wantEntropy != gotEntropy {
						t.Errorf("initial entropy board not what in correct state got %v, want %v", gotCell, wantCell)
					}
				}
			}

		}

	})

	t.Run("test entropy board with a card already placed", func(t *testing.T) {

		want := [][][]int{
			{
				{1, 2, 3, 4}, {2}, {1, 2, 3, 4},
			},
			{
				{2, 3}, {}, {2},
			},
			{
				{1, 2, 3, 4}, {2, 3, 4}, {1, 2, 3, 4},
			},
		}
		buildBoard := getInitialBuildBoard(getTestGame())
		buildBoard[1][1] = buildCell{placed: true, connectors: []Connector{Road, Road, Road, Road}}

		got := getEntropyBoard(buildBoard, getTestGame())
		getTestCards()
		compareEntropyBoard(t, got, want)

	})

}

func Test_getEntropicCard(t *testing.T) {

	// note i goes down, j goes across

	board := getInitialBuildBoard(getTestGame())

	t.Run("test that the entropic card is built correctly corner cards ", func(t *testing.T) {

		for i := 0; i < 3; i += 2 {
			got := getEntropicCard(board, i, i)

			want := []Connector{Grass + Road, Grass + Road, Grass + Road, Grass + Road}

			compareConnectors(t, got, want)
		}
	})

	t.Run("test that the entropic card is built correctly top middle ", func(t *testing.T) {
		got := getEntropicCard(board, 0, 1)

		want := []Connector{Grass + Road, Grass + Road, Road, Grass + Road}

		compareConnectors(t, got, want)
	})

	t.Run("test that the entropic card is built correctly left centre ", func(t *testing.T) {
		got := getEntropicCard(board, 1, 0)

		want := []Connector{Grass + Road, Road, Grass + Road, Grass + Road}

		compareConnectors(t, got, want)
	})

	t.Run("test that the entropic card is built correctly bottom middle ", func(t *testing.T) {
		got := getEntropicCard(board, 2, 1)

		want := []Connector{Road, Grass + Road, Grass + Road, Grass + Road}

		compareConnectors(t, got, want)
	})

	t.Run("test that the entropic card is built correctly right centre ", func(t *testing.T) {
		got := getEntropicCard(board, 1, 2)

		want := []Connector{Grass + Road, Grass + Road, Grass + Road, Road}

		compareConnectors(t, got, want)
	})

}

func Test_countEntropyBoard(t *testing.T) {
	want := []available{
		{x: 1, y: 0, ids: []int{2}},
		{x: 2, y: 1, ids: []int{2}},
	}
	g := getTestGame()

	got := countEntropyBoard(getEntropyBoard(getInitialBuildBoard(g), g))

	if len(want) != len(got) {

		t.Errorf("different number of low-entropy cards returned, got %v, want %v", got, want)
	}

	for i := range got {
		if want[i].x != got[i].x && want[i].y != got[i].y && len(want[i].ids) != len(got[i].ids) {
			t.Errorf("different low-entropy cards returned, got %v, want %v", got, want)
		}
		for j := range want[i].ids {
			if want[i].ids[j] != got[i].ids[j] {
				t.Errorf("different low-entropy cards ids returned, got %v, want %v", got, want)
			}
		}
	}

}

func Test_evolveBoard(t *testing.T) {

	// we're going to use a fake random number generator to ensure that the same cards are picked
	// it will always return 0 so the first available and the first id from that available will be picked

	g := getTestGame()

	initalBuildBoard := getInitialBuildBoard(g)

	want := [][]buildCell{
		{newBuildCell(false, 3, 3, 3, 3), newBuildCell(true, 2, 2, 2, 2), newBuildCell(false, 3, 3, 3, 3)},
		{newBuildCell(false, 3, 3, 3, 3), newBuildCell(true, 2, 2, 2, 2), newBuildCell(false, 3, 3, 3, 3)},
		{newBuildCell(false, 3, 3, 3, 3), newBuildCell(false, 3, 3, 3, 3), newBuildCell(false, 3, 3, 3, 3)},
	}

	g.evolveBoard(&initalBuildBoard)

	for i, wantRow := range want {
		gotRow := initalBuildBoard[i]
		for j, wantCell := range wantRow {
			gotCell := gotRow[j]
			for k, wantConnector := range wantCell.connectors {
				gotConnector := gotCell.connectors[k]

				if wantConnector != gotConnector {
					t.Errorf("evolved board not what in correct state got %v, want %v", gotCell, wantCell)
				}
			}
		}

	}

}

func Test_processEnds(t *testing.T) {
	t.Run("test that the process ends after 8 iterations with Cross in centre", func(t *testing.T) {
		g := getTestGame()

		buildBoard := getInitialBuildBoard(g)

		cnt := 0
		for {
			println("cnt", cnt)
			printBoard(g.Board)
			if !g.evolveBoard(&buildBoard) {
				break
			}
			cnt++
			println("\n\n")
		}

		if cnt != 8 {
			fmt.Printf("%v", buildBoard)
			t.Errorf("process ends did not end after 8 iterations, got %v", cnt)
		}

	})

	t.Run("test that the process ends after 8 iterations with L in centre", func(t *testing.T) {
		g := getTestGame()
		g.Rules.SeedTiles = []SeedTiles{{1, 1, 3}}
		buildBoard := getInitialBuildBoard(g)

		cnt := 0
		for {
			println("cnt", cnt)
			printBoard(g.Board)
			if !g.evolveBoard(&buildBoard) {
				break
			}
			cnt++
			println("\n\n")
		}

		if cnt != 8 {
			fmt.Printf("%v", buildBoard)
			t.Errorf("process ends did not end after 8 iterations, got %v", cnt)
		}
	})

}

// ======================== helper functions ========================

func getTestGame() *Game {
	rules := BasicRules{
		BoardWidth:  3,
		BoardHeight: 3,
		SeedTiles:   []SeedTiles{{1, 1, 2}},
	}
	tiles := make([][]Tile, rules.BoardWidth)
	rows := make([]Tile, rules.BoardHeight*rules.BoardWidth)
	for i, startRow := 0, 0; i < rules.BoardWidth; i, startRow = i+1, startRow+rules.BoardHeight {
		endRow := startRow + rules.BoardHeight
		tiles[i] = rows[startRow:endRow:endRow]
	}

	tiles[1][1] = Tile{X: 1, Y: 1, Card: Card{Id: 2, Connectors: []Connector{Road, Road, Road, Road}}}

	return &Game{
		Cards: getTestCards(),
		Rules: rules,
		Board: tiles,
		R:     &TestRnd{},
	}
}

func newBuildCell(placed bool, a, b, c, d Connector) buildCell {
	return buildCell{placed: placed, connectors: []Connector{a, b, c, d}}
}

type TestRnd struct {
}

func (t TestRnd) Intn(n int) int {
	return 0
}

func compareEntropyBoard(t *testing.T, got, want [][][]int) {
	t.Helper()
	for i, wantRow := range want {
		gotRow := got[i]
		for j, wantCell := range wantRow {
			gotCell := gotRow[j]
			for k, wantEntropy := range wantCell {
				gotEntropy := gotCell[k]

				if wantEntropy != gotEntropy {
					debugPrintEntropyBoard("got", got)
					debugPrintEntropyBoard("want", want)
					t.Errorf("initial entropy board not what in correct state got %v, want %v", gotCell, wantCell)
				}
			}
		}

	}

}

func compareConnectors(t *testing.T, got, want []Connector) {
	t.Helper()

	if len(want) != len(got) {
		t.Errorf("got %v, want %v", got, want)
	}

	for i := range want {
		if want[i] != got[i] {
			t.Errorf("got %v, want %v", got, want)
		}
	}

}

// func getTestBoardWithCrossInMiddle() [][]buildCell {
// 	board := getInitialBuildBoard(getTestGame())
// 	board[1][1] = buildCell{placed: true, connectors: []Connector{Road, Road, Road, Road}}
// 	return board
// }

func getTestCards() []Card {
	// get a set of cards with connectors and ids that can be used for testing
	allGrass := Card{Id: 1, Connectors: []Connector{Grass, Grass, Grass, Grass}}
	crossRoads := Card{Id: 2, Connectors: []Connector{Road, Road, Road, Road}}
	lRoad := Card{Id: 3, Connectors: []Connector{Road, Road, Grass, Grass}}
	deadEnd := Card{Id: 4, Connectors: []Connector{Road, Grass, Grass, Grass}}

	return []Card{allGrass, crossRoads, lRoad, deadEnd}

}

func PrintBuildBoard(buildBoard [][]buildCell) {

	for _, row := range buildBoard {

		for _, cell := range row {
			if cell.placed {
				print("X")
			} else {
				print("O")
			}
		}
		println()
	}

}

func printBoard(board [][]Tile) {
	for _, row := range board {
		for _, tile := range row {
			print(tile.Card.Id)
		}
		println()
	}
}
