package game

import (
	"fmt"
	"strings"

	"github.com/vincer2040/chess/internal/util"
)

var BOARD_IDXS = [8][8]int{
	[8]int{0, 1, 2, 3, 4, 5, 6, 7},
	[8]int{8, 9, 10, 11, 12, 13, 14, 15},
	[8]int{16, 17, 18, 19, 20, 21, 22, 23},
	[8]int{24, 25, 26, 27, 28, 29, 30, 31},
	[8]int{32, 33, 34, 35, 36, 37, 38, 39},
	[8]int{40, 41, 42, 43, 44, 45, 46, 47},
	[8]int{48, 49, 50, 51, 52, 53, 54, 55},
	[8]int{56, 57, 58, 59, 60, 61, 62, 63},
}

type Board []Piece

func newBoard(pos string) Board {
	var res Board
	b := strings.Split(pos, "/")
	for _, rank := range b {
		for _, ch := range rank {
			if util.IsDigit(byte(ch)) {
				skip := util.ByteToInt(byte(ch))
				for i := 0; i < skip; i++ {
					res = append(res, newPiece(' '))
				}
				continue
			}
			res = append(res, newPiece(byte(ch)))
		}
	}
	return res
}

func (b Board) copy() Board {
    var c Board
    for _, x := range b {
        c = append(c, x)
    }
    return c
}

func (b Board) print() {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			piece := b[(i*8)+j]
			b := piece.GetPieceByte()
			fmt.Printf("%c ", b)
		}
		fmt.Printf("\n")
	}
}

func (b Board) hasPieceOnIdx(idx int) bool {
	square := b[idx]
	return square != None
}

func (b Board) hasColorPieceOnIdx(idx int, color Piece) bool {
	piece := b[idx]
	if piece == None {
		return false
	}
	pieceColor := piece & COLORMASK
	return pieceColor == color
}

func getMinIdxForRank(rank int) int {
	return BOARD_IDXS[rank][0]
}

func getMaxIdxForRank(rank int) int {
	return BOARD_IDXS[rank][7]
}

func getRankForIdx(idx int) int {
	return idx / 8
}

func getFileForIdx(idx int) int {
	return idx % 8
}
