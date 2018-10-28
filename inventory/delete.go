package inventory

import (
	"encoding/json"
	"log"

	"github.com/TerrexTech/go-eventstore-models/model"
	"github.com/TerrexTech/go-mongoutils/mongo"
	"github.com/pkg/errors"
)

type deleteResult struct {
	DeletedCount int64 `json:"deletedCount,omitempty"`
}

// Delete handles "delete" events.
func Delete(collection *mongo.Collection, event *model.Event) *model.KafkaResponse {
	filter := map[string]interface{}{}

	err := json.Unmarshal(event.Data, &filter)
	if err != nil {
		err = errors.Wrap(err, "Delete: Error while unmarshalling Event-data")
		log.Println(err)
		return &model.KafkaResponse{
			AggregateID:   event.AggregateID,
			CorrelationID: event.CorrelationID,
			Error:         err.Error(),
			ErrorCode:     InternalError,
			UUID:          event.TimeUUID,
		}
	}

	if len(filter) == 0 {
		err = errors.New("blank filter provided")
		err = errors.Wrap(err, "Delete")
		log.Println(err)
		return &model.KafkaResponse{
			AggregateID:   event.AggregateID,
			CorrelationID: event.CorrelationID,
			Error:         err.Error(),
			ErrorCode:     InternalError,
			UUID:          event.TimeUUID,
		}
	}

	deleteStats, err := collection.DeleteMany(filter)
	if err != nil {
		err = errors.Wrap(err, "Delete: Error in DeleteMany")
		log.Println(err)
		return &model.KafkaResponse{
			AggregateID:   event.AggregateID,
			CorrelationID: event.CorrelationID,
			Error:         err.Error(),
			ErrorCode:     DatabaseError,
			UUID:          event.TimeUUID,
		}
	}

	result := &deleteResult{deleteStats.DeletedCount}
	resultMarshal, err := json.Marshal(result)
	if err != nil {
		err = errors.Wrap(err, "Delete: Error marshalling Inventory Delete-result")
		log.Println(err)
		return &model.KafkaResponse{
			AggregateID:   event.AggregateID,
			CorrelationID: event.CorrelationID,
			Error:         err.Error(),
			ErrorCode:     InternalError,
			UUID:          event.TimeUUID,
		}
	}

	return &model.KafkaResponse{
		AggregateID:   event.AggregateID,
		CorrelationID: event.CorrelationID,
		Result:        resultMarshal,
		UUID:          event.TimeUUID,
	}
}
