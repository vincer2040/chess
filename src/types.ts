export type RankFile = {
    rank: Rank;
    file: File;
}

export type Move = {
    from: number;
    to: number;
};

export const Ranks = [1, 2, 3, 4, 5, 6, 7, 8];
export type Rank = typeof Ranks[number];
export const Files = ["a", "b", "c", "d", "e", "f", "g", "h"];
export type File = typeof Files[number];

export const DataTypes = {
    Illegal: "illegal",
    Position: "position",
    Move: "move",
    Error: "error",
    Command: "command",
    LegalMoves: "legal moves"
} as const;

export type DataType = typeof DataTypes[keyof typeof DataTypes];

export type LegalMoves = Map<number, number[]>;

export type DataFromServer = {
    type: DataType,
    data: LegalMoves | string | Move | null;
}
