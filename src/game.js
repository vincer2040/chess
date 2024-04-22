import { getFileFromIdx, getRankFromIdx, getSquareIdx } from "./move";
import { Builder } from "./protocolBuilder";
import { Parser } from "./protocolParser";
import { DataTypes } from "./types";
import { isDigit, createPiece, getPieceUrl, getPieceFromUrl, pieceColor } from "./util";
import { Queue } from "./queue";

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

    /** @type {number}*/
    #moveFromIdx;

    /** @type {number}*/
    #moveToIdx;

    /** @type {WebSocket} */
    #ws;

    /** @type {import("./types").LegalMoves} */
    #legalMoves;

    /** @type {import("./types").AttackingMoves} */
    #attackingMoves;

    /** @type {Queue<ArrayBuffer>}*/
    #messageQueue;

    /**
     * @param {string} startingPosition
     * @param {WebSocket} ws
     */
    constructor(startingPosition, ws) {
        const split = startingPosition.split(" ");
        this.#position = split[0];
        this.#toMove = split[1];
        this.#moving = null;
        this.#fromSquare = null;
        this.#toSquare = null;
        this.#moveFromIdx = -1;
        this.#moveToIdx = -1;
        this.#board = /** @type {HTMLElement} */(document.getElementById("board"));
        this.#legalMoves = new Map();
        this.#attackingMoves = new Map();
        this.#ws = ws;
        this.#ws.addEventListener("message", Game.#handleMessageCallback);
        this.#messageQueue = new Queue();

        let s = new Builder().addCommand("START").getBuf();
        this.#messageQueue.enque(s);
        let lm = new Builder().addCommand("LEGAL_MOVES").getBuf();
        this.#messageQueue.enque(lm);
        let am = new Builder().addCommand("ATTACKING_MOVES").getBuf();
        this.#messageQueue.enque(am);
        this.#ws.addEventListener("open", () => {
            const start = this.#messageQueue.deque();
            if (start === null) {
                throw new Error("impossible");
            }
            this.#ws.send(start);
        });
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

    /**
     * @param {Uint8Array} message
     */
    #handleMessage(message) {
        const data = new Parser(message).parse();
        switch (data.type) {
            case DataTypes.Command:
                break;
            case DataTypes.LegalMoves:
                this.#legalMoves = /** @type {import("./types").LegalMoves} */(data.data);
                break
            case DataTypes.AttackingMoves:
                this.#attackingMoves = /** @type {import("./types").AttackingMoves} */(data.data);
                this.#showAttackingMoves();
                break
            case DataTypes.Move:
                break
        }
        const next = this.#messageQueue.deque();
        if (next) {
            this.#ws.send(next);
        }
    }

    /**
     * @param {import("./types").Move} move
     */
    #emitMove(move) {
        let m = new Builder().addMove(move).getBuf();
        this.#ws.send(m);
        const lm = new Builder().addCommand("LEGAL_MOVES").getBuf();
        this.#messageQueue.enque(lm);
        const am = new Builder().addCommand("ATTACKING_MOVES").getBuf();
        this.#messageQueue.enque(am);
    }

    /**
     * @param {import("./types").Promotion} promotion
     */
    #emitPromotion(promotion) {
        let m = new Builder().addPromotion(promotion).getBuf();
        this.#ws.send(m);
        const lm = new Builder().addCommand("LEGAL_MOVES").getBuf();
        this.#messageQueue.enque(lm);
        const am = new Builder().addCommand("ATTACKING_MOVES").getBuf();
        this.#messageQueue.enque(am);
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

    #resetBoardColorsFromShowLegalMoves() {
        for (let rank = 0; rank < 8; ++rank) {
            const rankEl = this.#board.children.item(rank);
            for (let file = 0; file < 8; ++file) {
                const square = rankEl?.children.item(file);
                if (square === this.#fromSquare) {
                    continue;
                }
                if ((file + rank) % 2 === 0) {
                    square?.classList.replace("bg-red-500", "bg-orange-100");
                } else {
                    square?.classList.replace("bg-red-500", "bg-sky-800");
                }
            }
        }
    }

    #showAttackingMoves() {
        for (const [_, value] of this.#attackingMoves) {
            for (const dir of value) {
                for (const idx of dir) {
                    const rank = getRankFromIdx(idx);
                    const file = getFileFromIdx(idx);
                    const square = this.#board.children.item(rank)?.children.item(file);
                    if ((file + rank) % 2 === 0) {
                        square?.classList.replace("bg-orange-100", "bg-green-500");
                    } else {
                        square?.classList.replace("bg-sky-800", "bg-green-500");
                    }
                }
            }
        }
    }

    #resetBoardColorsFromShowAttackingMoves() {
        for (let rank = 0; rank < 8; ++rank) {
            for (let file = 0; file < 8; ++file) {
                const square = this.#board.children.item(rank)?.children.item(file);
                if ((file + rank) % 2 === 0) {
                    square?.classList.replace("bg-green-500", "bg-orange-100");
                } else {
                    square?.classList.replace("bg-green-500", "bg-sky-800");
                }
            }
        }
    }

    /**
     * @param {import("./types").Move} move
     * @returns {boolean}
     */
    #isCastle(move) {
        // @ts-ignore:
        const piece = getPieceFromUrl(this.#moving.src.replace("http://localhost:8080", ""));
        if (piece !== "K" && piece !== "k") {
            return false
        }
        const amtMoved = Math.abs(move.to - move.from);
        return amtMoved === 2;
    }

    /**
     * @param {import("./types").Move} move
     * @returns {boolean}
     */
    #isPromotion(move) {
        // @ts-ignore:
        const piece = getPieceFromUrl(this.#moving.src.replace("http://localhost:8080", ""));
        if (piece !== "P" && piece !== "p") {
            return false
        }
        if (piece === 'P') {
            return move.to >= 0 && move.to <= 7;
        }
        return move.to >= 56 && move.to <= 63;
    }

    /**
     * @param {import("./types").Move} move
     * @returns {[boolean, HTMLImageElement | null]}
     */
    #isEnPassant(move) {
        // @ts-ignore:
        const piece = getPieceFromUrl(this.#moving.src.replace("http://localhost:8080", ""));
        if (piece !== 'P' && piece !== 'p') {
            return [false, null];
        }
        const amtMoved = Math.abs(move.to - move.from);
        if (amtMoved !== 7 && amtMoved !== 9) {
            return [false, null];
        }
        const rank = getRankFromIdx(this.#moveFromIdx);
        const file = getFileFromIdx(this.#moveFromIdx);
        if (amtMoved === 9) {
            const captured = this.#board?.children.item(rank)?.children.item(piece === 'P' ? file - 1 : file + 1)?.children.item(0);
            // @ts-ignore:
            return [true, captured];
        }
        const captured = this.#board?.children.item(rank)?.children.item(piece === 'P' ? file + 1 : file - 1)?.children.item(0);
        // @ts-ignore:
        return [true, captured];
    }

    /**
     * @param {import("./types").Move} move
     */
    #castle(move) {
        if (!this.#moving) {
            throw new Error("no moving");
        }
        const piece = getPieceFromUrl(this.#moving.src.replace("http://localhost:8080", ""));
        // @ts-ignore:
        const color = pieceColor(piece);
        if (color === "white") {
            if (move.to > move.from) {
                // castle king side
                const rook = this.#board.children.item(7)?.children.item(7)?.children.item(0);
                if (!rook) {
                    throw new Error("no rook");
                }
                rook.remove();
                this.#board.children.item(7)?.children.item(5)?.append(rook);
            } else {
                // castle queen side
                const rook = this.#board.children.item(7)?.children.item(0)?.children.item(0);
                if (!rook) {
                    throw new Error("no rook");
                }
                rook.remove();
                this.#board.children.item(7)?.children.item(3)?.append(rook);
            }
        } else {
            if (move.to > move.from) {
                // castle king side
                const rook = this.#board.children.item(0)?.children.item(7)?.children.item(0);
                if (!rook) {
                    throw new Error("no rook");
                }
                rook.remove();
                this.#board.children.item(0)?.children.item(5)?.append(rook);
            } else {
                // castle queen side
                const rook = this.#board.children.item(0)?.children.item(0)?.children.item(0);
                if (!rook) {
                    throw new Error("no rook");
                }
                rook.remove();
                this.#board.children.item(0)?.children.item(3)?.append(rook);
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

        this.#fromSquare = parent;
        this.#fromSquare?.classList.replace("bg-orange-100", "bg-sky-400");
        this.#fromSquare?.classList.replace("bg-sky-800", "bg-sky-400");
        this.#moving.style.position = "absolute"
        this.#moving.style.top = `${e.clientY - 48}px`;
        this.#moving.style.left = `${e.clientX - 48}px`

        this.#moveFromIdx = getSquareIdx(this.#board, parent);

        const pieceLegalMoves = this.#legalMoves.get(this.#moveFromIdx);
        this.#drawLegalMoves(pieceLegalMoves);
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
        let isCapture = false;
        // is this a capture?
        if (square?.tagName === "IMG") {
            isCapture = true;
            square = square.parentElement;
        }
        if (!square) {
            throw new Error("impossible");
        }

        this.#resetBoardColors();
        this.#resetBoardColorsFromShowLegalMoves();
        window.removeEventListener("mouseup", Game.#mouseUpCallback);
        window.removeEventListener("mousemove", Game.#mouseMoveCallback);

        this.#toSquare = /** @type {HTMLElement}*/(square);
        this.#toSquare.classList.replace("bg-sky-800", "bg-sky-400");
        this.#toSquare.classList.replace("bg-orange-100", "bg-sky-400");
        this.#moveToIdx = getSquareIdx(this.#board, square);

        const legalMovesForPiece = this.#legalMoves.get(this.#moveFromIdx);
        if (!legalMovesForPiece) {
            this.#fromSquare?.append(this.#moving);
            this.#fromSquare = null;
            this.#toSquare = null;
            this.#resetBoardColors();
            return;
        }
        if (!legalMovesForPiece.includes(this.#moveToIdx)) {
            this.#fromSquare?.append(this.#moving);
            this.#fromSquare = null;
            this.#toSquare = null;
            this.#resetBoardColors();
            return;
        }

        /** @type {import("./types").Move}*/
        const move = {
            from: this.#moveFromIdx,
            to: this.#moveToIdx,
        };

        const iscastle = this.#isCastle(move);

        if (!isCapture) {
            this.#toSquare.append(this.#moving);
            if (iscastle) {
                this.#castle(move);
            }
            const [isEnPassant, captured] = this.#isEnPassant(move);
            if (isEnPassant) {
                captured?.remove();
            }
        } else {
            this.#toSquare.replaceChildren(this.#moving);
        }

        this.#moving = null;
        this.#fromSquare = null;
        this.#toSquare = null;
        this.#emitMove(move);
        this.#resetBoardColorsFromShowAttackingMoves();
    }

    /**
     * @param {number[] | undefined} legalMoves
     */
    #drawLegalMoves(legalMoves) {
        if (!legalMoves) {
            return;
        }
        const ranks = this.#board.children;
        legalMoves.forEach((idx) => {
            const rank = getRankFromIdx(idx);
            const file = getFileFromIdx(idx);
            const square = ranks.item(rank)?.children.item(file);
            square?.classList.replace("bg-orange-100", "bg-red-500");
            square?.classList.replace("bg-sky-800", "bg-red-500");
            square?.classList.replace("bg-sky-400", "bg-red-500");
        });
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

    /**
     * @param {MessageEvent<any>} e
     */
    static #handleMessageCallback(e) {
        Game.#instance.#handleMessage(new TextEncoder().encode(e.data));
    }
}

