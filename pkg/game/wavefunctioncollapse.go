package game

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

func getInitialBuildBoard(width, height int) [][]buildCell {
	board := make([][]buildCell, width)
	for i := range board {
		board[i] = make([]buildCell, height)
		for j := range board[i] {
			board[i][j] = initialBuildCell
		}
	}
	return board
}

// loop through all of the buildCells in the board and set the entropy to the number of possible connectors
// we need to test the card's connectors against the the surrounding cells connectors
func getEntropyBoard(board [][]buildCell, cards []Card) [][][]int {

	entropyBoard := make([][][]int, len(board))

	for i, row := range board {

		entropyBoard[i] = make([][]int, len(row))

		for j, cell := range row {
			placed := cell.placed
			buildCell := buildCell{connectors: getEntropicCard(board, i, j)}

			if !placed {
				// compare the built cell to all of the cards
				for _, card := range cards {
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

func evolveBoard(r Rnd, buildBoard *[][]buildCell, cards []Card, board *[][]Tile) {

	// the whole loop

	entropyBoard := getEntropyBoard(*buildBoard, cards)
	availableCards := countEntropyBoard(entropyBoard)

	// select a random available
	selectedAvaialable := availableCards[r.Intn(len(availableCards))]

	// select a random id from the available
	selectedCardId := selectedAvaialable.ids[r.Intn(len(selectedAvaialable.ids))]
	// place the card in the buildBoard
	(*buildBoard)[selectedAvaialable.x][selectedAvaialable.y] = buildCell{placed: true, connectors: cards[selectedCardId].Connectors}

	// place the card on the board
	(*board)[selectedAvaialable.x][selectedAvaialable.y] = Tile{Card: cards[selectedCardId], X: selectedAvaialable.x, Y: selectedAvaialable.y}
}
