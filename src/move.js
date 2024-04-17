import { Files, Ranks } from "./types";

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
 * @param {HTMLElement} board
 * @param {Element} square
 * @returns {import("./types").RankFile}
 */
export function getSquare(board, square) {
    for (let i = 0; i < 8; ++i) {
        const rankEl = board.children.item(i);
        for (let j = 0; j < 8; ++j) {
            const s = rankEl?.children.item(j);
            if (s === square) {
                const rank = getRankFromNum(i);
                const file = getFileFromNum(j);
                return { rank, file };
            }
        }
    }

    throw new Error("invalid move");
}

/**
 * @param {number} num
 * @returns import("./types").Rank
 */
function getRankFromNum(num) {
    return [...Ranks].reverse()[num];
}

/**
 * @param {number} num
 * @returns import("./types").File
 */
function getFileFromNum(num) {
    return Files[num];
}
