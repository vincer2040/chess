import { DataTypes } from "./types"

const POSITION_BYTE = 43 // +
const MOVE_BYTE = 36 // $
const MOVE_SEPERATOR = 58 // :
const RET_CAR = 13 // \r
const NEW_LINE = 10 // \n
const ERROR_BYTE = 45; // -
const COMMAND_BYTE = 35 // #

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
        /** @type {import("./types").Move | string | null} */
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
        }
        return { type, data};
    }

    /**
     * @returns {import("./types").Move | null}
     */
    #parseMove() {
        this.#readByte();
        let from = "";
        let to = "";
        while (this.#byte !== MOVE_SEPERATOR && this.#byte !== 0) {
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
