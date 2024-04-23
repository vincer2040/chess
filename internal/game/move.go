package game

import (
	"fmt"
	"math"
)

type LegalMoves map[int][]int
type AttackingMoves map[int][][]int

type Direction int

type Check struct {
	from int
	to   []int
}

type Checks struct {
	inCheck bool
	kingIdx int
	checks  []Check
}

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

func getLegalMoves(board Board, toMove byte, castleRights *CastleRights, enPassant int, attackingMoves AttackingMoves) LegalMoves {
	var legalMoves LegalMoves = make(LegalMoves)
	checks := getChecks(board, toMove, attackingMoves)
	fmt.Printf("attacking: %+v\n", attackingMoves)
	fmt.Printf("checks: %+v\n", checks)
	if checks.inCheck && len(checks.checks) > 1 {
		// we are in double check and can only move the king
		var color Piece
		if toMove == 'w' {
			color = White
		} else {
			color = Black
		}
		moves := getKingMoves(board, checks.kingIdx, color, castleRights, &checks)
		if len(moves) == 0 {
			return legalMoves
		}
		legalMoves[checks.kingIdx] = moves
		return legalMoves
	}
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
			moves := getPawnMoves(board, idx, color, enPassant, &checks)
			if len(moves) == 0 {
				break
			}
			legalMoves[idx] = moves
			break
		case Knight:
			moves := getKnightMoves(board, idx, color, &checks)
			if len(moves) == 0 {
				break
			}
			legalMoves[idx] = moves
			break
		case Bishop:
			moves := getDiagonalMoves(board, idx, color, &checks)
			if len(moves) == 0 {
				break
			}
			legalMoves[idx] = moves
			break
		case Rook:
			moves := getStraightMoves(board, idx, color, &checks)
			if len(moves) == 0 {
				break
			}
			legalMoves[idx] = moves
			break
		case Queen:
			moves := getStraightMoves(board, idx, color, &checks)
			moves = append(moves, getDiagonalMoves(board, idx, color, &checks)...)
			if len(moves) == 0 {
				break
			}
			legalMoves[idx] = moves
			break
		case King:
			moves := getKingMoves(board, idx, color, castleRights, &checks)
			if len(moves) == 0 {
				break
			}
			legalMoves[idx] = moves
			break
		}
	}
	return legalMoves
}

func getLegalMovesOtherSide(board Board, toMove byte) LegalMoves {
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
			moves := getPawnMoves(board, idx, color, -1, &Checks{inCheck: false})
			if len(moves) == 0 {
				break
			}
			legalMoves[idx] = moves
			break
		case Knight:
			moves := getKnightMoves(board, idx, color, &Checks{inCheck: false})
			if len(moves) == 0 {
				break
			}
			legalMoves[idx] = moves
			break
		case Bishop:
			moves := getDiagonalMoves(board, idx, color, &Checks{inCheck: false})
			if len(moves) == 0 {
				break
			}
			legalMoves[idx] = moves
			break
		case Rook:
			moves := getStraightMoves(board, idx, color, &Checks{inCheck: false})
			if len(moves) == 0 {
				break
			}
			legalMoves[idx] = moves
			break
		case Queen:
			moves := getStraightMoves(board, idx, color, &Checks{inCheck: false})
			moves = append(moves, getDiagonalMoves(board, idx, color, &Checks{inCheck: false})...)
			if len(moves) == 0 {
				break
			}
			legalMoves[idx] = moves
			break
		}
	}
	return legalMoves

}

func getAttackingMoves(board Board, toMove byte) AttackingMoves {
	attackingMoves := make(AttackingMoves)
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
		}
	}
	return attackingMoves
}

func getPawnMoves(board Board, idx int, color Piece, enPassant int, checks *Checks) []int {
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
		if checks.inCheck {
			if moveResolvesCheck(sq, checks) {
				res = append(res, sq)
			}
		} else {
			res = append(res, sq)
		}
	}
	rank := getRankForIdx(sq)
	left := sq - 1
	right := sq + 1
	if onStartSquare {
		sq += (8 * sign)
		if !board.hasPieceOnIdx(sq) {
			if checks.inCheck {
				if moveResolvesCheck(sq, checks) {
					res = append(res, sq)
				}
			} else {
				res = append(res, sq)
			}
		}
	}
	leftRank := getRankForIdx(left)
	rightRank := getRankForIdx(right)
	if board.hasPieceOnIdx(left) && !board.hasColorPieceOnIdx(left, color) && rank == leftRank {
		if checks.inCheck {
			if moveResolvesCheck(left, checks) {
				res = append(res, left)
			}
		} else {
			res = append(res, left)
		}
	}
	if board.hasPieceOnIdx(right) && !board.hasColorPieceOnIdx(right, color) && rank == rightRank {
		if checks.inCheck {
			if moveResolvesCheck(right, checks) {
				res = append(res, right)
			}
		} else {
			res = append(res, right)
		}
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
				if checks.inCheck {
					if moveResolvesCheck(left-8, checks) {
						res = append(res, left-8)
					}
				} else {
					res = append(res, left-8)
				}
			} else {
				if checks.inCheck {
					if moveResolvesCheck(left+8, checks) {
						res = append(res, left+8)
					}
				} else {
					res = append(res, left+8)
				}
			}
		}
	} else if enPassant == right {
		if rightRank == rank {
			if color == White {
				if checks.inCheck {
					if moveResolvesCheck(left-8, checks) {
						res = append(res, left-8)
					}
				} else {
					res = append(res, left-8)
				}
			} else {
				if checks.inCheck {
					if moveResolvesCheck(left+8, checks) {
						res = append(res, left+8)
					}
				} else {
					res = append(res, left+8)
				}
			}
		}
	}
	return res
}

func getDiagonalMoves(board Board, idx int, color Piece, checks *Checks) []int {
	var res []int
	offsets := []int{9, 7, -9, -7}
	directions := []Direction{SouthEast, SouthWest, NorthWest, NorthEast}
	for i, offset := range offsets {
		dir := directions[i]
		n := getMaxToEdge(idx, dir)
		sq := idx + offset
		for j := 0; j < n; j++ {
			if !board.hasPieceOnIdx(sq) {
				if checks.inCheck {
					if moveResolvesCheck(sq, checks) {
						res = append(res, sq)
					}
				} else {
					res = append(res, sq)
				}
				sq += offset
				continue
			}
			if !board.hasColorPieceOnIdx(sq, color) {
				if checks.inCheck {
					if moveResolvesCheck(sq, checks) {
						res = append(res, sq)
					}
				} else {
					res = append(res, sq)
				}
				break
			}
			break
		}
		sq = idx
	}
	return res
}

func getStraightMoves(board Board, idx int, color Piece, checks *Checks) []int {
	var res []int
	offsets := []int{8, 1, -8, -1}
	dirs := []Direction{South, East, North, West}
	for i, offset := range offsets {
		dir := dirs[i]
		n := getMaxToEdge(idx, dir)
		sq := idx + offset
		for j := 0; j < n; j++ {
			if !board.hasPieceOnIdx(sq) {
				if checks.inCheck {
					if moveResolvesCheck(sq, checks) {
						res = append(res, sq)
					}
				} else {
					res = append(res, sq)
				}
				sq += offset
				continue
			}
			if !board.hasColorPieceOnIdx(sq, color) {
				if checks.inCheck {
					if moveResolvesCheck(sq, checks) {
						res = append(res, sq)
					}
				} else {
					res = append(res, sq)
				}
				break
			}
			break
		}
		sq = idx
	}
	return res
}

func getKnightMoves(board Board, idx int, color Piece, checks *Checks) []int {
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
			if checks.inCheck {
				if moveResolvesCheck(sq, checks) {
					res = append(res, sq)
				}
			} else {
				res = append(res, sq)
			}
			continue
		}
		if !board.hasColorPieceOnIdx(sq, color) {
			if checks.inCheck {
				if moveResolvesCheck(sq, checks) {
					res = append(res, sq)
				}
			} else {
				res = append(res, sq)
			}
			continue
		}
	}
	return res
}

func getKingMoves(board Board, idx int, color Piece, castleRights *CastleRights, checks *Checks) []int {
	var res []int
	curRank := getRankForIdx(idx)
	minIdx := getMinIdxForRank(curRank)
	maxIdx := getMaxIdxForRank(curRank)
	straightOffsets := []int{8, 1, -1, -8}
	diagOffsets := []int{7, 9, -7, -9}
	var nextToMove byte
	if color == White {
		nextToMove = 'b'
	} else {
		nextToMove = 'w'
	}
	boardCopy := board.copy()
	for _, offset := range straightOffsets {
		sq := idx + offset
		if sq >= 64 || sq < 0 {
			continue
		}
		if int(math.Abs(float64(offset))) == 1 && (sq < minIdx || sq > maxIdx) {
			continue
		}
		if !board.hasPieceOnIdx(sq) {
			boardCopy[idx] = None
			boardCopy[sq] = King | color
			newLegalMoves := getLegalMovesOtherSide(boardCopy, nextToMove)
			if !legalMovesContainsCaptureOfIdx(sq, newLegalMoves) {
				res = append(res, sq)
			}
			boardCopy[idx] = King | color
			boardCopy[sq] = None
			continue
		}
		if !board.hasColorPieceOnIdx(sq, color) {
			old := boardCopy[sq]
			boardCopy[idx] = None
			boardCopy[sq] = King | color
			newLegalMoves := getLegalMovesOtherSide(boardCopy, nextToMove)
			if !legalMovesContainsCaptureOfIdx(sq, newLegalMoves) {
				res = append(res, sq)
			}
			boardCopy[idx] = King | color
			boardCopy[sq] = old
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
			boardCopy[idx] = None
			boardCopy[sq] = King | color
			newLegalMoves := getLegalMovesOtherSide(boardCopy, nextToMove)
			if !legalMovesContainsCaptureOfIdx(sq, newLegalMoves) {
				res = append(res, sq)
			}
			boardCopy[idx] = King | color
			boardCopy[sq] = None
			continue
		}
		if !board.hasColorPieceOnIdx(sq, color) {
			old := boardCopy[sq]
			boardCopy[idx] = None
			boardCopy[sq] = King | color
			newLegalMoves := getLegalMovesOtherSide(boardCopy, nextToMove)
			if !legalMovesContainsCaptureOfIdx(sq, newLegalMoves) {
				res = append(res, sq)
			}
			boardCopy[idx] = King | color
			boardCopy[sq] = old
			continue
		}
	}

	if !checks.inCheck {
		if color == White {
			if castleRights.WhiteKing && board[61] == None && board[62] == None {
				boardCopy[idx] = None
				boardCopy[62] = King | color
				newLegalMoves := getLegalMovesOtherSide(boardCopy, nextToMove)
				if !legalMovesContainsCaptureOfIdx(62, newLegalMoves) && !legalMovesContainsCaptureOfIdx(61, newLegalMoves) {
					res = append(res, 62)
				}
				boardCopy[idx] = King | color
				boardCopy[62] = None
			}
			if castleRights.WhiteQueen && board[59] == None && board[58] == None && board[57] == None {
				boardCopy[idx] = None
				boardCopy[58] = King | color
				newLegalMoves := getLegalMovesOtherSide(boardCopy, nextToMove)
				if !legalMovesContainsCaptureOfIdx(58, newLegalMoves) && !legalMovesContainsCaptureOfIdx(59, newLegalMoves) && !legalMovesContainsCaptureOfIdx(57, newLegalMoves) {
					res = append(res, 58)
				}
				boardCopy[idx] = King | color
				boardCopy[58] = None
			}
		} else {
			if castleRights.BlackKing && board[5] == None && board[6] == None {
				boardCopy[idx] = None
				boardCopy[6] = King | color
				newLegalMoves := getLegalMovesOtherSide(boardCopy, nextToMove)
				if !legalMovesContainsCaptureOfIdx(6, newLegalMoves) && !legalMovesContainsCaptureOfIdx(5, newLegalMoves) {
					res = append(res, 6)
				}
				boardCopy[idx] = King | color
				boardCopy[6] = None
			}
			if castleRights.BlackQueen && board[3] == None && board[2] == None && board[1] == None {
				boardCopy[idx] = None
				boardCopy[2] = King | color
				newLegalMoves := getLegalMovesOtherSide(boardCopy, nextToMove)
				if !legalMovesContainsCaptureOfIdx(2, newLegalMoves) && !legalMovesContainsCaptureOfIdx(2, newLegalMoves) && !legalMovesContainsCaptureOfIdx(1, newLegalMoves) {
					res = append(res, 2)
				}
				boardCopy[idx] = King | color
				boardCopy[2] = None
			}
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

func getChecks(board Board, toMove byte, attackingMoves AttackingMoves) Checks {
	var res Checks
	var king int
	for idx, pieceInfo := range board {
		color := pieceInfo & COLORMASK
		piece := pieceInfo & PIECEMASK
		if color == White && toMove == 'b' {
			continue
		}
		if color == Black && toMove == 'w' {
			continue
		}
		if piece == King {
			king = idx
			break
		}
	}
	res.kingIdx = king
	for from, dirs := range attackingMoves {
		for _, moves := range dirs {
			last := moves[len(moves)-1]
			if last == king {
				res.inCheck = true
				check := Check{
					from: from,
					to:   moves,
				}
				res.checks = append(res.checks, check)
				break
			}
		}
	}
	return res
}

func getAttackingPawnMoves(board Board, idx int, color Piece) [][]int {
	var res [][]int
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
		res = append(res, []int{left})
	}
	if board.hasPieceOnIdx(right) && !board.hasColorPieceOnIdx(right, color) && rank == rightRank {
		res = append(res, []int{right})
	}
	return res
}

func getAttackingKnightMoves(board Board, idx int, color Piece) [][]int {
	var res [][]int
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
			res = append(res, []int{sq})
			continue
		}
	}
	return res
}

func getAttackingDiagonalMoves(board Board, idx int, color Piece) [][]int {
	var res [][]int
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
		res = append(res, tmp)
		sq = idx
	}
	return res
}

func getAttackingStraightMoves(board Board, idx int, color Piece) [][]int {
	var res [][]int
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
		res = append(res, tmp)
		sq = idx
	}
	return res
}

func legalMovesContainsCaptureOfIdx(idx int, legalMoves LegalMoves) bool {
	for _, moves := range legalMoves {
		for _, move := range moves {
			if move == idx {
				return true
			}
		}
	}
	return false
}

func moveResolvesCheck(sq int, checks *Checks) bool {
	if len(checks.checks) != 1 {
		panic("not a single check")
	}
	check := checks.checks[0]
	if sq == check.from {
		return true
	}
	for _, i := range check.to {
		if i == sq {
			return true
		}
	}
	return false
}
