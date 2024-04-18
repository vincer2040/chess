const POSITION_BYTE = 43 // +
const MOVE_BYTE = 36 // $
const MOVE_SEPERATOR = 58 // :
const RET_CAR = 13 // \r
const NEW_LINE = 10 // \n
const COMMAND_BYTE = 35 // #

export class Builder {
    /** @type {ArrayBuffer} */
    #buf;

    /** @type {Uint8Array} */
    #view;

    /** @type {number} */
    #len;

    /** @type {number} */
    #capacity;

    constructor() {
        this.#len = 0;
        this.#capacity = 32;
        this.#buf = new ArrayBuffer(this.#capacity);
        this.#view = new Uint8Array(this.#buf);
    }

    /**
     * @param {string} position
     * @returns {Builder}
     */
    addPosition(position) {
        if (this.#len == this.#capacity) {
            this.#resize();
        }
        this.#view[this.#len] = POSITION_BYTE;
        this.#len++;
        for (let i = 0; i < position.length; ++i) {
            if (this.#len == this.#capacity) {
                this.#resize();
            }
            this.#view[this.#len] = position.charCodeAt(i);
            this.#len++;
        }
        this.#addEnd();
        return this;
    }

    /**
     * @param {import("./types").Move} move
     * @returns {Builder}
     */
    addMove(move) {
        if (this.#len == this.#capacity) {
            this.#resize();
        }
        this.#view[this.#len] = MOVE_BYTE;
        this.#len++;

        const fromString = move.from.toString();
        for (let i = 0; i < fromString.length; ++i) {
            if (this.#len == this.#capacity) {
                this.#resize();
            }
            this.#view[this.#len] = fromString.charCodeAt(i);
            this.#len++;
        }

        if (this.#len == this.#capacity) {
            this.#resize();
        }
        this.#view[this.#len] = MOVE_SEPERATOR;
        this.#len++;

        const toString = move.to.toString();
        for (let i = 0; i < toString.length; ++i) {
            if (this.#len == this.#capacity) {
                this.#resize();
            }
            this.#view[this.#len] = toString.charCodeAt(i);
            this.#len++;
        }
        this.#addEnd();
        return this;
    }

    /**
     * @param {string} command
     * @returns {Builder}
     */
    addCommand(command) {
        if (this.#len == this.#capacity) {
            this.#resize();
        }
        this.#view[this.#len] = COMMAND_BYTE;
        this.#len++;
        for (let i = 0; i < command.length; ++i) {
            if (this.#len === this.#capacity) {
                this.#resize();
            }
            this.#view[this.#len] = command.charCodeAt(i);
            this.#len++;
        }
        this.#addEnd();
        return this;
    }

    getBuf() {
        this.#resize();
        return this.#buf
    }

    reset() {
        this.#len = 0;
        this.#view.fill(0);
        return this;
    }

    #addEnd() {
        if (this.#len == this.#capacity) {
            this.#resize();
        }
        this.#view[this.#len] = RET_CAR;
        this.#len++;
        if (this.#len == this.#capacity) {
            this.#resize();
        }
        this.#view[this.#len] = NEW_LINE;
        this.#len++;
    }

    #resize() {
        this.#capacity <<= 1
        const newBuf = new ArrayBuffer(this.#capacity);
        const newView = new Uint8Array(newBuf);
        for (let i = 0; i < this.#view.length; ++i) {
            newView[i] = this.#view[i];
        }
        this.#buf = newBuf;
        this.#view = newView;
    }
}
