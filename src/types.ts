export type Move = {
    from: number;
    to: number;
};

export type Promotion = Move & {
    promoteTo: string;
}

export const DataTypes = {
    Illegal: "illegal",
    Position: "position",
    Move: "move",
    Error: "error",
    Command: "command",
    LegalMoves: "legal moves",
    AttackingMoves: "attacking moves",
    Promotion: "promotion",
} as const;

export type DataType = typeof DataTypes[keyof typeof DataTypes];

export type LegalMoves = Map<number, number[]>;
export type AttackingMoves = Map<number, number[][]>;

export type DataFromServer = {
    type: DataType,
    data: LegalMoves | AttackingMoves | string | Move | Promotion | null;
}
