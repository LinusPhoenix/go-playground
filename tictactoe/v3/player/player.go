package player

import (
	"errors"
	"linusphoenix/tictactoe/v3/game"
	"math/rand"
	"time"
)

var winDirections = [][3]move{
	// Rows.
	[3]move{move{0, 0}, move{0, 1}, move{0, 2}},
	[3]move{move{1, 0}, move{1, 1}, move{1, 2}},
	[3]move{move{2, 0}, move{2, 1}, move{2, 2}},
	// Columns
	[3]move{move{0, 0}, move{1, 0}, move{2, 0}},
	[3]move{move{0, 1}, move{1, 1}, move{2, 1}},
	[3]move{move{0, 2}, move{1, 2}, move{2, 2}},
	// Diagonals
	[3]move{move{0, 0}, move{1, 1}, move{2, 2}},
	[3]move{move{2, 0}, move{1, 1}, move{0, 2}},
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Player interface {
	MakeTurn(game *game.Game) error
	GetPlayer() game.CellState
}

type aiPlayer struct {
	player game.CellState
}

type move struct {
	X int
	Y int
}

func New(player game.CellState) (Player, error) {
	if player == game.None {
		return nil, errors.New("Player must play as X or O")
	}
	return &aiPlayer{player}, nil
}

func (ai *aiPlayer) GetPlayer() game.CellState {
	return ai.player
}

func (ai *aiPlayer) MakeTurn(tictactoe *game.Game) error {
	if tictactoe.GetPlayerActive() != ai.player {
		return errors.New("It's not the AI player's turn")
	}

	move := ai.calcMove(tictactoe)
	err := tictactoe.MakeTurn(move.X, move.Y)
	if err != nil {
		return err
	}
	return nil
}

func (ai *aiPlayer) calcMove(tictactoe *game.Game) move {
	var opp game.CellState
	if ai.player == game.X {
		opp = game.O
	} else {
		opp = game.X
	}

	// If there's a winning move on the board, play it.
	win, ok := findWinningMove(ai.player, tictactoe)
	if ok {
		return win
	}

	// If the opponent has a winning move, prevent it.
	loss, ok := findWinningMove(opp, tictactoe)
	if ok {
		return loss
	}

	// Collect all possible moves.
	var movesOpen []move
	for i := 0; i <= 2; i++ {
		for j := 0; j <= 2; j++ {
			if cell, _ := tictactoe.GetCellState(i, j); cell == game.None {
				movesOpen = append(movesOpen, move{i, j})
			}
		}
	}

	// Pick one at random.
	return movesOpen[rand.Intn(len(movesOpen))]
}

func findWinningMove(player game.CellState, tictactoe *game.Game) (move, bool) {
	for _, dir := range winDirections {
		move, ok := isAlmostWin(player, tictactoe, dir)
		if ok {
			return move, true
		}
	}
	return move{0, 0}, false
}

func isAlmostWin(player game.CellState, tictactoe *game.Game, winDirection [3]move) (move, bool) {
	// This direction is an almost win if the player has taken exactly two of the three cells, and the last one is taken by None.
	playerCells := 0
	winMoveFound := false
	var winMove move
	for _, move := range winDirection {
		cell, _ := tictactoe.GetCellState(move.X, move.Y)
		if cell == player {
			playerCells++
		}
		if cell == game.None {
			winMoveFound = true
			winMove = move
		}
	}
	if playerCells == 2 && winMoveFound {
		return winMove, true
	} else {
		return move{0, 0}, false
	}
}
