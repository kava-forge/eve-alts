package operators

//go:generate go-enum --sql --marshal --names --values

// ENUM(all, any, none)
type Operator string
