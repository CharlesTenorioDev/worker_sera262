package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type PerguntIADTO struct {
	Perguntas string `json:"perguntas"`
}

type QuestionarioParaFilaDTO struct {
	ID     primitive.ObjectID `json:"id"`
	Titulo string             `json:"titulo"`
}

type GeneratePayloadGroq struct {
	Model    string        `json:"model"`
	Messages []MessageGroq `json:"messages"`
}

type MessageGroq struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
