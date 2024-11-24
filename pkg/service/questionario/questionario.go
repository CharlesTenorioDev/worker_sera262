package questionario

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/sera/back-end/worker/internal/config/logger"
	"github.com/sera/back-end/worker/pkg/adapter/mongodb"
	"github.com/sera/back-end/worker/pkg/adapter/rabbitmq"
	"github.com/sera/back-end/worker/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QuestionarioServiceInterface interface {
	Update(ctx context.Context, ID string, QuestionarioToChange *model.Questionario) (bool, error)
	GetByID(ctx context.Context, ID string) (*model.Questionario, error)
	GetAll(ctx context.Context, filters model.FilterQuestionario, limit, page int64) (*model.Paginate, error)
}

type QuestionarioDataService struct {
	mdb mongodb.MongoDBInterface
	rmb rabbitmq.RabbitInterface
}

func NewQuestionarioervice(mongo_connection mongodb.MongoDBInterface, rabbitmq_connection rabbitmq.RabbitInterface) *QuestionarioDataService {
	return &QuestionarioDataService{
		mdb: mongo_connection,
		rmb: rabbitmq_connection,
	}
}

func (cat *QuestionarioDataService) Update(ctx context.Context, ID string, Questionario *model.Questionario) (bool, error) {
	collection := cat.mdb.GetCollection("cfSera")

	opts := options.Update().SetUpsert(true)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return false, err
	}

	filter := bson.D{

		{Key: "_id", Value: objectID},
		{Key: "data_type", Value: "Questionario"},
	}

	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "titulo", Value: Questionario.Titulo},
			{Key: "enabled", Value: Questionario.Enabled},
			{Key: "updated_at", Value: time.Now().Format(time.RFC3339)},
		},
	}}

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		logger.Error("Error while updating data", err)

		return false, err
	}

	return true, nil
}

func (cat *QuestionarioDataService) GetByID(ctx context.Context, ID string) (*model.Questionario, error) {

	collection := cat.mdb.GetCollection("cfSera")

	Questionario := &model.Questionario{}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return nil, err
	}

	filter := bson.D{
		{Key: "data_type", Value: "Questionario"},
		{Key: "_id", Value: objectID},
	}

	err = collection.FindOne(ctx, filter).Decode(Questionario)
	if err != nil {
		logger.Error("erro ao consultar Questionario", err)
		return nil, err
	}

	return Questionario, nil
}

func (cat *QuestionarioDataService) GetAll(ctx context.Context, filters model.FilterQuestionario, limit, page int64) (*model.Paginate, error) {
	collection := cat.mdb.GetCollection("cfSera")

	query := bson.M{"data_type": "Questionario"}

	if filters.Titulo != "" || filters.Enabled != "" {
		if filters.Titulo != "" {
			query["nome"] = bson.M{"$regex": fmt.Sprintf(".*%s.*", filters.Titulo), "$options": "i"}
		}
		if filters.Enabled != "" {
			enable, err := strconv.ParseBool(filters.Enabled)
			if err != nil {
				logger.Error("erro converter campo enabled", err)
				return nil, err
			}
			query["enabled"] = enable
		}
	}
	count, err := collection.CountDocuments(ctx, query, &options.CountOptions{})

	if err != nil {
		logger.Error("erro ao consultar todas as Alunos", err)
		return nil, err
	}

	pagination := model.NewPaginate(limit, page, count)

	curr, err := collection.Find(ctx, query, pagination.GetPaginatedOpts())
	if err != nil {
		return nil, err
	}

	result := make([]*model.Questionario, 0)
	for curr.Next(ctx) {
		cat := &model.Questionario{}
		if err := curr.Decode(cat); err != nil {
			logger.Error("erro ao consulta todos os Questionarios", err)
		}
		result = append(result, cat)
	}

	pagination.Paginate(result)

	return pagination, nil
}
