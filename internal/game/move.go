package game

import (
	"math"
)

type LegalMoves map[int][]int

type Direction int

const (
	North Direction = iota
	East
	South
	West
	NorthWest
	NorthEast
	SouthWest
	SouthEast
)

func getLegalMoves(board Board, toMove byte, castleRights *CastleRights, enPassant int) LegalMoves {
	var legalMoves LegalMoves = make(LegalMoves)
	for idx, pieceInfo := range board {
		color := pieceInfo & COLORMASK
		piece := pieceInfo & PIECEMASK
		if toMove == 'w' && color == Black {
			continue
		}
		if toMove == 'b' && color == White {
			continue
		}

		switch piece {
		case Pawn:
			moves := getPawnMoves(board, idx, color, enPassant)
			if len(moves) == 0 {
				break
			}
			legalMoves[idx] = moves
			break
		case Knight:
			moves := getKnightMoves(board, idx, color)
			if len(moves) == 0 {
				break
			}
			legalMoves[idx] = moves
			break
		case Bishop:
			moves := getDiagonalMoves(board, idx, color)
			if len(moves) == 0 {
				break
			}
			legalMoves[idx] = moves
			break
		case Rook:
			moves := getStraightMoves(board, idx, color)
			if len(moves) == 0 {
				break
			}
			legalMoves[idx] = moves
			break
		case Queen:
			moves := getStraightMoves(board, idx, color)
			moves = append(moves, getDiagonalMoves(board, idx, color)...)
			if len(moves) == 0 {
				break
			}
			legalMoves[idx] = moves
			break
		case King:
			moves := getKingMoves(board, idx, color, castleRights)
			if len(moves) == 0 {
				break
			}
			legalMoves[idx] = moves
			break
		}
	}
	return legalMoves
}

func getAttackingMoves(board Board, toMove byte) LegalMoves {
	attackingMoves := make(LegalMoves)
	for idx, pieceInfo := range board {
		piece := pieceInfo & PIECEMASK
		color := pieceInfo & COLORMASK
		if toMove == 'w' && color == White {
			continue
		}
		if toMove == 'b' && color == Black {
			continue
		}
		switch piece {
		case Pawn:
			moves := getAttackingPawnMoves(board, idx, color)
			if len(moves) == 0 {
				break
			}
			attackingMoves[idx] = moves
			break
		case Knight:
			moves := getAttackingKnightMoves(board, idx, color)
			if len(moves) == 0 {
				break
			}
			attackingMoves[idx] = moves
			break
		case Bishop:
			moves := getAttackingDiagonalMoves(board, idx, color)
			if len(moves) == 0 {
				break
			}
			attackingMoves[idx] = moves
			break
		case Rook:
			moves := getAttackingStraightMoves(board, idx, color)
			if len(moves) == 0 {
				break
			}
			attackingMoves[idx] = moves
			break
        case Queen:
            moves := getAttackingDiagonalMoves(board, idx, color)
            moves = append(moves, getAttackingStraightMoves(board, idx, color)...)
            if len(moves) == 0 {
                break
            }
            attackingMoves[idx] = moves
            break
        case King:
            moves := getAttackingKingMoves(board, idx, color)
            if len(moves) == 0 {
                break
            }
            attackingMoves[idx] = moves
            break
		}
	}
	return attackingMoves
}

func getAttackingPawnMoves(board Board, idx int, color Piece) []int {
	var res []int
	var sign int
	if color == White {
		sign = -1
	} else {
		sign = 1
	}
	sq := idx + (8 * sign)
	rank := getRankForIdx(sq)
	left := sq - 1
	right := sq + 1
	leftRank := getRankForIdx(left)
	rightRank := getRankForIdx(right)
	if board.hasPieceOnIdx(left) && !board.hasColorPieceOnIdx(left, color) && rank == leftRank {
		res = append(res, left)
	}
	if board.hasPieceOnIdx(right) && !board.hasColorPieceOnIdx(right, color) && rank == rightRank {
		res = append(res, right)
	}
	return res
}

func getAttackingKnightMoves(board Board, idx int, color Piece) []int {
	var res []int
	offsets := []struct {
		x int
		y int
	}{
		{x: 1, y: 16},
		{x: -1, y: 16},
		{x: 2, y: 8},
		{x: -2, y: 8},
		{x: 1, y: -16},
		{x: -1, y: -16},
		{x: 2, y: -8},
		{x: -2, y: -8},
	}
	for _, offset := range offsets {
		sq := idx + offset.y
		if sq >= 64 || sq < 0 {
			continue
		}
		rank := getRankForIdx(sq)
		sq += offset.x
		if sq < getMinIdxForRank(rank) || sq > getMaxIdxForRank(rank) {
			continue
		}
		if board.hasPieceOnIdx(sq) && !board.hasColorPieceOnIdx(sq, color) {
			res = append(res, sq)
			continue
		}
	}
	return res
}

func getAttackingDiagonalMoves(board Board, idx int, color Piece) []int {
	var res []int
	offsets := []int{9, 7, -9, -7}
	directions := []Direction{SouthEast, SouthWest, NorthWest, NorthEast}
	for i, offset := range offsets {
		dir := directions[i]
		n := getMaxToEdge(idx, dir)
		sq := idx + offset
		var tmp []int
		for j := 0; j < n; j++ {
			if !board.hasPieceOnIdx(sq) {
				tmp = append(tmp, sq)
				sq += offset
				continue
			}
			if !board.hasColorPieceOnIdx(sq, color) {
				tmp = append(tmp, sq)
				break
			}
			break
		}
		tmpLen := len(tmp)
		if tmpLen == 0 {
			sq = idx
			continue
		}
		lastChecked := tmp[tmpLen-1]
		if !board.hasPieceOnIdx(lastChecked) {
			sq = idx
			continue
		}
		res = append(res, tmp...)
		sq = idx
	}
	return res
}

func getAttackingStraightMoves(board Board, idx int, color Piece) []int {
	var res []int
	offsets := []int{8, 1, -8, -1}
	dirs := []Direction{South, East, North, West}
	for i, offset := range offsets {
		dir := dirs[i]
		n := getMaxToEdge(idx, dir)
		sq := idx + offset
		var tmp []int
		for j := 0; j < n; j++ {
			if !board.hasPieceOnIdx(sq) {
				tmp = append(tmp, sq)
				sq += offset
				continue
			}
			if !board.hasColorPieceOnIdx(sq, color) {
				tmp = append(tmp, sq)
				break
			}
			break
		}
		tmpLen := len(tmp)
		if tmpLen == 0 {
			sq = idx
			continue
		}
		lastChecked := tmp[tmpLen-1]
		if !board.hasPieceOnIdx(lastChecked) {
			sq = idx
			continue
		}
		res = append(res, tmp...)
		sq = idx
	}
	return res
}

func getAttackingKingMoves(board Board, idx int, color Piece) []int {
	var res []int
	curRank := getRankForIdx(idx)
	minIdx := getMinIdxForRank(curRank)
	maxIdx := getMaxIdxForRank(curRank)
	straightOffsets := []int{8, 1, -1, -8}
	diagOffsets := []int{7, 9, -7, -9}
    var tmp1 []int
	for _, offset := range straightOffsets {
		sq := idx + offset
		if sq >= 64 || sq < 0 {
			continue
		}
		if int(math.Abs(float64(offset))) == 1 && (sq < minIdx || sq > maxIdx) {
			continue
		}
		if board.hasPieceOnIdx(sq) && !board.hasColorPieceOnIdx(sq, color) {
			tmp1 = append(tmp1, sq)
			continue
		}
	}

    if len(tmp1) > 0 {
        res = append(res, tmp1...)
    }

    var tmp2 []int
	for _, offset := range diagOffsets {
		sq := idx + offset
		if sq > 64 || sq < 0 {
			continue
		}
		rank := getRankForIdx(sq)
		if rank != (curRank-1) && rank != (curRank+1) {
			continue
		}
		if board.hasPieceOnIdx(sq) && !board.hasColorPieceOnIdx(sq, color) {
			tmp2 = append(tmp2, sq)
			continue
		}
	}

    if len(tmp2) > 0 {
        res = append(res, tmp2...)
    }
    return res
}

func getPawnMoves(board Board, idx int, color Piece, enPassant int) []int {
	var res []int
	var sign int
	var onStartSquare bool
	if color == White {
		sign = -1
		onStartSquare = idx >= 48 && idx <= 55
	} else {
		sign = 1
		onStartSquare = idx >= 8 && idx <= 15
	}
	sq := idx + (8 * sign)
	if !board.hasPieceOnIdx(sq) {
		res = append(res, sq)
	}
	rank := getRankForIdx(sq)
	left := sq - 1
	right := sq + 1
	if onStartSquare {
		sq += (8 * sign)
		if !board.hasPieceOnIdx(sq) {
			res = append(res, sq)
		}
	}
	leftRank := getRankForIdx(left)
	rightRank := getRankForIdx(right)
	if board.hasPieceOnIdx(left) && !board.hasColorPieceOnIdx(left, color) && rank == leftRank {
		res = append(res, left)
	}
	if board.hasPieceOnIdx(right) && !board.hasColorPieceOnIdx(right, color) && rank == rightRank {
		res = append(res, right)
	}
	if enPassant == -1 {
		return res
	}
	rank = getRankForIdx(idx)
	left = idx - 1
	right = idx + 1
	leftRank = getRankForIdx(left)
	rightRank = getRankForIdx(right)
	if enPassant == left {
		if leftRank == rank {
			if color == White {
				res = append(res, left-8)
			} else {
				res = append(res, left+8)
			}
		}
	} else if enPassant == right {
		if rightRank == rank {
			if color == White {
				res = append(res, right-8)
			} else {
				res = append(res, right+8)
			}
		}
	}
	return res
}

func getDiagonalMoves(board Board, idx int, color Piece) []int {
	var res []int
	offsets := []int{9, 7, -9, -7}
	directions := []Direction{SouthEast, SouthWest, NorthWest, NorthEast}
	for i, offset := range offsets {
		dir := directions[i]
		n := getMaxToEdge(idx, dir)
		sq := idx + offset
		for j := 0; j < n; j++ {
			if !board.hasPieceOnIdx(sq) {
				res = append(res, sq)
				sq += offset
				continue
			}
			if !board.hasColorPieceOnIdx(sq, color) {
				res = append(res, sq)
				break
			}
			break
		}
		sq = idx
	}
	return res
}

func getStraightMoves(board Board, idx int, color Piece) []int {
	var res []int
	offsets := []int{8, 1, -8, -1}
	dirs := []Direction{South, East, North, West}
	for i, offset := range offsets {
		dir := dirs[i]
		n := getMaxToEdge(idx, dir)
		sq := idx + offset
		for j := 0; j < n; j++ {
			if !board.hasPieceOnIdx(sq) {
				res = append(res, sq)
				sq += offset
				continue
			}
			if !board.hasColorPieceOnIdx(sq, color) {
				res = append(res, sq)
				break
			}
			break
		}
		sq = idx
	}
	return res
}

func getKnightMoves(board Board, idx int, color Piece) []int {
	var res []int
	offsets := []struct {
		x int
		y int
	}{
		{x: 1, y: 16},
		{x: -1, y: 16},
		{x: 2, y: 8},
		{x: -2, y: 8},
		{x: 1, y: -16},
		{x: -1, y: -16},
		{x: 2, y: -8},
		{x: -2, y: -8},
	}
	for _, offset := range offsets {
		sq := idx + offset.y
		if sq >= 64 || sq < 0 {
			continue
		}
		rank := getRankForIdx(sq)
		sq += offset.x
		if sq < getMinIdxForRank(rank) || sq > getMaxIdxForRank(rank) {
			continue
		}
		if !board.hasPieceOnIdx(sq) {
			res = append(res, sq)
			continue
		}
		if !board.hasColorPieceOnIdx(sq, color) {
			res = append(res, sq)
			continue
		}
	}
	return res
}

func getKingMoves(board Board, idx int, color Piece, castleRights *CastleRights) []int {
	var res []int
	curRank := getRankForIdx(idx)
	minIdx := getMinIdxForRank(curRank)
	maxIdx := getMaxIdxForRank(curRank)
	straightOffsets := []int{8, 1, -1, -8}
	diagOffsets := []int{7, 9, -7, -9}
	for _, offset := range straightOffsets {
		sq := idx + offset
		if sq >= 64 || sq < 0 {
			continue
		}
		if int(math.Abs(float64(offset))) == 1 && (sq < minIdx || sq > maxIdx) {
			continue
		}
		if !board.hasPieceOnIdx(sq) {
			res = append(res, sq)
			continue
		}
		if !board.hasColorPieceOnIdx(sq, color) {
			res = append(res, sq)
			continue
		}
	}

	for _, offset := range diagOffsets {
		sq := idx + offset
		if sq > 64 || sq < 0 {
			continue
		}
		rank := getRankForIdx(sq)
		if rank != (curRank-1) && rank != (curRank+1) {
			continue
		}
		if !board.hasPieceOnIdx(sq) {
			res = append(res, sq)
			continue
		}
		if !board.hasColorPieceOnIdx(sq, color) {
			res = append(res, sq)
			continue
		}
	}

	if color == White {
		if castleRights.WhiteKing && board[61] == None && board[62] == None {
			res = append(res, 62)
		}
		if castleRights.WhiteQueen && board[59] == None && board[58] == None && board[57] == None {
			res = append(res, 58)
		}
	} else {
		if castleRights.BlackKing && board[5] == None && board[6] == None {
			res = append(res, 6)
		}
		if castleRights.BlackQueen && board[3] == None && board[2] == None && board[1] == None {
			res = append(res, 2)
		}
	}
	return res
}

func getMaxToEdge(idx int, dir Direction) int {
	rank := getRankForIdx(idx)
	file := getFileForIdx(idx)
	switch dir {
	case North:
		return rank
	case South:
		return 7 - rank
	case West:
		return file
	case East:
		return 7 - file
	case NorthWest:
		n := rank
		w := file
		return int(math.Min(float64(n), float64(w)))
	case NorthEast:
		n := rank
		e := 7 - file
		return int(math.Min(float64(n), float64(e)))
	case SouthWest:
		s := 7 - rank
		w := file
		return int(math.Min(float64(s), float64(w)))
	case SouthEast:
		s := 7 - rank
		e := 7 - file
		return int(math.Min(float64(s), float64(e)))
	}
	return 0
}
