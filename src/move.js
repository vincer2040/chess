/**
 * @param {HTMLElement} board
 * @param {Element} square
 * @returns {number}
 */
export function getSquareIdx(board, square) {
    for (let i = 0; i < 8; ++i) {
        const rankEl = board.children.item(i);
        for (let j = 0; j < 8; ++j) {
            const s = rankEl?.children.item(j);
            if (s === square) {
                return (8 * i) + j;
            }
        }
    }

    throw new Error("invalid move");
}

/**
 * @param {number} num
 * @returns {number}
 */
export function getRankFromIdx(num) {
    return Math.floor(num / 8);
}

/**
 * @param {number} num
 * @returns {number}
 */
export function getFileFromIdx(num) {
    return num % 8;
}

