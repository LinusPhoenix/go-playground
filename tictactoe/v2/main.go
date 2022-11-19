package main

import (
	"errors"
	"fmt"
	"linusphoenix/tictactoe/v2/game"
	"linusphoenix/tictactoe/v2/player"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Goals of this version:
// - Add a computer player that uses a naive strategy

func init() {
	rand.Seed(time.Now().UnixNano())
}

type model struct {
	game    *game.Game
	ai      player.Player
	cursorX int
	cursorY int
	err     error
}

func initialModel() model {
	// Determine which player the AI plays as.
	i := rand.Intn(2)
	var p game.CellState
	if i == 0 {
		p = game.X
	} else {
		p = game.O
	}
	ai, _ := player.New(p)
	tictactoe := game.New()
	// If ai is X, then it makes the first move.
	if ai.GetPlayer() == game.X {
		ai.MakeTurn(tictactoe)
	}
	cursorX, cursorY := 0, 0
	return model{tictactoe, ai, cursorX, cursorY, nil}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up":
			if m.cursorX > 0 {
				m.cursorX--
			}
		case "down":
			if m.cursorX < 2 {
				m.cursorX++
			}

		case "left":
			if m.cursorY > 0 {
				m.cursorY--
			}
		case "right":
			if m.cursorY < 2 {
				m.cursorY++
			}

		// The "enter" key makes the player take their turn.
		case "enter":
			if !m.game.IsGameOver() {
				// Process the player's turn.
				m.err = m.game.MakeTurn(m.cursorX, m.cursorY)
				if m.game.IsGameOver() {
					return m, tea.Quit
				}
				// Process the AI's turn.
				err := m.ai.MakeTurn(m.game)
				if err != nil {
					m.err = errors.New("The AI crashed! The game is over")
					return m, tea.Quit
				}
				if m.game.IsGameOver() {
					return m, tea.Quit
				}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	var humanp game.CellState
	if m.ai.GetPlayer() == game.X {
		humanp = game.O
	} else {
		humanp = game.X
	}
	s := fmt.Sprintf("You're playing as: %v\n", humanp)
	s += "Use the arrow keys to navigate on the board.\n"
	s += "Use enter to make your move.\n\n"

	playerActive := m.game.GetPlayerActive()
	gameOver := m.game.IsGameOver()
	s += fmt.Sprintf("Active player: %v\n\n", playerActive)

	// Render the board row by row.
	for i := 0; i <= 2; i++ {
		// Render each cell.
		for j := 0; j <= 2; j++ {
			if j != 0 {
				s += " "
			}
			cellState, _ := m.game.GetCellState(i, j)
			// Color the cell value if it's the most recent move.
			var coloredState string
			if move := m.game.GetLastMove(); move != nil && move.X == i && move.Y == j {
				coloredState = fmt.Sprintf("\033[4;32m%v\033[0m", cellState)
			} else {
				coloredState = fmt.Sprintf("%v", cellState)
			}
			s += fmt.Sprintf("[%v]", coloredState)
		}
		s += "\n"
		// If this is the cursor's row, render the cursor. Otherwise, render an empty row.
		if !gameOver && i == m.cursorX {
			// The cursor is either the 2nd, 5th, or 10th position in the row.
			// So we can multiply cursorY by 4 and append that many plus one to spaces to the string.
			s += strings.Repeat(" ", m.cursorY*4+1) + "^\n"
		} else {
			s += "\n"
		}
	}

	// Render the current error, if any, from the game.
	if m.err != nil {
		s += fmt.Sprintf("%v.\n", m.err)
	} else {
		s += "\n"
	}

	// Announce the winner. Otherwise, explain how to quit.
	if gameOver {
		winner := m.game.GetWinner()
		if winner == game.None {
			s += "The game is a draw!\n"
		} else {
			s += fmt.Sprintf("%v wins!\n", winner)
		}
	} else {
		s += "\nPress q to quit.\n"
	}

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
