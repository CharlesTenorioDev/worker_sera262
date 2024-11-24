package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sera/back-end/worker/internal/config/logger"
)

func (rbm *rbm_pool) SenderRb(ctx context.Context, exchangeName string, routingKey string, msg *Message) error {

	// Verificar se o canal está inicializado
	if rbm.channel == nil {
		logger.Info("RMB.CHANNEL E NULL")
		return fmt.Errorf("canal RabbitMQ não está inicializado")
	}

	// Verificar se o contexto está inicializado
	if ctx == nil {
		logger.Info("ctx E NULL")
		return fmt.Errorf("contexto não está inicializado")
	}

	// Publicar a mensagem
	err := rbm.channel.PublishWithContext(ctx,
		exchangeName, // Nome da exchange
		routingKey,   // Chave de roteamento (nome da fila)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Body:        msg.Data,
			ContentType: msg.ContentType,
		})

	if err != nil {
		logger.Error("Erro ao publicar mensagem", err)
		return err
	}

	logger.Info("Mensagem enviada com sucesso!")
	return nil
}
