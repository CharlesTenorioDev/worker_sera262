package main

import (
	"log"

	"github.com/sera/back-end/worker/internal/config"
	"github.com/sera/back-end/worker/internal/config/logger"
	"github.com/sera/back-end/worker/pkg/adapter/mongodb"
	"github.com/sera/back-end/worker/pkg/adapter/rabbitmq"
	"github.com/sera/back-end/worker/pkg/service"
)

var (
	VERSION = "0.1.0-dev"
	COMMIT  = "ABCDEFG-dev"
)

func main() {

	logger.Info("start Application Sera 462 API")
	conf := config.NewConfig()

	mogDbConn := mongodb.New(conf)

	fila := []rabbitmq.Fila{
		{
			Name:       "QUEUE_ENVIAR_IA",
			Durable:    true,
			AutoDelete: false,
			Exclusive:  false,
		},
	}

	rbtMQConn := rabbitmq.NewRabbitMQ(fila, conf)

	if err := rbtMQConn.Connect(); err != nil {
		logger.Error("Falha ao conectar no RabbitMQ", err)

	}

	if err := rbtMQConn.DeclareQueues(); err != nil {
		logger.Error("Falha ao declarar filas no RabbitMQ", err)

	}

	isAlive, err := rbtMQConn.IsAlive()
	if !isAlive || err != nil {
		logger.Error("RabbitMQ não está disponível", err)

	}

	rbmq_conn := rabbitmq.NewRabbitMQ(fila, conf)
	task_service := service.NewTaskMessageCounterService(rbmq_conn, db_conn)

	done := make(chan bool)
	go task_service.Run()
	log.Printf("Worker Running [Mode: %s], [Version: %s], [Commit: %s]", conf.Mode, VERSION, COMMIT)
	<-done
}
