import { getSquare, getSquareIdx } from "./move";
import { isDigit, createPiece, getPieceUrl } from "./util";

export class Game {
    /** @type {Game} */
    static #instance;

    /** @type {string}*/
    #position;

    /** @type {HTMLImageElement | null}*/
    #moving;

    /** @type {HTMLElement | null}*/
    #fromSquare;

    /** @type {HTMLElement | null}*/
    #toSquare;

    /** @type {HTMLElement}*/
    #board;

    /** @type {string}*/
    #toMove;

    /** @type import("./types").RankFile */
    #moveFrom;

    /** @type import("./types").RankFile */
    #moveTo;

    /** @type {number}*/
    #moveFromIdx;

    /** @type {number}*/
    #moveToIdx;

    /** @type {WebSocket} */
    #ws;

    /**
     * @param {string} startingPosition
     * @param {WebSocket} ws
     */
    constructor(startingPosition, ws) {
        const split = startingPosition.split(" ");
        this.#position = split[0];
        this.#toMove = split[1];
        this.#moveFrom = { rank: 1, file: "a" }
        this.#moveTo = { rank: 1, file: "a" }
        this.#moving = null;
        this.#fromSquare = null;
        this.#toSquare = null;
        this.#moveFromIdx = -1;
        this.#moveToIdx = -1;
        this.#board = /** @type {HTMLElement} */(document.getElementById("board"));
        this.#ws = ws;
        Game.#instance = this;
    }

    drawBoard() {
        const split = this.#position.split("/");
        for (let rank = 0; rank < split.length; ++rank) {
            const curRank = split[rank];
            const expanded = this.#expandRank(curRank);
            const rankEl = this.#board.children.item(rank);
            for (let file = 0; file < expanded.length; ++file) {
                const p = expanded[file];
                const squareEl = rankEl?.children.item(file);
                if (p == " ") {
                    continue;
                }
                const url = /** @type {string}*/(getPieceUrl(p));
                const piece = createPiece(url);
                piece.addEventListener("mousedown", e => Game.#mouseDownCallback(e));
                squareEl?.append(piece);
            }
        }
    }

    #resetBoardColors() {
        for (let rank = 0; rank < 8; ++rank) {
            const rankEl = this.#board.children.item(rank);
            for (let file = 0; file < 8; ++file) {
                const square = rankEl?.children.item(file);
                if (square === this.#fromSquare) {
                    continue;
                }
                if ((file + rank) % 2 === 0) {
                    square?.classList.replace("bg-sky-400", "bg-orange-100");
                } else {
                    square?.classList.replace("bg-sky-400", "bg-sky-800");
                }
            }
        }
    }

    /**
     * @param {MouseEvent} e
     */
    #mouseDown(e) {
        e.preventDefault();
        // @ts-ignore:
        window.addEventListener("mousemove", Game.#mouseMoveCallback);
        window.addEventListener("mouseup", Game.#mouseUpCallback);
        this.#moving = /**@type {HTMLImageElement}*/(e.target);
        const parent = this.#moving.parentElement;
        if (!parent) {
            throw new Error("impossible");
        }
        this.#fromSquare?.classList.replace("bg-orange-100", "bg-sky-400");
        this.#fromSquare?.classList.replace("bg-sky-800", "bg-sky-400");
        this.#moving.style.position = "absolute"
        this.#moving.style.top = `${e.clientY - 48}px`;
        this.#moving.style.left = `${e.clientX - 48}px`

        this.#moveFrom = getSquare(this.#board, parent);
        this.#fromSquare = parent;

        this.#moveFromIdx = getSquareIdx(this.#board, parent);
    }

    /**
     * @param {MouseEvent} e
     */
    #mouseMove(e) {
        if (!this.#moving) {
            throw new Error("we have messed up");
        }
        this.#moving.style.position = "absolute"
        this.#moving.style.top = `${e.clientY - 48}px`;
        this.#moving.style.left = `${e.clientX - 48}px`
    }

    /**
     * @param {MouseEvent} e
     */
    #mouseUp(e) {
        if (!this.#moving) {
            throw new Error("we have messed up");
        }
        this.#moving.remove();
        let square = document.elementFromPoint(e.clientX, e.clientY);
        this.#moving.style.removeProperty("position");
        this.#moving.style.removeProperty("top");
        this.#moving.style.removeProperty("left");
        // is this a capture?
        if (square?.tagName === "IMG") {
            square.replaceWith(this.#moving);
            square = this.#moving.parentElement;
        } else {
            square?.append(this.#moving);
        }
        if (!square) {
            throw new Error("impossible");
        }
        this.#resetBoardColors();
        this.#toSquare = /** @type {HTMLElement}*/(square);
        this.#toSquare.classList.replace("bg-orange-100", "bg-sky-400");
        this.#toSquare.classList.replace("bg-sky-800", "bg-sky-400");
        window.removeEventListener("mouseup", Game.#mouseUpCallback);
        window.removeEventListener("mousemove", Game.#mouseMoveCallback);
        this.#moveTo = getSquare(this.#board, square);
        this.#moveToIdx = getSquareIdx(this.#board, square);
        console.log(this.#moveFrom);
        console.log(this.#moveTo);
        console.log(this.#moveFromIdx, this.#moveToIdx);
    }

    /**
     * @param {string} rank
     */
    #expandRank(rank) {
        let res = "";
        for (let i = 0; i < rank.length; ++i) {
            const cur = rank[i];
            if (isDigit(cur)) {
                const bound = parseInt(cur);
                for (let j = 0; j < bound; ++j) {
                    res += " ";
                }
                continue;
            }
            res += cur;
        }
        return res;
    }

    /**
     * @param {MouseEvent} e
     */
    static #mouseMoveCallback(e) {
        Game.#instance.#mouseMove(e);
    }

    /**
     * @param {MouseEvent} e
     */
    static #mouseDownCallback(e) {
        Game.#instance.#mouseDown(e);
    }

    /**
     * @param {MouseEvent} e
     */
    static #mouseUpCallback(e) {
        Game.#instance.#mouseUp(e);
    }
}
