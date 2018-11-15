package inventory

import (
	"github.com/TerrexTech/go-eventstore-models/model"
	"github.com/TerrexTech/go-mongoutils/mongo"
	"github.com/coreos/etcd/clientv3"
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
func Update(
	etcd *clientv3.Client,
	collection *mongo.Collection,
	event *model.Event,
) *model.Document {
	switch event.ServiceAction {
	case "createSale":
		return createSale(etcd, collection, event)
	default:
		return updateInventory(collection, event)
	}
}
