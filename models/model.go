package models


type Result struct {
	Source        string		`json:"source"`
	SparsityAfter string		`json:"sparsityAfter"`
	RowsAfter     int64			`json:"rowsAfter"`
	Task          interface{}	`json:"task"`
}
