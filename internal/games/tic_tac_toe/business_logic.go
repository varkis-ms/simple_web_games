package tic_tac_toe

import (
	"math/rand"
	"simple_web_games/internal/apperror"
	"simple_web_games/internal/games"
	"time"
)

type TttGameField struct {
	CheckWin  bool
	Status    bool
	Last      int
	MovesLeft int
	Size      int
	Field     [][]int
}

// NewGameField конструктор интерфейса
func NewGameField(size int) games.GameField {
	field := make([][]int, size)
	for i := range field {
		field[i] = make([]int, size)
	}
	rand.Seed(time.Now().UnixNano())
	return &TttGameField{
		CheckWin:  false,
		Status:    true,
		Last:      rand.Intn(2) + 1,
		MovesLeft: size * size,
		Size:      size,
		Field:     field,
	}
}

// Progress Валидация данных и возврат результата хода игрока
func (f *TttGameField) Progress(row, col, user int) (int, bool, error) {
	if f.Status == false {
		return 0, false, apperror.ErrNoGame
	}
	if f.Last == user {
		return 0, false, apperror.ErrWrongUser
	}
	field := f.Field
	if len(field) <= row || len(field) <= col {
		return 0, false, apperror.ErrIncorrectMove
	}
	if field[row][col] != 0 {
		return 0, false, apperror.ErrFilledCell
	}
	field[row][col] = user
	f.Last = user
	f.MovesLeft--
	if !f.CheckWin {
		if f.Size*f.Size-f.MovesLeft >= f.Size {
			f.CheckWin = true
		} else {
			return 0, false, nil
		}
	}
	winner := f.checkProgress(row, col)
	if winner == 0 {
		if f.MovesLeft == 0 {
			f.Status = false
			return 0, true, nil
		}
		return 0, false, nil
	}
	f.Status = false
	return winner, true, nil
}

// checkProgress функцция для проверки окончания игры и определения победиты/ничьи по последнему ходу игрока
func (f *TttGameField) checkProgress(row, col int) int {
	field := f.Field
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

// GetPrettyField Возвращает игровое поле в виде string с отформатированными значениями
func (f *TttGameField) GetPrettyField() string {
	prettyField := ""
	for _, i := range f.Field {
		for _, j := range i {
			sym := ""
			switch j {
			case 1:
				sym = "[O]"
			case 2:
				sym = "[X]"
			default:
				sym = "[ ]"
			}
			prettyField += sym
		}
		prettyField += "|"
	}
	return prettyField
}
