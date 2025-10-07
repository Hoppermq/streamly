package domain

import "context"

type Datasource string

type AggFct struct {
	Function string   `json:"function"`
	Args     []string `json:"args"`
	Alias    string   `json:"alias"`
}

type SelectClause struct {
	Type     string
	Field    *string
	Function *AggFct
}

type QueryAstRequest struct {
	Select []string   `json:"select" binding:"required"`
	From   Datasource `json:"from" binding:"required"`
}

type QueryRepository interface{}

type QueryUseCase interface {
	QueryEvents(ctx context.Context)
}
