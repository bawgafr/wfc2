package game

import "fmt"

const fullConnector = Grass + Road

var initialBuildCell = buildCell{placed: false, connectors: []Connector{Grass + Road, Grass + Road, Grass + Road, Grass + Road}}

type buildCell struct {
	placed     bool
	connectors []Connector
}

type available struct {
	x   int
	y   int
	ids []int
}

func getBuildBoard(g *Game) [][]buildCell {
	board := make([][]buildCell, g.Rules.BoardWidth)
	for i := range board {
		board[i] = make([]buildCell, g.Rules.BoardHeight)
		for j := range board[i] {
			if g.Board[i][j].Card != nil {
				// if the cell is already placed
				board[i][j] = buildCell{placed: true, connectors: g.Board[i][j].Card.Connectors}
			} else {
				board[i][j] = initialBuildCell
			}
		}
	}

	return board
}

// loop through all of the buildCells in the board and set the entropy to the number of possible connectors
// we need to test the card's connectors against the the surrounding cells connectors
func getEntropyBoard(board [][]buildCell, g *Game) [][][]int {

	entropyBoard := make([][][]int, len(board))

	for i, row := range board {

		entropyBoard[i] = make([][]int, len(row))

		for j, cell := range row {
			placed := cell.placed
			buildCell := buildCell{connectors: getEntropicCard(board, i, j)}

			if !placed {
				// compare the built cell to all of the cards
				for l := 1; l <= len(g.Cards); l++ {
					card := g.Cards[l]
					match := true
					for k := 0; k < 4; k++ {
						if (buildCell.connectors[k] & card.Connectors[k]) == 0 {
							match = false
							break
						}

					}
					if match {
						entropyBoard[i][j] = append(entropyBoard[i][j], card.Id)
					}
				}
			} else {
				entropyBoard[i][j] = []int{}
			}

		}
	}

	return entropyBoard

}

func getEntropicCard(board [][]buildCell, i, j int) []Connector {
	n, e, s, w := fullConnector, fullConnector, fullConnector, fullConnector
	row := len(board[0])

	if i > 0 {
		n = board[i-1][j].connectors[2]
	}
	if j < row-1 {
		e = board[i][j+1].connectors[3]
	}
	if i < len(board)-1 {
		s = board[i+1][j].connectors[0]
	}
	if j > 0 {
		w = board[i][j-1].connectors[1]
	}
	return []Connector{n, e, s, w}

}

// return a list of the locations that have the fewest possible cards
func countEntropyBoard(entropyBoard [][][]int) []available {
	smallestPossible := 100000
	availableCells := []available{}
	for i, row := range entropyBoard {
		for j, cell := range row {
			l := len(cell)
			if l < smallestPossible && l > 0 {
				smallestPossible = len(cell)
				availableCells = []available{{x: i, y: j, ids: cell}}
			} else if len(cell) == smallestPossible {
				availableCells = append(availableCells, available{x: i, y: j, ids: cell})
			}
		}
	}
	return availableCells
}

func (g *Game) evolveBoard(buildBoard *[][]buildCell) bool {

	// the whole loop

	entropyBoard := getEntropyBoard(*buildBoard, g)

	// get a list of the lowest entropy
	availableCards := countEntropyBoard(entropyBoard)

	if len(availableCards) == 0 {
		return false
	}

	// randomisation functionality

	// select a random location to play a card -- from the list of locations with the fewest range of cards
	selectedAvaialable := availableCards[g.R.Intn(len(availableCards))]

	// select a random id from the available
	var selectedCardId int
	switch g.Rules.Randomiser {
	case Basic:
		selectedCardId = selectedAvaialable.ids[g.R.Intn(len(selectedAvaialable.ids))]
	case SimpleWeighted:
		selectedCardId = basicWeightedRandom(g, selectedAvaialable.ids)
	}

	// place the card in the buildBoard
	(*buildBoard)[selectedAvaialable.x][selectedAvaialable.y] = buildCell{placed: true, connectors: g.Cards[selectedCardId].Connectors}

	// place the card on the board
	g.Board[selectedAvaialable.x][selectedAvaialable.y] = Tile{Card: g.Cards[selectedCardId], X: selectedAvaialable.x, Y: selectedAvaialable.y}

	return true
}

func basicWeightedRandom(g *Game, ids []int) int {
	// get the total weight of the cards
	totalWeight := 0
	for _, id := range ids {
		totalWeight += g.Cards[id].chance
	}

	// get a random number between 0 and the total weight
	r := g.R.Intn(totalWeight)

	// loop through the ids and subtract the weight from the random number until it is less than 0
	// then return that id
	for _, id := range ids {
		r -= g.Cards[id].chance
		if r < 0 {
			return id
		}
	}

	// must be the first one then...
	return ids[0]
}

func debugPrintEntropyBoard(title string, entropyBoard [][][]int) {
	fmt.Println(title)
	for _, row := range entropyBoard {
		for _, cell := range row {
			fmt.Print(cell)
		}
		fmt.Println("")
	}

}
func (g Game) DebugPrintBoard() {
	fmt.Println("** Board **")
	fmt.Printf("Board width: %d, Board height: %d\n", g.Rules.BoardWidth, g.Rules.BoardHeight)
	fmt.Println()
	for _, row := range g.Board {
		for _, tile := range row {
			if tile.Card != nil {
				fmt.Printf("[%02d]", tile.Card.Id)
			} else {
				fmt.Printf("[*]")
			}
		}
		println()
	}
}
