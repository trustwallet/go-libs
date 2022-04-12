package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/trustwallet/go-libs/logging"
)

type RequiredAcks int

const (
	// NoResponse doesn't send any response, the TCP ACK is all you get.
	NoResponse RequiredAcks = 0
	// WaitForLocal waits for only the local commit to succeed before responding.
	WaitForLocal RequiredAcks = 1
	// WaitForAll waits for all in-sync replicas to commit before responding.
	WaitForAll RequiredAcks = -1
)

type ProducerConfig struct {
	Brokers []string
	// Topic must be specified either in ProducerConfig or in WriteMessage method. It can't be specified in both places.
	Topic        string
	MaxAttempts  int
	BatchSize    int
	BatchBytes   int
	BatchTimeout time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	RequiredAcks RequiredAcks
	Async        bool
	Log          bool
}

func (c *ProducerConfig) ToKafkaWriterConfig() kafka.WriterConfig {
	config := kafka.WriterConfig{
		Brokers:      c.Brokers,
		Topic:        c.Topic,
		MaxAttempts:  c.MaxAttempts,
		BatchSize:    c.BatchSize,
		BatchBytes:   c.BatchBytes,
		BatchTimeout: c.BatchTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		RequiredAcks: int(c.RequiredAcks),
		Async:        c.Async,
	}

	if c.Log {
		config.Logger = logging.GetLogger()
	}

	return config
}

type ConsumerConfig struct {
	Brokers                []string
	GroupID                string
	GroupTopics            []string
	Topic                  string
	Partition              int
	QueueCapacity          int
	MinBytes               int
	MaxBytes               int
	MaxWait                time.Duration
	ReadLagInterval        time.Duration
	HeartbeatInterval      time.Duration
	CommitInterval         time.Duration
	PartitionWatchInterval time.Duration
	WatchPartitionChanges  bool
	SessionTimeout         time.Duration
	RebalanceTimeout       time.Duration
	JoinGroupBackoff       time.Duration
	RetentionTime          time.Duration
	StartOffset            int64
	ReadBackoffMin         time.Duration
	ReadBackoffMax         time.Duration
	MaxAttempts            int
	Log                    bool
}

func (c *ConsumerConfig) ToKafkaReaderConfig() kafka.ReaderConfig {
	config := kafka.ReaderConfig{
		Brokers:                c.Brokers,
		GroupID:                c.GroupID,
		GroupTopics:            c.GroupTopics,
		Topic:                  c.Topic,
		Partition:              c.Partition,
		QueueCapacity:          c.QueueCapacity,
		MinBytes:               c.MinBytes,
		MaxBytes:               c.MaxBytes,
		MaxWait:                c.MaxWait,
		ReadLagInterval:        c.ReadLagInterval,
		HeartbeatInterval:      c.HeartbeatInterval,
		CommitInterval:         c.CommitInterval,
		PartitionWatchInterval: c.PartitionWatchInterval,
		WatchPartitionChanges:  c.WatchPartitionChanges,
		SessionTimeout:         c.SessionTimeout,
		RebalanceTimeout:       c.RebalanceTimeout,
		JoinGroupBackoff:       c.JoinGroupBackoff,
		RetentionTime:          c.RetentionTime,
		StartOffset:            c.StartOffset,
		ReadBackoffMin:         c.ReadBackoffMin,
		ReadBackoffMax:         c.ReadBackoffMax,
		MaxAttempts:            c.MaxAttempts,
	}

	if c.Log {
		config.Logger = logging.GetLogger()
	}

	return config
}
