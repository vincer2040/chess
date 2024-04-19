/**
 * @param {string} ch
 */
export function isDigit(ch) {
    const x = ch.charCodeAt(0);
    return 48 <= x && x <= 57;
}

/**
 * @param {string} url
 */
export function createPiece(url) {
    const img = document.createElement("img");
    img.src = url;
    img.classList.add("w-24", "h-24", "cursor-pointer");
    img.draggable = true;
    return img;
}

/**
 * @param {number} num
 */
export function numToRank(num) {
    const els = [8, 7, 6, 5, 4, 3, 2, 1];
    return els[num];
}

/**
 * @param {number} num
 */
export function numToFile(num) {
    const els = [
        "a", "b", "c", "d", "e", "f", "g", "h"
    ];
    return els[num];
}

const pieceToUrlMap = new Map([
    ["B", "/pieces/B.svg"],
    ["K", "/pieces/K.svg"],
    ["N", "/pieces/N.svg"],
    ["P", "/pieces/P.svg"],
    ["Q", "/pieces/Q.svg"],
    ["R", "/pieces/R.svg"],
    ["b", "/pieces/b.svg"],
    ["k", "/pieces/k.svg"],
    ["n", "/pieces/n.svg"],
    ["p", "/pieces/p.svg"],
    ["q", "/pieces/q.svg"],
    ["r", "/pieces/r.svg"],
]);

const urlToPieceMap = new Map([
    ["/pieces/B.svg", "B"],
    ["/pieces/K.svg", "K"],
    ["/pieces/N.svg", "N"],
    ["/pieces/P.svg", "P"],
    ["/pieces/Q.svg", "Q"],
    ["/pieces/R.svg", "R"],
    ["/pieces/b.svg", "b"],
    ["/pieces/k.svg", "k"],
    ["/pieces/n.svg", "n"],
    ["/pieces/p.svg", "p"],
    ["/pieces/q.svg", "q"],
    ["/pieces/r.svg", "r"],
]);

/**
 * @param {string} piece
 */
export function getPieceUrl(piece) {
    return pieceToUrlMap.get(piece);
}

/**
 * @param {string} url
 */
export function getPieceFromUrl(url) {
    return urlToPieceMap.get(url);
}

/**
 * @param {string} piece
 */
export function pieceColor(piece) {
    return piece.toUpperCase() === piece ? "white" : "black";
}
