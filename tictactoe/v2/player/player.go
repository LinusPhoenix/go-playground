package player

import (
	"errors"
	"linusphoenix/tictactoe/v2/game"
	"math/rand"
	"time"
)

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

	// Collect all possible moves.
	var movesOpen []move
	for i := 0; i <= 2; i++ {
		for j := 0; j <= 2; j++ {
			if cell, err := tictactoe.GetCellState(i, j); err == nil && cell == game.None {
				movesOpen = append(movesOpen, move{i, j})
			}
		}
	}

	// Pick one at random.
	move := movesOpen[rand.Intn(len(movesOpen))]
	err := tictactoe.MakeTurn(move.X, move.Y)
	if err != nil {
		return err
	}

	return nil
}
