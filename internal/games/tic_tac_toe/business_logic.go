package tic_tac_toe

import (
	"math/rand"
	"simple_web_games/internal/apperror"
	"simple_web_games/internal/games"
	"time"
)

type TttGameField struct {
	checkWin  bool
	last      int
	movesLeft int
	size      int
	field     [][]int
}

func NewGameField(size *int) games.GameField {
	field := make([][]int, *size)
	for i := range field {
		field[i] = make([]int, *size)
	}
	rand.Seed(time.Now().UnixNano())
	return &TttGameField{
		field:     field,
		last:      rand.Intn(2) + 1,
		movesLeft: *size * *size,
		size:      *size,
		checkWin:  false,
	}
}

func (f *TttGameField) Progress(row, col, user int) (int, bool, error) {
	if f.last == user {
		return 0, false, apperror.ErrWrongUser
	}
	field := f.field
	if len(field) <= row || len(field) <= col {
		return 0, false, apperror.ErrIncorrectMove
	}
	if field[row][col] != 0 {
		return 0, false, apperror.ErrFilledCell
	}
	field[row][col] = user
	f.last = user
	f.movesLeft--
	if !f.checkWin {
		if f.size*f.size-f.movesLeft >= f.size {
			f.checkWin = true
		} else {
			return 0, false, nil
		}
	}
	winner := f.checkProgress(row, col)
	if winner == 0 {
		return 0, false, nil
	}
	return winner, true, nil
}

func (f *TttGameField) checkProgress(row, col int) int {
	field := f.field
	checkWin := field[row][0]
	for _, v := range field[row] {
		if v != checkWin {
			checkWin = 0
			break
		}
	}
	if checkWin != 0 {
		return checkWin
	}
	checkWin = field[0][col]
	for i := range field {
		if field[i][col] != checkWin {
			checkWin = 0
			break
		}
	}
	if checkWin != 0 {
		return checkWin
	}
	if row == col {
		checkWin = field[0][0]
		for i := range field {
			if field[i][i] != checkWin {
				checkWin = 0
				break
			}
		}
		if checkWin != 0 {
			return checkWin
		}
		checkWin = field[len(field)-1][len(field)-1]
		for i := range field {
			if field[i][len(field)-1-i] != checkWin {
				checkWin = 0
				break
			}
		}
		if checkWin != 0 {
			return checkWin
		}
	}
	return 0
}
