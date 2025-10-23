// internal/domain/shared/domain_event.go
package shared

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// DomainEvent representa um evento que ocorreu no domínio
type DomainEvent interface {
	// EventType retorna o tipo do evento
	EventType() string

	// OccurredAt retorna quando o evento ocorreu
	OccurredAt() time.Time

	// AggregateID retorna o ID da entidade que gerou o evento
	AggregateID() uuid.UUID

	// EventData retorna os dados do evento
	EventData() interface{}
}

// BaseDomainEvent implementação base
type BaseDomainEvent struct {
	eventType   string
	occurredAt  time.Time
	aggregateID uuid.UUID
	eventData   interface{}
}

// NewDomainEvent cria um novo evento de domínio
func NewDomainEvent(eventType string, aggregateID uuid.UUID, data interface{}) DomainEvent {
	return &BaseDomainEvent{
		eventType:   eventType,
		occurredAt:  time.Now(),
		aggregateID: aggregateID,
		eventData:   data,
	}
}

func (e *BaseDomainEvent) EventType() string {
	return e.eventType
}

func (e *BaseDomainEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e *BaseDomainEvent) AggregateID() uuid.UUID {
	return e.aggregateID
}

func (e *BaseDomainEvent) EventData() interface{} {
	return e.eventData
}

// EventStore interface para armazenar eventos
type EventStore interface {
	Save(ctx context.Context, event DomainEvent) error
	GetByAggregateID(ctx context.Context, aggregateID uuid.UUID) ([]DomainEvent, error)
	GetByEventType(ctx context.Context, eventType string) ([]DomainEvent, error)
	GetByDateRange(ctx context.Context, start, end time.Time) ([]DomainEvent, error)
}

// EventPublisher interface para publicar eventos
type EventPublisher interface {
	Publish(ctx context.Context, event DomainEvent) error
	PublishMany(ctx context.Context, events []DomainEvent) error
}

// EventHandler interface para handlers de eventos
type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
	CanHandle(eventType string) bool
}

// EventBus interface para gerenciar eventos
type EventBus interface {
	Subscribe(eventType string, handler EventHandler)
	Unsubscribe(eventType string, handler EventHandler)
	Publish(ctx context.Context, event DomainEvent) error
	PublishMany(ctx context.Context, events []DomainEvent) error
}
