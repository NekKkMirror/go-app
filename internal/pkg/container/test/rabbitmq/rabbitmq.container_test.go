package rabbitmqcontainer

import (
	"context"
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_StartRabbitMQContainer(t *testing.T) {
	conn, err := Start(context.Background(), t)
	require.NoError(t, err)

	assert.NotNil(t, conn)

	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	verifyDefaultExchange(t, conn)
}

func Test_StartUsesDefaultRabbitMQOptions(t *testing.T) {
	ctx := context.Background()
	conn, err := Start(ctx, t)
	require.NoError(t, err)
	require.NotNil(t, conn)

	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	verifyDefaultExchange(t, conn)
}

func Test_StartCleansUpContainer(t *testing.T) {
	ctx := context.Background()
	conn, err := Start(ctx, t)
	require.NoError(t, err)
	require.NotNil(t, conn)

	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	verifyDefaultExchange(t, conn)
}

func Test_StartHandlesContextCancellationOrTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	conn, err := Start(ctx, t)
	require.Error(t, err)
	require.Nil(t, conn)
}

func verifyDefaultExchange(t *testing.T, conn *amqp.Connection) {
	channel, err := conn.Channel()
	require.NoError(t, err)
	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {

		}
	}(channel)

	exchangeName := "default"

	err = channel.ExchangeDeclarePassive(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	require.NoError(t, err)
}
