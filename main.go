package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	Empty = ' '
	X     = 'X'
	O     = 'O'
)

type Board [3][3]rune

func (b *Board) Reset() {
	for i := range b {
		for j := range b[i] {
			b[i][j] = Empty
		}
	}
}

func (b *Board) String() string {
	var sb strings.Builder
	sb.WriteString("\n")
	for i := 0; i < 3; i++ {
		sb.WriteString(" ")
		for j := 0; j < 3; j++ {
			cell := b[i][j]
			if cell == Empty {
				// show numbers 1..9 for empty cells
				num := i*3 + j + 1
				sb.WriteString(fmt.Sprintf("%d", num))
			} else {
				sb.WriteString(string(cell))
			}
			if j < 2 {
				sb.WriteString(" | ")
			}
		}
		if i < 2 {
			sb.WriteString("\n-----------\n")
		}
	}
	sb.WriteString("\n\n")
	return sb.String()
}

func (b *Board) IsFull() bool {
	for i := range b {
		for j := range b[i] {
			if b[i][j] == Empty {
				return false
			}
		}
	}
	return true
}

func (b *Board) Winner() rune {
	// rows, cols
	for i := 0; i < 3; i++ {
		if b[i][0] != Empty && b[i][0] == b[i][1] && b[i][1] == b[i][2] {
			return b[i][0]
		}
		if b[0][i] != Empty && b[0][i] == b[1][i] && b[1][i] == b[2][i] {
			return b[0][i]
		}
	}
	// diagonals
	if b[0][0] != Empty && b[0][0] == b[1][1] && b[1][1] == b[2][2] {
		return b[0][0]
	}
	if b[0][2] != Empty && b[0][2] == b[1][1] && b[1][1] == b[2][0] {
		return b[0][2]
	}
	return Empty
}

func humanMove(b *Board, player rune, reader *bufio.Reader) {
	for {
		fmt.Printf("Player %c, enter cell (1-9): ", player)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		num, err := strconv.Atoi(text)
		if err != nil || num < 1 || num > 9 {
			fmt.Println("Invalid input. Type a number from 1 to 9 (shown on board). Try again.")
			continue
		}
		i := (num - 1) / 3
		j := (num - 1) % 3
		if b[i][j] != Empty {
			fmt.Println("Cell already taken. Pick another.")
			continue
		}
		b[i][j] = player
		break
	}
}

func availableMoves(b *Board) []int {
	moves := []int{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if b[i][j] == Empty {
				moves = append(moves, i*3+j)
			}
		}
	}
	return moves
}

// Minimax to make the AI unbeatable.
func minimax(b *Board, depth int, maximizing bool, ai, human rune) int {
	winner := b.Winner()
	if winner == ai {
		return 10 - depth
	}
	if winner == human {
		return depth - 10
	}
	if b.IsFull() {
		return 0
	}

	if maximizing {
		best := -1000
		for _, m := range availableMoves(b) {
			i := m / 3
			j := m % 3
			b[i][j] = ai
			score := minimax(b, depth+1, false, ai, human)
			b[i][j] = Empty
			if score > best {
				best = score
			}
		}
		return best
	} else {
		best := 1000
		for _, m := range availableMoves(b) {
			i := m / 3
			j := m % 3
			b[i][j] = human
			score := minimax(b, depth+1, true, ai, human)
			b[i][j] = Empty
			if score < best {
				best = score
			}
		}
		return best
	}
}

func aiBestMove(b *Board, ai, human rune) int {
	bestScore := -10000
	bestMove := -1
	moves := availableMoves(b)
	// if multiple equal moves, pick randomly among them for variety
	rand.Shuffle(len(moves), func(i, j int) { moves[i], moves[j] = moves[j], moves[i] })
	for _, m := range moves {
		i := m / 3
		j := m % 3
		b[i][j] = ai
		score := minimax(b, 0, false, ai, human)
		b[i][j] = Empty
		if score > bestScore {
			bestScore = score
			bestMove = m
		}
	}
	return bestMove
}

func aiMove(b *Board, player rune, opponent rune) {
	m := aiBestMove(b, player, opponent)
	if m == -1 {
		// fallback: pick first available
		moves := availableMoves(b)
		if len(moves) == 0 {
			return
		}
		m = moves[0]
	}
	i := m / 3
	j := m % 3
	b[i][j] = player
	fmt.Printf("Computer (%c) plays cell %d\n", player, m+1)
}

func chooseMode(reader *bufio.Reader) int {
	for {
		fmt.Println("Choose mode:")
		fmt.Println("1) Two players (local)")
		fmt.Println("2) Play vs Computer (you choose X or O)")
		fmt.Print("Enter 1 or 2: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "1" {
			return 1
		}
		if text == "2" {
			return 2
		}
		fmt.Println("Invalid option. Try again.")
	}
}

func choosePlayer(reader *bufio.Reader) rune {
	for {
		fmt.Print("Do you want to play as X or O? (X goes first): ")
		text, _ := reader.ReadString('\n')
		text = strings.ToUpper(strings.TrimSpace(text))
		if text == "X" {
			return X
		}
		if text == "O" {
			return O
		}
		fmt.Println("Invalid choice. Enter X or O.")
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)

	// Create board
	var board Board
	board.Reset()

	mode := chooseMode(reader)
	playerTurn := X // X always starts
	humanIs := X
	vsComputer := false

	if mode == 2 {
		vsComputer = true
		humanIs = choosePlayer(reader)
		fmt.Printf("You are %c. Let's play!\n", humanIs)
	} else {
		fmt.Println("Two-player mode. X goes first.")
	}

	for {
		fmt.Print(board.String())
		if board.Winner() != Empty || board.IsFull() {
			break
		}
		if vsComputer {
			if playerTurn == humanIs {
				humanMove(&board, playerTurn, reader)
			} else {
				aiMove(&board, playerTurn, humanIs)
			}
		} else {
			humanMove(&board, playerTurn, reader)
		}
		// swap
		if playerTurn == X {
			playerTurn = O
		} else {
			playerTurn = X
		}
	}

	fmt.Print(board.String())
	winner := board.Winner()
	if winner != Empty {
		fmt.Printf("Winner: %c\n", winner)
	} else {
		fmt.Println("It's a draw!")
	}
	fmt.Println("Thanks for playing.")
}
