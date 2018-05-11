package models

type MsgResult struct {
	Source        string      `json:"source"`
	SparsityAfter string      `json:"sparsityAfter"`
	RowsAfter     int64       `json:"rowsAfter"`
	Task          interface{} `json:"task"`
}

type MsgSource struct {
	Source string
	Uuid   string
}

type MsgTask struct {
	Dimension [][]string
}