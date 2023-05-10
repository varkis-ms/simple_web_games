package tic_tac_toe_test

import (
	"math/rand"
	"simple_web_games/internal/apperror"
	"simple_web_games/internal/games"
	"simple_web_games/internal/games/tic_tac_toe"
	"testing"
	"time"
)

func NewGameFieldTest(size int) games.GameField {
	field := make([][]int, size)
	for i := range field {
		field[i] = make([]int, size)
	}
	rand.Seed(time.Now().UnixNano())
	return &tic_tac_toe.TttGameField{
		CheckWin:  false,
		Status:    true,
		Last:      1,
		MovesLeft: size * size,
		Size:      size,
		Field:     field,
	}
}

func TestTttGameField_Progress(t *testing.T) {
	// successful move
	field := NewGameFieldTest(2)
	_, check, err := field.Progress(0, 0, 2)
	if err != nil {
		t.Errorf("Incorrect result. Expect: nil, got: %s", err.Error())
	}
	if check != false {
		t.Errorf("Incorrect result. Expect: false, got: %v", check)
	}

	// user made a move out of turn
	_, check, err = field.Progress(0, 0, 2)
	if err != apperror.ErrWrongUser {
		t.Errorf("Incorrect result. Expect: %v, got: %v", apperror.ErrWrongUser.Error(), err.Error())
	}
	if check != false {
		t.Errorf("Incorrect result. Expect: false, got: %v", check)
	}

	// user made a move in the occupied cell
	_, check, err = field.Progress(0, 0, 1)
	if err != apperror.ErrFilledCell {
		t.Errorf("Incorrect result. Expect: %v, got: %v", apperror.ErrFilledCell.Error(), err.Error())
	}
	if check != false {
		t.Errorf("Incorrect result. Expect: false, got: %v", check)
	}

	// user has gone out of bounds
	_, check, err = field.Progress(10, 10, 1)
	if err != apperror.ErrIncorrectMove {
		t.Errorf("Incorrect result. Expect: %v, got: %v", apperror.ErrIncorrectMove.Error(), err.Error())
	}
	if check != false {
		t.Errorf("Incorrect result. Expect: false, got: %v", check)
	}

	// user has gone out of bounds
	_, check, err = field.Progress(10, 10, 1)
	if err != apperror.ErrIncorrectMove {
		t.Errorf("Incorrect result. Expect: %v, got: %v", apperror.ErrIncorrectMove.Error(), err.Error())
	}
	if check != false {
		t.Errorf("Incorrect result. Expect: false, got: %v", check)
	}

	// user win result
	//
	_, _, _ = field.Progress(1, 1, 1)
	winner, check, err := field.Progress(0, 1, 2)
	if err != nil {
		t.Errorf("Incorrect result. Expect: nil, got: %v", err.Error())
	}
	if check != true {
		t.Errorf("Incorrect result. Expect: false, got: %v", check)
	}
	if winner != 2 {
		t.Errorf("Incorrect result. Expect: 2, got: %v", winner)
	}

	// turn to the field after the end of the game (first user)
	_, _, err = field.Progress(0, 1, 2)
	if err != apperror.ErrNoGame {
		t.Errorf("Incorrect result. Expect: %v, got: %v", apperror.ErrNoGame.Error(), err.Error())
	}
	// second user
	_, _, err = field.Progress(0, 1, 1)
	if err != apperror.ErrNoGame {
		t.Errorf("Incorrect result. Expect: %v, got: %v", apperror.ErrNoGame.Error(), err.Error())
	}

	// draw situation
	field = NewGameFieldTest(3)
	_, _, _ = field.Progress(0, 1, 2)
	_, _, _ = field.Progress(0, 0, 1)
	_, _, _ = field.Progress(0, 2, 2)
	_, _, _ = field.Progress(1, 1, 1)
	_, _, _ = field.Progress(1, 0, 2)
	_, _, _ = field.Progress(1, 2, 1)
	_, _, _ = field.Progress(2, 0, 2)
	_, _, _ = field.Progress(2, 1, 1)
	winner, check, err = field.Progress(2, 2, 2)
	if err != nil {
		t.Errorf("Incorrect result. Expect: nil, got: %s", err.Error())
	}
	if check != true {
		t.Errorf("Incorrect result. Expect: nil, got: %s", err.Error())
	}
	if winner != 0 {
		t.Errorf("Incorrect result. Expect: 0, got: %d", winner)
	}
}
