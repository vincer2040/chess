import { For, Match, Switch, createSignal, on } from "solid-js"

const pieceToUrlMap = new Map([
    ["p", "/pieces/p.svg"],
    ["n", "/pieces/n.svg"],
    ["b", "/pieces/b.svg"],
    ["r", "/pieces/r.svg"],
    ["q", "/pieces/q.svg"],
    ["k", "/pieces/k.svg"],
    ["P", "/pieces/P.svg"],
    ["N", "/pieces/N.svg"],
    ["B", "/pieces/B.svg"],
    ["R", "/pieces/R.svg"],
    ["Q", "/pieces/Q.svg"],
    ["K", "/pieces/K.svg"],
]);

export default function Board() {

    const [position, setPosition] = createSignal("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR");

    return (
        <div class="flex flex-col">
            <For each={new Array(8)}>
                {(_, i) =>
                    <div class="flex">
                        <For each={new Array(8)}>
                            {(_, j) =>
                                <Square rank={i()} file={j()} position={position} />
                            }
                        </For>
                    </div>
                }
            </For>
        </div>
    )
}

/**
 * @param {{ rank: number, file: number, position: import("solid-js/types/reactive/signal").Accessor<string> }} props
 */
function Square({ rank, file, position }) {
    const piece = findPiece(rank, file, position());
    const url = piece ? /**@type {string}*/(pieceToUrlMap.get(piece)) : null;
    return (
        <Switch fallback={
            <div class="w-24 h-24 bg-sky-800">
                <Piece url = {url} />
            </div>
        }>
            <Match when={(rank + file) % 2 === 0}>
                <div class="w-24 h-24 bg-orange-100">
                    <Piece url = {url}/>
                </div>
            </Match>
        </Switch>
    )
}

/**
 * @param {{ url: string | null }} props
 */
function Piece({ url }) {
    if (url === null) {
        return (<></>);
    }
    return(
        <img class="w-24 h-24" src={url} alt=":(" />
    )
}

/**
 * @param {number} rank
 * @param {number} file
 * @param {string} position
 * @returns {string | null}
 */
function findPiece(rank, file, position) {
    const split = position.split("/");
    const curRank = split[rank];
    const expanded = expandRank(curRank);
    const piece = expanded[file];
    return piece === " " ? null : piece;
}

/**
 * @param {string} rank
 * @returns {string}
 */
function expandRank(rank) {
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
 * @param {string} ch
 * @returns {boolean}
 */
function isDigit(ch) {
    const num = ch.charCodeAt(0);
    return 48 <= num && num <= 57;
}
