import { DataTypes } from "./types"

const POSITION_BYTE = 43; // +
const MOVE_BYTE = 36; // $
const SEPARATOR = 58; // :
const RET_CAR = 13; // \r
const NEW_LINE = 10; // \n
const ERROR_BYTE = 45; // -
const COMMAND_BYTE = 35; // #
const LEGAL_MOVES = 126; // ~
const ATTACKING_MOVES = 94; // ^
const ARRAY_BYTE = 42; // *
const PROMOTION_BYTE = 33; // !
const ZERO_BYTE = 48; // 0

export class Parser {
    /** @type {Uint8Array} */
    #buf;

    /** @type {number} */
    #pos;

    /** @type {number} */
    #byte;

    /**
     * @param {Uint8Array} buf
     */
    constructor(buf) {
        this.#buf = buf;
        this.#pos = 0;
        this.#byte = 0;
        this.#readByte();
    }

    /**
     * @returns {import("./types").DataFromServer}
     */
    parse() {
        /** @type {import("./types").LegalMoves | import("./types").AttackingMoves | import("./types").Move | string | null} */
        let data = null;
        /** @type {import("./types").DataType} */
        let type = DataTypes.Illegal;
        switch (this.#byte) {
            case POSITION_BYTE:
                data = this.#parsePosition();
                if (data !== null) {
                    type = DataTypes.Position;
                }
                break;
            case MOVE_BYTE:
                data = this.#parseMove();
                if (data !== null) {
                    type = DataTypes.Move;
                }
                break;
            case COMMAND_BYTE:
                data = this.#parseCommand();
                if (data !== null) {
                    type = DataTypes.Command;
                }
                break
            case ERROR_BYTE:
                data = this.#parseError();
                if (data !== null) {
                    type = DataTypes.Error;
                }
                break
            case LEGAL_MOVES:
                data = this.#parseLegalMoves();
                if (data !== null) {
                    type = DataTypes.LegalMoves;
                }
                break
            case ATTACKING_MOVES:
                data = this.#parseAttackingMoves();
                if (data !== null) {
                    type = DataTypes.AttackingMoves;
                }
                break
            case PROMOTION_BYTE:
                data = this.#parsePromotion();
                if (data !== null) {
                    type = DataTypes.Promotion;
                }
                break
        }
        return { type, data };
    }

    /**
     * @returns {import("./types").Move | null}
     */
    #parseMove() {
        this.#readByte();
        let from = "";
        let to = "";
        while (this.#byte !== SEPARATOR && this.#byte !== 0) {
            from += String.fromCharCode(this.#byte);
            this.#readByte();
        }
        this.#readByte();
        // @ts-ignore:
        while (this.#byte !== RET_CAR && this.#byte !== 0) {
            to += String.fromCharCode(this.#byte);
            this.#readByte();
        }
        if (!this.#expectEnd()) {
            return null;
        }
        return { from: parseInt(from), to: parseInt(to) };
    }

    /**
     * @returns {import("./types").Promotion | null}
     */
    #parsePromotion() {
        this.#readByte()
        let from = "";
        let to = "";
        let promoteTo = "";
        while (this.#byte !== SEPARATOR && this.#byte !== 0) {
            from += String.fromCharCode(this.#byte);
            this.#readByte();
        }
        this.#readByte();
        while (this.#byte !== SEPARATOR && this.#byte !== 0) {
            to += String.fromCharCode(this.#byte);
            this.#readByte();
        }
        this.#readByte();
        promoteTo = String.fromCharCode(this.#byte);
        this.#readByte();
        if (!this.#expectEnd()) {
            return null;
        }
        return { from: parseInt(from), to: parseInt(to), promoteTo };
    }

    /**
     * @returns {string | null}
     */
    #parsePosition() {
        let res = "";
        this.#readByte();
        while (this.#byte != RET_CAR && this.#byte != 0) {
            res += String.fromCharCode(this.#byte);
            this.#readByte();
        }
        if (!this.#expectEnd()) {
            return null;
        }
        return res;
    }

    /**
     * @returns {string | null}
     */
    #parseCommand() {
        let res = "";
        this.#readByte();
        while (this.#byte != RET_CAR && this.#byte != 0) {
            res += String.fromCharCode(this.#byte);
            this.#readByte();
        }
        if (!this.#expectEnd()) {
            return null;
        }
        return res;
    }

    /**
     * @returns {string | null}
     */
    #parseError() {
        let res = "";
        this.#readByte();
        while (this.#byte != RET_CAR && this.#byte != 0) {
            res += String.fromCharCode(this.#byte);
            this.#readByte();
        }
        if (!this.#expectEnd()) {
            return null;
        }
        return res;
    }

    /**
     * @returns {import("./types").LegalMoves | null}
     */
    #parseLegalMoves() {
        /** @type {import("./types").LegalMoves}*/
        const res = new Map();
        this.#readByte();
        let len = this.#parseNumber();
        if (!this.#expectEnd()) {
            return null;
        }
        this.#readByte();
        for (let i = 0; i < len; ++i) {
            let key = this.#parseNumber();
            if (!this.#expectEnd()) {
                return null;
            }
            if (!this.#expectPeek(ARRAY_BYTE)) {
                return null;
            }
            this.#readByte();
            let numMoves = this.#parseNumber();
            if (!this.#expectEnd()) {
                return null;
            }
            this.#readByte();
            /** @type {number[]}*/
            let moves = [];
            while (this.#byte !== RET_CAR && this.#byte !== 0) {
                let move = 0;
                // @ts-ignore:
                while (this.#byte !== SEPARATOR && this.#byte != RET_CAR && this.#byte !== 0) {
                    move = (move * 10) + (this.#byte - ZERO_BYTE);
                    this.#readByte();
                }
                moves.push(move);
                // @ts-ignore:
                if (this.#byte === SEPARATOR) {
                    this.#readByte();
                }
            }
            if (!this.#expectEnd()) {
                return null;
            }
            if (moves.length !== numMoves) {
                return null;
            }
            this.#readByte();
            res.set(key, moves);
        }
        return res;
    }

    /**
     * @returns {import("./types").AttackingMoves | null}
     */
    #parseAttackingMoves() {
        /** @type {import("./types").AttackingMoves}*/
        const res = new Map();
        this.#readByte();
        let len = this.#parseNumber();
        if (!this.#expectEnd()) {
            return null;
        }
        this.#readByte();
        for (let i = 0; i < len; ++i) {
            let key = this.#parseNumber();
            if (!this.#expectEnd()) {
                return null;
            }
            if (!this.#expectPeek(ARRAY_BYTE)) {
                return null;
            }
            this.#readByte();
            let numDirections = this.#parseNumber();
            if (!this.#expectEnd()) {
                return null;
            }
            /** @type {number[][]}*/
            let allMoves = [];
            for (let j = 0; j < numDirections; ++j) {
                if (!this.#expectPeek(ARRAY_BYTE)) {
                    return null;
                }
                this.#readByte();
                let numMoves = this.#parseNumber();
                if (!this.#expectEnd()) {
                    return null;
                }
                this.#readByte();
                /** @type {number[]}*/
                let moves = [];
                for (let k = 0; k < numMoves; ++k) {
                    let move = 0;
                    // @ts-ignore:
                    while (this.#byte !== SEPARATOR && this.#byte != RET_CAR && this.#byte !== 0) {
                        move = (move * 10) + (this.#byte - ZERO_BYTE);
                        this.#readByte();
                    }
                    moves.push(move);
                    if (k !== (numMoves - 1)) {
                        this.#readByte();
                    }
                }
                if (!this.#expectEnd()) {
                    return null;
                }
                this.#readByte();
                allMoves.push(moves);
            }
            res.set(key, allMoves);
        }
        return res;
    }

    /**
     * @returns {number}
     */
    #parseNumber() {
        let res = 0;
        while (this.#byte !== RET_CAR && this.#byte !== 0) {
            res = (res * 10) + (this.#byte - ZERO_BYTE);
            this.#readByte();
        }
        return res;
    }

    /**
     * @param {number} byte
     * @returns {boolean}
     */
    #expectPeek(byte) {
        const peek = this.#peek();
        if (peek !== byte) {
            return false;
        }
        this.#readByte();
        return true;
    }

    /**
     * @returns {number}
     */
    #peek() {
        if (this.#pos >= this.#buf.length) {
            return 0;
        }
        return this.#buf[this.#pos];
    }

    #readByte() {
        if (this.#pos >= this.#buf.length) {
            this.#byte = 0;
            return;
        }
        this.#byte = this.#buf[this.#pos];
        this.#pos++;
    }

    /**
     * @returns {boolean}
     */
    #expectEnd() {
        if (this.#byte !== RET_CAR) {
            return false
        }
        this.#readByte();
        // @ts-ignore:
        if (this.#byte !== NEW_LINE) {
            return false
        }
        return true;
    }
}
