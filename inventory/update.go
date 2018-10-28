package inventory

import (
	"encoding/json"
	"log"

	"github.com/TerrexTech/uuuid"

	"github.com/TerrexTech/go-eventstore-models/model"
	"github.com/TerrexTech/go-mongoutils/mongo"
	"github.com/pkg/errors"
)

type inventoryUpdate struct {
	Filter map[string]interface{} `json:"filter"`
	Update map[string]interface{} `json:"update"`
}

type updateResult struct {
	MatchedCount  int64 `json:"matchedCount,omitempty"`
	ModifiedCount int64 `json:"modifiedCount,omitempty"`
}

// Update handles "update" events.
func Update(collection *mongo.Collection, event *model.Event) *model.KafkaResponse {
	invUpdate := &inventoryUpdate{}

	err := json.Unmarshal(event.Data, invUpdate)
	if err != nil {
		err = errors.Wrap(err, "Update: Error while unmarshalling Event-data")
		log.Println(err)
		return &model.KafkaResponse{
			AggregateID:   event.AggregateID,
			CorrelationID: event.CorrelationID,
			Error:         err.Error(),
			ErrorCode:     InternalError,
			UUID:          event.TimeUUID,
		}
	}

	if len(invUpdate.Filter) == 0 {
		err = errors.New("blank filter provided")
		err = errors.Wrap(err, "Update")
		log.Println(err)
		return &model.KafkaResponse{
			AggregateID:   event.AggregateID,
			CorrelationID: event.CorrelationID,
			Error:         err.Error(),
			ErrorCode:     InternalError,
			UUID:          event.TimeUUID,
		}
	}
	if len(invUpdate.Update) == 0 {
		err = errors.New("blank update provided")
		err = errors.Wrap(err, "Update")
		log.Println(err)
		return &model.KafkaResponse{
			AggregateID:   event.AggregateID,
			CorrelationID: event.CorrelationID,
			Error:         err.Error(),
			ErrorCode:     InternalError,
			UUID:          event.TimeUUID,
		}
	}
	if invUpdate.Update["itemID"] == (uuuid.UUID{}).String() {
		err = errors.New("found blank itemID in update")
		err = errors.Wrap(err, "Update")
		log.Println(err)
		return &model.KafkaResponse{
			AggregateID:   event.AggregateID,
			CorrelationID: event.CorrelationID,
			Error:         err.Error(),
			ErrorCode:     InternalError,
			UUID:          event.TimeUUID,
		}
	}

	updateStats, err := collection.UpdateMany(invUpdate.Filter, invUpdate.Update)
	if err != nil {
		err = errors.Wrap(err, "Update: Error in UpdateMany")
		log.Println(err)
		return &model.KafkaResponse{
			AggregateID:   event.AggregateID,
			CorrelationID: event.CorrelationID,
			Error:         err.Error(),
			ErrorCode:     DatabaseError,
			UUID:          event.TimeUUID,
		}
	}

	result := &updateResult{
		MatchedCount:  updateStats.MatchedCount,
		ModifiedCount: updateStats.ModifiedCount,
	}
	resultMarshal, err := json.Marshal(result)
	if err != nil {
		err = errors.Wrap(err, "Update: Error marshalling Inventory Update-result")
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
