import { Game } from "./game";
const url = "http://localhost:8080";

const startingPosition = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1";

const ws = new WebSocket(url.replace("http", "ws") + "/game")

const game = new Game(startingPosition, ws);

game.drawBoard();

