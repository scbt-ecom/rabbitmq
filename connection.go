package rabbitmq

import (
	"context"
	"log/slog"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/scbt-ecom/slogging"
)

func createConnection(rabbitmqURL string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

type Publisher struct {
	url      string
	conn     *amqp.Connection
	channels map[string]*amqp.Channel
	mu       sync.RWMutex
}

func NewPublisher(rabbitmqURL string, channelNames ...string) *Publisher {
	conn, err := createConnection(rabbitmqURL)
	if err != nil {
		slogging.L(context.Background()).Fatal("failed to init RabbitMQ connection, exiting",
			slogging.ErrAttr(err))
	}

	return &Publisher{
		url:      rabbitmqURL,
		conn:     conn,
		channels: createChannels(conn, channelNames...),
	}
}

func (p *Publisher) StartLifeSupport(reconnectInterval time.Duration) {
	connClose := p.conn.NotifyClose(make(chan *amqp.Error))
	//chansClose := p.createChannelsSigClose()

	for {
		select {
		case amqpErr := <-connClose:
			slog.Error("disconnected from RabbitMQ",
				slogging.ErrAttr(amqpErr))
			for {
				newConn, err := createConnection(p.url)
				if err != nil {
					slog.Error("failed to reconnect to RabbitMQ",
						slogging.ErrAttr(err))
					time.Sleep(reconnectInterval)
					continue
				}

				p.mu.Lock()
				p.conn = newConn
				connClose = p.conn.NotifyClose(make(chan *amqp.Error))
				err = p.recreateChannels()
				if err != nil {
					slog.Error("failed to recreate channels",
						slogging.ErrAttr(err))
				}

				//chansClose = p.createChannelsSigClose()

				p.mu.Unlock()
				slog.Info("successfully reconnected to RabbitMQ")

				break
			}
			//case amqpErr := <-chansClose:
			//	slog.Error("connection with RabbitMQ channel closed",
			//		slog.String("error", amqpErr.Error()))
			//
			//	if p.conn.IsClosed() {
			//		continue
			//	}
			//
			//	p.mu.Lock()
			//	err := p.recreateChannels()
			//	if err != nil {
			//		slog.Error("failed to recreate channels",
			//			slog.String("error", err.Error()))
			//	}
			//
			//	chansClose = p.createChannelsSigClose()
			//	p.mu.Unlock()
		}
	}
}

func (p *Publisher) recreateChannels() error {
	for k, _ := range p.channels {
		newCh, err := p.conn.Channel()
		if err != nil {
			return err
		}

		p.channels[k] = newCh
	}

	return nil
}

func (p *Publisher) createChannelsSigClose() <-chan *amqp.Error {
	chans := make([]<-chan *amqp.Error, 1)

	for _, v := range p.channels {
		chans = append(chans, v.NotifyClose(make(chan *amqp.Error, 1)))
	}

	return fanin(context.Background(), chans...)
}

func fanin(ctx context.Context, chans ...<-chan *amqp.Error) <-chan *amqp.Error {
	out := make(chan *amqp.Error, 1)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, ch := range chans {
		go func(ch <-chan *amqp.Error) {
			select {
			case <-ctx.Done():
				return
			case err, ok := <-ch:
				if !ok {
					return
				}

				select {
				case out <- err:
					return
				default:
				}
				cancel()
			}
		}(ch)
	}

	return out
}

func (p *Publisher) GetChannel(name string) *amqp.Channel {
	p.mu.RLock()
	ch := p.channels[name]
	p.mu.RUnlock()

	return ch
}

func createChannels(conn *amqp.Connection, channelNames ...string) map[string]*amqp.Channel {
	chs := make(map[string]*amqp.Channel, len(channelNames))

	for _, name := range channelNames {

		ch, err := conn.Channel()
		if err != nil {
			slogging.L(context.Background()).Fatal("failed to init RabbitMQ channel, exiting",
				slogging.ErrAttr(err),
				slogging.StringAttr("channel_name", name))
		}

		chs[name] = ch
	}

	return chs
}
