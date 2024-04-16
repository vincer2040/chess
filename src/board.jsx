import { For, Match, Switch } from "solid-js"

export default function Board() {
    return (
        <div class="flex flex-col">
            <For each={new Array(8)}>
                {(_, i) =>
                    <div class="flex">
                        <For each={new Array(8)}>
                            {(_, j) => <Square rank={i()} file={j()} />}
                        </For>
                    </div>
                }
            </For>
        </div>
    )
}

/**
 * @param {{ rank: number, file: number }} props
 */
function Square({ rank, file }) {
    return (
        <Switch fallback={<div class="w-24 h-24 bg-sky-800"></div>}>
            <Match when={(rank + file) % 2 === 0}>
                <div class="w-24 h-24 bg-orange-100"></div>
            </Match>
        </Switch>
    )
}
