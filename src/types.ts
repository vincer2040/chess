export type Move = {
    from: number;
    to: number;
};

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
