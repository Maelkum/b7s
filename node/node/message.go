package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

type topicInfo struct {
	handle       *pubsub.Topic
	subscription *pubsub.Subscription
}

// TODO: Reintroduce telemetry here

func (c *Core) SubscribeToTopics(ctx context.Context, topics []string) error {

	err := c.Host.InitPubSub(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize pubsub: %w", err)
	}

	c.Log.Info().Strs("topics", topics).Msg("topics node will subscribe to")

	// c.metrics.IncrCounter(subscriptionsMetric, float32(len(topics)))

	// TODO: If some topics/subscriptions failed, cleanup those already subscribed to.
	for _, topicName := range topics {

		topic, subscription, err := c.Host.Subscribe(topicName)
		if err != nil {
			return fmt.Errorf("could not subscribe to topic (name: %s): %w", topicName, err)
		}

		ti := &topicInfo{
			handle:       topic,
			subscription: subscription,
		}

		_ = ti

		// TODO: Handle this.

		// No need for locking since this initialization is done once on start.
		// c.subgroups.topics[topicName] = ti
	}

	return nil
}

// send serializes the message and sends it to the specified peer.
func (c *Core) Send(ctx context.Context, to peer.ID, msg blockless.Message) error {

	// opts := new(msgSpanConfig).pipeline(pipeline.DirectMessagePipeline()).peer(to).spanOpts()
	// ctx, span := c.tracer.Start(ctx, msgSendSpanName(spanMessageSend, msg.Type()), opts...)
	// defer span.End()

	// saveTraceContext(ctx, msg)

	// Serialize the message.
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not encode record: %w", err)
	}

	// Send message.
	err = c.Host.SendMessage(ctx, to, payload)
	if err != nil {
		return fmt.Errorf("could not send message: %w", err)
	}

	// c.metrics.IncrCounterWithLabels(messagesSentMetric, 1, []metrics.Label{{Name: "type", Value: msg.Type()}})

	return nil
}

// sendToMany serializes the message and sends it to a number of peers. `requireAll` dictates how we treat partial errors.
func (c *Core) SendToMany(ctx context.Context, peers []peer.ID, msg blockless.Message, requireAll bool) error {

	// opts := new(msgSpanConfig).pipeline(pipeline.DirectMessagePipeline()).peers(peers...).spanOpts()
	// ctx, span := c.tracer.Start(ctx, msgSendSpanName(spanMessageSend, msg.Type()), opts...)
	// defer span.End()

	// saveTraceContext(ctx, msg)

	// Serialize the message.
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not encode record: %w", err)
	}

	var errGroup multierror.Group
	for i, peer := range peers {
		i := i
		peer := peer

		errGroup.Go(func() error {
			err := c.Host.SendMessage(ctx, peer, payload)
			if err != nil {
				return fmt.Errorf("peer %v/%v send error (peer: %v): %w", i+1, len(peers), peer.String(), err)
			}

			return nil
		})
	}

	// c.metrics.IncrCounterWithLabels(messagesSentMetric, float32(len(peers)), []metrics.Label{{Name: "type", Value: msg.Type()}})

	retErr := errGroup.Wait()
	if retErr == nil || len(retErr.Errors) == 0 {
		// If everything succeeded => ok.
		return nil
	}

	switch len(retErr.Errors) {
	case len(peers):
		// If everything failed => error.
		return fmt.Errorf("all sends failed: %w", retErr)

	default:
		// Some sends failed - do as requested by `requireAll`.
		if requireAll {
			return fmt.Errorf("some sends failed: %w", retErr)
		}

		c.Log.Warn().Err(retErr).Msg("some sends failed, proceeding")

		return nil
	}
}

func (c *Core) Publish(ctx context.Context, msg blockless.Message) error {
	return c.PublishToTopic(ctx, blockless.DefaultTopic, msg)
}

func (c *Core) PublishToTopic(ctx context.Context, topic string, msg blockless.Message) error {

	// opts := new(msgSpanConfig).pipeline(pipeline.PubSubPipeline(topic)).spanOpts()
	// ctx, span := c.tracer.Start(ctx, msgSendSpanName(spanMessagePublish, msg.Type()), opts...)
	// defer span.End()

	// saveTraceContext(ctx, msg)

	// Serialize the message.
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not encode record: %w", err)
	}

	// TODO: fix this
	var (
		topicInfo *topicInfo
		ok        bool
	)
	// c.subgroups.RLock()
	// topicInfo, ok := c.subgroups.topics[topic]
	// c.subgroups.RUnlock()

	if !ok {
		c.Log.Info().Str("topic", topic).Msg("unknown topic, joining now")

		var err error
		topicInfo, err = c.JoinTopic(topic)
		if err != nil {
			return fmt.Errorf("could not join topic (topic: %s): %w", topic, err)
		}
	}

	// Publish message.
	err = c.Host.Publish(ctx, topicInfo.handle, payload)
	if err != nil {
		return fmt.Errorf("could not publish message: %w", err)
	}

	// c.metrics.IncrCounterWithLabels(messagesPublishedMetric, 1,
	// 	[]metrics.Label{
	// 		{Name: "type", Value: msg.Type()},
	// 		{Name: "topic", Value: topic},
	// 	})

	return nil
}

// wrapper around topic joining + housekeeping.
func (c *Core) JoinTopic(topic string) (*topicInfo, error) {

	// c.subgroups.Lock()
	// defer c.subgroups.Unlock()

	th, err := c.Host.JoinTopic(topic)
	if err != nil {
		return nil, fmt.Errorf("could not join topic (topic: %s): %w", topic, err)
	}

	// NOTE: No subscription, joining topic only.
	ti := &topicInfo{
		handle: th,
	}

	// c.subgroups.topics[topic] = ti

	return ti, nil
}

func (c *Core) Connected(peer peer.ID) bool {
	connections := c.Host.Network().ConnsToPeer(peer)
	return len(connections) > 0
}
