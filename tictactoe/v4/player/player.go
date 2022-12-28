package player

import (
	"errors"
	"linusphoenix/tictactoe/v4/game"
	"math"
)

type Player interface {
	MakeTurn(game *game.Game) error
	GetPlayer() game.CellState
}

type aiPlayer struct {
	player game.CellState
}

type move struct {
	X     int
	Y     int
	Score int
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

	move, err := ai.calcMove(tictactoe)
	if err != nil {
		return err
	}
	err = tictactoe.MakeTurn(move.X, move.Y)
	if err != nil {
		return err
	}
	return nil
}

func (ai *aiPlayer) calcMove(tictactoe *game.Game) (*move, error) {
	var moves []move
	for i := 0; i <= 2; i++ {
		for j := 0; j <= 2; j++ {
			if cell, _ := tictactoe.GetCellState(i, j); cell == game.None {
				nextGame := tictactoe.Copy()
				err := nextGame.MakeTurn(i, j)
				if err != nil {
					return nil, err
				}
				score, err := calcScore(nextGame)
				if err != nil {
					return nil, err
				}
				moves = append(moves, move{i, j, score})
			}
		}
	}

	if ai.player == game.X {
		best := move{-1, -1, math.MaxInt}
		for _, move := range moves {
			if move.Score <= best.Score {
				best = move
			}
		}
		return &best, nil
	} else {
		best := move{-1, -1, math.MinInt}
		for _, move := range moves {
			if move.Score >= best.Score {
				best = move
			}
		}
		return &best, nil
	}
}

func calcScore(tictactoe *game.Game) (int, error) {
	if tictactoe.IsGameOver() {
		switch tictactoe.GetWinner() {
		case game.X:
			return -1, nil
		case game.O:
			return 1, nil
		case game.None:
			return 0, nil
		}
	}

	sum := 0
	for i := 0; i <= 2; i++ {
		for j := 0; j <= 2; j++ {
			if cell, _ := tictactoe.GetCellState(i, j); cell == game.None {
				nextGame := tictactoe.Copy()
				err := nextGame.MakeTurn(i, j)
				if err != nil {
					return 0, err
				}
				score, err := calcScore(nextGame)
				if err != nil {
					return 0, err
				}
				sum += score
			}
		}
	}

	return sum, nil
}
