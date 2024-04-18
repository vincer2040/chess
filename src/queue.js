/** @template T */
class QueueNode {
    /** @type {T} */
    value;

    /** @type {QueueNode<T> | null} */
    next;

    /**
     * @param {T} val
     */
    constructor(val) {
        this.value = val;
        this.next = null;
    }
}

/** @template T*/
export class Queue {
    /** @type {QueueNode<T> | null} */
    #head;

    /** @type {QueueNode<T> | null} */
    #tail;

    /** @type {number} */
    #len;

    constructor() {
        this.#head = this.#tail = null;
        this.#len = 0;
    }

    /**
     * @returns {T | null}
     */
    deque() {
        if (this.#len === 0) {
            return null;
        }
        // @ts-ignore:
        const v = this.#head.value;
        if (this.#len === 1) {
            this.#head = this.#tail = null;
            this.#len--;
            return v;
        }
        // @ts-ignore:
        this.#head = this.#head?.next;
        this.#len--;
        return v;
    }

    /**
     * @param {T} value
     */
    enque(value) {
        const node = new QueueNode(value);
        if (this.#len === 0) {
            this.#head = this.#tail = node;
            this.#len++;
            return;
        }
        // @ts-ignore:
        this.#tail.next = node;
        this.#tail = node;
        this.#len++;
    }

    /**
     * @returns {number}
     */
    getLen() {
        return this.#len;
    }
}
