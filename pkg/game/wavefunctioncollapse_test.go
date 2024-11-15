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

	t.Run("test initial build board", func(t *testing.T) {
		got := getInitialBuildBoard(3, 3)

		want := make([][]buildCell, 3)
		for i := range want {
			want[i] = make([]buildCell, 3)
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

	t.Run("test initial entropy board", func(t *testing.T) {

		// test in the initial state, and all cards should be possible

		cards := getTestCards()

		want := make([][][]int, 3)
		for i := range want {
			want[i] = make([][]int, 3)
			for j := range want[i] {
				want[i][j] = make([]int, 0)

				for _, card := range cards {
					want[i][j] = append(want[i][j], card.Id)
				}
			}
		}

		got := getEntropyBoard(getInitialBuildBoard(3, 3), cards)

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
		buildBoard := getInitialBuildBoard(3, 3)
		buildBoard[1][1] = buildCell{placed: true, connectors: []Connector{Road, Road, Road, Road}}

		got := getEntropyBoard(buildBoard, getTestCards())

		compareEntropyBoard(t, got, want)

	})

}

func Test_getEntropicCard(t *testing.T) {

	// note i goes down, j goes across

	board := getTestBoardWithCrossInMiddle()

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

	got := countEntropyBoard(getEntropyBoard(getTestBoardWithCrossInMiddle(), getTestCards()))

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

	rules := BasicRules{
		BoardWidth:  3,
		BoardHeight: 3,
	}
	tiles := make([][]Tile, rules.BoardWidth)
	rows := make([]Tile, rules.BoardHeight*rules.BoardWidth)
	for i, startRow := 0, 0; i < rules.BoardWidth; i, startRow = i+1, startRow+rules.BoardHeight {
		endRow := startRow + rules.BoardHeight
		tiles[i] = rows[startRow:endRow:endRow]
	}

	tiles[1][1] = Tile{X: 1, Y: 1, Card: Card{Id: 2, Connectors: []Connector{Road, Road, Road, Road}}}

	g := Game{
		Cards: getTestCards(),
		Rules: rules,
		Board: tiles,
		R:     &TestRnd{},
	}

	initalBuildBoard := getTestBoardWithCrossInMiddle()

	fmt.Println("BB before: ", initalBuildBoard)
	fmt.Println("Board before: ", g.Board)
	evolveBoard(g.R, &initalBuildBoard, getTestCards(), &g.Board)
	fmt.Println("BB after: ", initalBuildBoard)
	fmt.Println("Board after: ", g.Board)
	t.Fail()
}

// ======================== helper functions ========================

type TestRnd struct {
}

func (t TestRnd) Intn(n int) int {
	return 0
}

func printEntropyBoard(title string, entropyBoard [][][]int) {
	fmt.Println(title)
	for _, row := range entropyBoard {
		for _, cell := range row {
			fmt.Print(cell)
		}
		fmt.Println("")
	}

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
					printEntropyBoard("got", got)
					printEntropyBoard("want", want)
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

func getTestBoardWithCrossInMiddle() [][]buildCell {
	board := getInitialBuildBoard(3, 3)
	board[1][1] = buildCell{placed: true, connectors: []Connector{Road, Road, Road, Road}}
	return board
}

func getTestCards() []Card {
	// get a set of cards with connectors and ids that can be used for testing
	allGrass := Card{Id: 1, Connectors: []Connector{Grass, Grass, Grass, Grass}}
	crossRoads := Card{Id: 2, Connectors: []Connector{Road, Road, Road, Road}}
	lRoad := Card{Id: 3, Connectors: []Connector{Road, Road, Grass, Grass}}
	deadEnd := Card{Id: 4, Connectors: []Connector{Road, Grass, Grass, Grass}}

	return []Card{allGrass, crossRoads, lRoad, deadEnd}

}
