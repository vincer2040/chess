import { render } from "solid-js/web"

import Board from "./board"

const root = document.getElementById("root");
if (!root) {
    throw new Error("no root :(");
}

render(() => Board(), root);
