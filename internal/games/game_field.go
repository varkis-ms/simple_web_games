package games

type GameField interface {
	Progress(row, col, user int) (int, bool, error)
}
