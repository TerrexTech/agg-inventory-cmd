package inventory

import (
	"github.com/TerrexTech/go-eventstore-models/model"
	"github.com/TerrexTech/go-mongoutils/mongo"
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
	switch event.ServiceAction {
	case "createSale":
		return createSale(collection, event)
	default:
		return updateInventory(collection, event)
	}
}
