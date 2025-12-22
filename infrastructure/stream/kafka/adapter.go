package kafka

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	kafkago "github.com/segmentio/kafka-go"

	adapter "src/application/adapter/stream"
)

type KafkaStreamAdapter struct {
	brokers []string
	writer  *kafkago.Writer
}

var _ adapter.IStreamAdapter = (*KafkaStreamAdapter)(nil)

func NewKafkaStreamAdapter(config *adapter.StreamConfig) *KafkaStreamAdapter {
	if config == nil {
		panic("stream/kafka: config is nil")
	}

	brokers, err := parseKafkaBrokers(config.KafkaURI)
	if err != nil {
		panic(fmt.Errorf("stream/kafka: invalid KafkaURI %q: %w", config.KafkaURI, err))
	}

	writer := &kafkago.Writer{
		Addr:                   kafkago.TCP(brokers...),
		AllowAutoTopicCreation: true,
	}

	return &KafkaStreamAdapter{
		brokers: brokers,
		writer:  writer,
	}
}

func (a *KafkaStreamAdapter) Ping(ctx context.Context) error {
	if len(a.brokers) == 0 {
		return fmt.Errorf("stream/kafka: no brokers configured")
	}

	conn, err := kafkago.DialContext(ctx, "tcp", a.brokers[0])
	if err != nil {
		return err
	}

	return conn.Close()
}

func (a *KafkaStreamAdapter) Publish(ctx context.Context, topic string, payload adapter.Payload) error {
	if topic == "" {
		return fmt.Errorf("stream/kafka: topic is required")
	}

	var headers []kafkago.Header

	if payload.TTL != nil {
		ttlMilliseconds := payload.TTL.Milliseconds()
		headers = append(headers, kafkago.Header{
			Key:   "ttl_ms",
			Value: []byte(strconv.FormatInt(ttlMilliseconds, 10)),
		})
	}

	if payload.MaxRetries != nil {
		headers = append(headers, kafkago.Header{
			Key:   "max_retries",
			Value: []byte(strconv.Itoa(*payload.MaxRetries)),
		})
	}

	message := kafkago.Message{
		Topic:   topic,
		Value:   payload.Message,
		Headers: headers,
	}

	if payload.Key != nil {
		message.Key = []byte(*payload.Key)
	}

	return a.writer.WriteMessages(ctx, message)
}

func (a *KafkaStreamAdapter) Subscribe(ctx context.Context, topic string, handler func(payload adapter.Payload) error) error {
	if topic == "" {
		return fmt.Errorf("stream/kafka: topic is required")
	}

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:     a.brokers,
		Topic:       topic,
		StartOffset: kafkago.FirstOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
	})

	go func() {
		defer reader.Close()

		for {
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				return
			}

			payload := adapter.Payload{
				Message: msg.Value,
			}

			if len(msg.Key) > 0 {
				key := string(msg.Key)
				payload.Key = &key
			}

			for _, header := range msg.Headers {
				switch header.Key {
				case "ttl_ms":
					if millis, err := strconv.ParseInt(string(header.Value), 10, 64); err == nil {
						duration := time.Duration(millis) * time.Millisecond
						payload.TTL = &duration
					}
				case "max_retries":
					if retries, err := strconv.Atoi(string(header.Value)); err == nil {
						payload.MaxRetries = &retries
					}
				}
			}

			_ = handler(payload)
		}
	}()

	return nil
}

func parseKafkaBrokers(rawURI string) ([]string, error) {
	value := strings.TrimSpace(rawURI)
	if value == "" {
		return nil, fmt.Errorf("KafkaURI is empty")
	}

	if strings.Contains(value, "://") {
		return nil, fmt.Errorf(
			"KafkaURI must be a comma-separated list of brokers in the form host:port (no scheme). Example: \"localhost:29092,localhost:29093\"",
		)
	}

	parts := strings.Split(value, ",")
	brokers := make([]string, 0, len(parts))

	for _, part := range parts {
		hostPort := strings.TrimSpace(part)
		if hostPort == "" {
			continue
		}

		if err := validateHostPort(hostPort); err != nil {
			return nil, fmt.Errorf("invalid broker %q: %w", hostPort, err)
		}

		brokers = append(brokers, hostPort)
	}

	if len(brokers) == 0 {
		return nil, fmt.Errorf("no brokers parsed from %q", rawURI)
	}

	return brokers, nil
}

func validateHostPort(hostPort string) error {
	host, portStr, err := net.SplitHostPort(hostPort)
	if err != nil {
		return fmt.Errorf("expected host:port")
	}

	host = strings.TrimSpace(host)
	if host == "" {
		return fmt.Errorf("host is empty")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("port is not a number")
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("port out of range")
	}

	if strings.ContainsAny(host, "/?") || strings.ContainsAny(portStr, "/?") {
		return fmt.Errorf("broker must not contain path or query")
	}

	return nil
}
