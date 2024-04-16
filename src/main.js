import { isDigit, createPiece } from "./util";
// const url = "http://localhost:8080";

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

/**
 * @param {string} piece
 */
function getPieceUrl(piece) {
    return pieceToUrlMap.get(piece);
}

class Game {

    /** @type {Game} */
    static #instance;

    /** @type {string}*/
    #position;

    /** @type {HTMLImageElement | null}*/
    #moving;

    /** @type {HTMLElement}*/
    #board;

    /**
     * @param {string} startingPosition
     */
    constructor(startingPosition) {
        this.#position = startingPosition;
        this.#moving = null;
        this.#board = /** @type {HTMLElement} */(document.getElementById("board"));
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
        this.#resetBoardColors();
        // @ts-ignore:
        window.addEventListener("mousemove", Game.#mouseMoveCallback);
        window.addEventListener("mouseup", Game.#mouseUpCallback);
        this.#moving = /**@type {HTMLImageElement}*/(e.target);
        const parent = this.#moving.parentElement;
        parent?.classList.replace("bg-orange-100", "bg-sky-400");
        parent?.classList.replace("bg-sky-800", "bg-sky-400");
        this.#moving.style.position = "absolute"
        this.#moving.style.top = `${e.clientY - 48}px`;
        this.#moving.style.left = `${e.clientX - 48}px`
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
        square?.classList.replace("bg-orange-100", "bg-sky-400");
        square?.classList.replace("bg-sky-800", "bg-sky-400");
        window.removeEventListener("mouseup", Game.#mouseUpCallback);
        window.removeEventListener("mousemove", Game.#mouseMoveCallback);
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

const startingPosition = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR";

const game = new Game(startingPosition);

game.drawBoard();
