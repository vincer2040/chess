import { Game } from "./game";
import { Builder } from "./protocolBuilder";
import { Parser } from "./protocolParser";
const url = "http://localhost:8080";

const startingPosition = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1";

const ws = new WebSocket(url.replace("http", "ws") + "/game")

const game = new Game(startingPosition, ws);

ws.addEventListener("message", (e) => {
    const p = new Parser(new TextEncoder().encode(e.data));
    const d = p.parse();
    console.log(d);
});

const buf = new Builder().addCommand("LEGAL_MOVES").getBuf();

ws.addEventListener("open", () => {
    ws.send(buf);
});

game.drawBoard();

