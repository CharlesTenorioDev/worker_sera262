package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Questionario struct {
	ID              primitive.ObjectID
	IDTurma         primitive.ObjectID `bson:"turma_id" json:"turma_id"`
	IDMateria       primitive.ObjectID `bson:"materia_id" json:"materia_id"`
	IDProfessor     primitive.ObjectID `bson:"professor_id" json:"professor_id"`
	DataType        string             `bson:"data_type" json:"-"`
	Titulo          string             `bson:"titulo" json:"titulo"`
	Questoes        []Questao          `bson:"questoes" json:"questoes"`
	EnvaidaParaIA   bool               `bson:"enviado_ia" json:"enviado_ia"`
	EnviadaParaFila bool               `bson:"enviado_fila" json:"enviado_fila"`
	CreatedAt       string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt       string             `bson:"updated_at" json:"updated_at,omitempty"`
	Enabled         bool               `bson:"enabled" json:"enabled"`
}

type Questao struct {
	Pergunta     string   `bson:"pergunta" json:"pergunta"`
	DataType     string   `bson:"data_type" json:"-"`
	Alternativas []string `bson:"alternativas" json:"alternativas"`
	Correta      string   `bson:"correta" json:"correta"`
}

type FilterQuestionario struct {
	ID          primitive.ObjectID
	IDTurma     primitive.ObjectID `bson:"turma_id" json:"turma_id"`
	DataType    string             `bson:"data_type" json:"-"`
	IDMateria   primitive.ObjectID `bson:"materia_id" json:"materia_id"`
	IDProfessor primitive.ObjectID `bson:"professor_id" json:"professor_id"`
	Titulo      string             `json:"titulo"`
	Enabled     string             `json:"enabled"`
}
