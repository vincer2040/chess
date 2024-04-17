package game

import (
	"math"
)

type LegalMoves map[int][]int

func getLegalMoves(board Board, toMove byte, castleRights *CastleRights, enPessant int) LegalMoves {
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
			moves := getPawnMoves(board, idx, color)
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
			moves := GetKingMoves(board, idx, color, castleRights)
			if len(moves) == 0 {
				break
			}
			legalMoves[idx] = moves
		}
	}
	return legalMoves
}

func getPawnMoves(board Board, idx int, color Piece) []int {
	var res []int
	var sign int
	var onStartSquare bool
	if color == White {
		onStartSquare = idx >= 48 && idx <= 55
		sign = -1
	} else {
		onStartSquare = idx >= 8 && idx <= 15
		sign = 1
	}
	sq := idx + (sign * 8)
	if !board.hasPieceOnIdx(sq) && sq < 64 {
		res = append(res, sq)
	}
	rank := getRankForIdx(sq)
	maxIdxForTake := getMaxIdxForRank(rank)
	minIdxForTake := getMinIdxForRank(rank)
	left := sq - 1
	right := sq + 1
	if left >= minIdxForTake && board.hasPieceOnIdx(left) && !board.hasColorPieceOnIdx(left, color) {
		res = append(res, left)
	}
	if right <= maxIdxForTake && board.hasPieceOnIdx(left) && !board.hasColorPieceOnIdx(right, color) {
		res = append(res, right)
	}
	if onStartSquare {
		sq = idx + (sign * 16)
		res = append(res, sq)
	}
	return res
}

func getDiagonalMoves(board Board, idx int, color Piece) []int {
	var res []int
	offsets := []int{9, 7, -9, -7}
	sq := idx
	for _, offset := range offsets {
		sq += offset
		for sq < 64 && sq >= 0 {
			if !board.hasPieceOnIdx(sq) {
				res = append(res, sq)
				sq += offset
				continue
			}
			if board.hasColorPieceOnIdx(sq, color) {
				break
			}
			if squareIsEdge(sq) {
				if !board.hasPieceOnIdx(sq) {
					res = append(res, sq)
					sq += offset
				} else if !board.hasColorPieceOnIdx(sq, color) {
					res = append(res, sq)
				}
				break
			}
			res = append(res, sq)
			break
		}
		sq = idx
	}
	return res
}

func getStraightMoves(board Board, idx int, color Piece) []int {
	var res []int
	offsets := []int{8, 1, -8, -1}
	sq := idx
	for _, offset := range offsets {
		sq += offset
		for sq < 64 && sq >= 0 {
			if !board.hasPieceOnIdx(sq) {
				res = append(res, sq)
				sq += offset
				continue
			}
			if board.hasColorPieceOnIdx(sq, color) {
				break
			}
			res = append(res, sq)
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
		{x: -2, y: 9},
		{x: 1, y: -16},
		{x: -1, y: -16},
		{x: 2, y: -8},
		{x: -2, y: -9},
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

func GetKingMoves(board Board, idx int, color Piece, castleRights *CastleRights) []int {
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
			res = append(res, 2)
		}
		if castleRights.BlackQueen && board[3] == None && board[2] == None && board[1] == None {
			res = append(res, 6)
		}
	}
	return res
}
