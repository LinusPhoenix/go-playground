package game

import (
	"errors"
	"fmt"
)

type CellState int

func (c CellState) String() string {
	switch c {
	case None:
		return " "
	case X:
		return "X"
	case O:
		return "O"
	}
	return "Unknown CellState!"
}

const (
	None CellState = iota
	X    CellState = iota
	O    CellState = iota
)

type Move struct {
	X      int
	Y      int
	player CellState
}

type Game struct {
	// The tictactoe board.
	// First dimension represents rows, second dimension represents columns.
	// [0][0] is the cell in the top left-hand corner,
	// [2][2] is the cell in the bottom right-hand corner.
	board [3][3]CellState
	// While the game is active, swaps between X and O. X goes first.
	// Once the game is over, this is set to None
	playerActive CellState
	moveLast     *Move
	// While the game is going on, this is false, and true when the game is finished.
	gameOver bool
	// While gameOver is false, this is set to None.
	// When gameOver is true, this is X or O if there is a winner, or None is the game is a draw.
	winner CellState
}

func New() *Game {
	var board [3][3]CellState
	for i, row := range board {
		for j := range row {
			board[i][j] = None
		}
	}

	return &Game{board, X, nil, false, None}
}

func (game *Game) Copy() *Game {
	newGame := New()
	for i, row := range newGame.board {
		for j := range row {
			newGame.board[i][j] = game.board[i][j]
		}
	}
	newGame.playerActive = game.playerActive
	moveLast := game.moveLast
	if moveLast != nil {
		newGame.moveLast = &Move{moveLast.X, moveLast.Y, moveLast.player}
	}
	newGame.gameOver = game.gameOver
	newGame.winner = game.winner

	return newGame
}

func (game *Game) GetPlayerActive() CellState {
	return game.playerActive
}

func (game *Game) GetCellState(x int, y int) (CellState, error) {
	if !areCoordsInBounds(x, y) {
		return None, errors.New("Invalid move. Each coordinate must be in the interval [0, 2]")
	}
	return game.board[x][y], nil
}

func (game *Game) IsGameOver() bool {
	return game.gameOver
}

func (game *Game) GetWinner() CellState {
	return game.winner
}

func (game *Game) GetLastMove() *Move {
	return game.moveLast
}

func (game *Game) MakeTurn(x int, y int) error {
	if game.gameOver {
		return errors.New("The game is over. No more moves can be made")
	}
	if !areCoordsInBounds(x, y) {
		return errors.New("Invalid move. Each coordinate must be in the interval [0, 2]")
	}

	if cell := game.board[x][y]; cell != None {
		return fmt.Errorf("Invalid move. The cell at [%d][%d] is already claimed by %#v", x, y, cell)
	}
	game.board[x][y] = game.playerActive
	game.moveLast = &Move{x, y, game.playerActive}

	switch game.playerActive {
	case X:
		game.playerActive = O
	case O:
		game.playerActive = X
	}

	isOver, winner := game.isGameOver()
	game.gameOver = isOver
	game.winner = winner

	return nil
}

func (game *Game) isGameOver() (bool, CellState) {
	// If there is no last move, there can be no winner.
	if game.moveLast == nil {
		return false, None
	}

	// Only the player who made the last move can win.
	winCandidate := game.moveLast.player
	x := game.moveLast.X
	y := game.moveLast.Y

	// Should the last move be made by None, there can be no winner.
	if winCandidate == None {
		return false, None
	}

	// We need to check the row, column, and diagonals for a potential win.
	row := game.checkRowForWinner(x, y, winCandidate)
	column := game.checkColumnForWinner(x, y, winCandidate)
	diagonal := game.checkDiagonalForWinner(x, y, winCandidate)
	antidiagonal := game.checkAntidiagonalForWinner(x, y, winCandidate)
	if row == winCandidate || column == winCandidate || diagonal == winCandidate || antidiagonal == winCandidate {
		return true, winCandidate
	}

	isDraw := true
	for i, row := range game.board {
		for j := range row {
			// If any cell is still None, then the game is not a draw.
			isDraw = isDraw && game.board[i][j] != None
		}
	}
	return isDraw, None
}

func (game *Game) checkRowForWinner(x int, y int, winCandidate CellState) CellState {
	if winCandidate == None {
		return None
	}
	for j := 0; j <= 2; j++ {
		if game.board[x][j] != winCandidate {
			return None
		}
	}

	return winCandidate
}

func (game *Game) checkColumnForWinner(x int, y int, winCandidate CellState) CellState {
	if winCandidate == None {
		return None
	}

	for i := 0; i <= 2; i++ {
		if game.board[i][y] != winCandidate {
			return None
		}
	}

	return winCandidate
}

func (game *Game) checkDiagonalForWinner(x int, y int, winCandidate CellState) CellState {
	if winCandidate == None {
		return None
	}
	// If x does not equal y, the move was not made on a diagonal.
	if x != y {
		return None
	}

	// Check top left to bottom right
	for i := 0; i <= 2; i++ {
		if game.board[i][i] != winCandidate {
			return None
		}
	}

	return winCandidate
}

func (game *Game) checkAntidiagonalForWinner(x int, y int, winCandidate CellState) CellState {
	if winCandidate == None {
		return None
	}

	// Check bottom left to top right
	for i := 0; i <= 2; i++ {
		if game.board[2-i][i] != winCandidate {
			return None
		}
	}

	return winCandidate
}

func areCoordsInBounds(x int, y int) bool {
	return x >= 0 || x <= 2 && y >= 0 || y <= 2
}
