package inventory

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"

	"github.com/TerrexTech/go-kafkautils/kafka"

	"github.com/TerrexTech/go-commonutils/commonutil"
	"github.com/TerrexTech/go-eventstore-models/model"
	"github.com/TerrexTech/go-mongoutils/mongo"
	"github.com/TerrexTech/uuuid"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/pkg/errors"
)

// SaleItemResult is the result from updating the sale-item.
type SaleItemResult struct {
	ItemID          uuuid.UUID `json:"itemID,omitempty"`
	Error           string     `json:"error,omitempty"`
	ErrorCode       int        `json:"errorCode,omitempty"`
	TotalSoldWeight float64    `json:"totalSoldWeight,omitempty"`
	TotalWeight     float64    `json:"totalWeight,omitempty"`
}

// SaleValidationResp is the response when a sale is validated.
type SaleValidationResp struct {
	OriginalRequest map[string]interface{} `json:"originalRequest,omitempty"`
	Result          []SaleItemResult       `json:"result,omitempty"`
}

var producer *kafka.Producer

func createSale(
	etcd *clientv3.Client,
	collection *mongo.Collection,
	event *model.Event,
) *model.Document {
	m := map[string]interface{}{}
	err := json.Unmarshal(event.Data, &m)
	if err != nil {
		err = errors.Wrap(err, "SaleCreated-Event: Error unmarshalling sale-data")
		log.Println(err)
		return &model.Document{
			AggregateID:   event.AggregateID,
			CorrelationID: event.CorrelationID,
			Error:         err.Error(),
			ErrorCode:     InternalError,
			EventAction:   "insert",
			ServiceAction: "saleValidated",
			UUID:          event.UUID,
		}
	}

	if m["items"] == nil {
		return &model.Document{
			AggregateID:   event.AggregateID,
			CorrelationID: event.CorrelationID,
			Error:         err.Error(),
			ErrorCode:     UserError,
			EventAction:   "insert",
			ServiceAction: "saleValidated",
			UUID:          event.UUID,
		}
	}

	items, assertOK := m["items"].([]interface{})
	if !assertOK {
		err = errors.New("error asserting Items to array")
		err = errors.Wrap(err, "SaleCreated-Event")
		log.Println(err)
		return &model.Document{
			AggregateID:   event.AggregateID,
			CorrelationID: event.CorrelationID,
			Error:         err.Error(),
			ErrorCode:     InternalError,
			EventAction:   "insert",
			ServiceAction: "saleValidated",
			UUID:          event.UUID,
		}
	}

	result := validateSaleItems(etcd, collection, items)

	marshalResult, err := json.Marshal(SaleValidationResp{
		OriginalRequest: m,
		Result:          result,
	})
	if err != nil {
		err = errors.Wrap(err, "SaleCreated-Event: Error marshalling result")
		log.Println(err)
		return &model.Document{
			AggregateID:   event.AggregateID,
			CorrelationID: event.CorrelationID,
			Error:         err.Error(),
			ErrorCode:     InternalError,
			EventAction:   "insert",
			ServiceAction: "saleValidated",
			UUID:          event.UUID,
		}
	}

	if producer == nil {
		kafkaBrokersStr := os.Getenv("KAFKA_BROKERS")
		producer, err = kafka.NewProducer(&kafka.ProducerConfig{
			KafkaBrokers: *commonutil.ParseHosts(kafkaBrokersStr),
		})
		if err != nil {
			err = errors.Wrap(err, "ValidateSale: Error creating producer")
			log.Println(err)
			return nil
		}
	}

	uuid, err := uuuid.NewV4()
	if err != nil {
		err = errors.Wrap(err, "UpdateSale: Error generating UUID for SaleValidation")
		log.Println(err)
		return nil
	}
	topic := os.Getenv("KAFKA_PRODUCER_EVENT_TOPIC")
	saleEvent, err := json.Marshal(model.Event{
		AggregateID:   3,
		CorrelationID: event.CorrelationID,
		Data:          marshalResult,
		EventAction:   "insert",
		NanoTime:      time.Now().UnixNano(),
		ServiceAction: "saleValidated",
		UUID:          uuid,
		YearBucket:    2018,
	})
	if err != nil {
		err = errors.Wrap(err, "CreateSale: Error Marshalling result")
		log.Println(err)
		return nil
	}

	producer.Input() <- kafka.CreateMessage(topic, saleEvent)

	return &model.Document{
		AggregateID:   event.AggregateID,
		CorrelationID: event.CorrelationID,
		EventAction:   "insert",
		Result:        marshalResult,
		ServiceAction: "saleValidated",
		UUID:          event.UUID,
	}
}

func validateSaleItems(
	etcd *clientv3.Client,
	collection *mongo.Collection,
	items []interface{},
) []SaleItemResult {
	result := []SaleItemResult{}

	for _, item := range items {
		itemMap, assertOK := item.(map[string]interface{})
		if !assertOK {
			err := errors.New("error asserting Item to Map")
			err = errors.Wrap(err, "SaleCreated-Event")
			log.Println(err)
			result = append(result, SaleItemResult{
				Error:     err.Error(),
				ErrorCode: InternalError,
			})
			continue
		}
		if itemMap["itemID"] == nil {
			err := errors.New("missing ItemID")
			err = errors.Wrap(err, "SaleCreated-Event")
			log.Println(err)
			result = append(result, SaleItemResult{
				Error:     err.Error(),
				ErrorCode: UserError,
			})
			continue
		}

		itemIDStr, assertOK := itemMap["itemID"].(string)
		if !assertOK {
			err := errors.New("error asserting ItemID to string")
			err = errors.Wrap(err, "SaleCreated-Event")
			log.Println(err)
			result = append(result, SaleItemResult{
				Error:     err.Error(),
				ErrorCode: UserError,
			})
		}
		itemID, err := uuuid.FromString(itemIDStr)
		if err != nil {
			err := errors.New("error parsing ItemID")
			err = errors.Wrap(err, "SaleCreated-Event")
			log.Println(err)
			result = append(result, SaleItemResult{
				Error:     err.Error(),
				ErrorCode: InternalError,
			})
		}

		if itemMap["weight"] == nil {
			err := errors.New("missing weight")
			err = errors.Wrap(err, "SaleCreated-Event")
			log.Println(err)
			result = append(result, SaleItemResult{
				ItemID:    itemID,
				Error:     err.Error(),
				ErrorCode: UserError,
			})
			continue
		}

		// Get this early as we remove it below to match document in Mongo
		soldWeight, err := commonutil.AssertFloat64(itemMap["weight"])
		if err != nil {
			err := errors.New("error asserting sold-item weight")
			err = errors.Wrap(err, "SaleCreated-Event")
			log.Println(err)
			result = append(result, SaleItemResult{
				ItemID:    itemID,
				Error:     err.Error(),
				ErrorCode: UserError,
			})
			continue
		}

		// So it can match document in Mongo
		delete(itemMap, "weight")

		lockTTL := concurrency.WithTTL(25)
		lockSession, err := concurrency.NewSession(etcd, lockTTL)
		if err != nil {
			err = errors.Wrapf(
				err,
				"SaleCreatedEvent: Failed to obtain lock for Updating ItemID: %s",
				itemIDStr,
			)
			log.Println(err)
			result = append(result, SaleItemResult{
				ItemID:    itemID,
				Error:     err.Error(),
				ErrorCode: InternalError,
			})
			continue
		}
		defer lockSession.Close()
		prefix := fmt.Sprintf("/%s/", itemIDStr)
		mx := concurrency.NewMutex(lockSession, prefix)

		lockCtx, lockCancel := context.WithTimeout(context.Background(), 25*time.Second)
		defer lockCancel()
		err = mx.Lock(lockCtx)

		unlockCtx, unlockCancel := context.WithTimeout(context.Background(), 25*time.Second)
		defer unlockCancel()
		defer mx.Unlock(unlockCtx)
		if err != nil {
			err = errors.Wrapf(
				err,
				"SaleCreatedEvent: Failed to apply obtained lock for ItemID: %s",
				itemIDStr,
			)
			log.Println(err)
			result = append(result, SaleItemResult{
				ItemID:    itemID,
				Error:     err.Error(),
				ErrorCode: InternalError,
			})
			lockSession.Close()
			continue
		}
		findArgs := map[string]interface{}{
			"itemID": itemIDStr,
		}
		findResult, err := collection.FindOne(findArgs)
		if err != nil {
			err := errors.Wrap(err, "SaleCreated-Event: Error getting Item from database")
			log.Println(err)
			result = append(result, SaleItemResult{
				ItemID:    itemID,
				Error:     err.Error(),
				ErrorCode: DatabaseError,
			})
			continue
		}

		inv, assertOK := findResult.(*Inventory)
		if !assertOK {
			err := errors.New("error asserting database-result to Inventory-Item")
			err = errors.Wrap(err, "SaleCreated-Event")
			log.Println(err)
			result = append(result, SaleItemResult{
				ItemID:    itemID,
				Error:     err.Error(),
				ErrorCode: InternalError,
			})
			continue
		}

		// Temporary fix for duplicated message-issue
		totalSoldWeight := inv.SoldWeight + (soldWeight / 2)
		cumWeight := totalSoldWeight + inv.WasteWeight + inv.DonateWeight
		if cumWeight > inv.TotalWeight {
			err := errors.New("sale-weight exceeds the total available weight")
			err = errors.Wrap(err, "SaleCreated-Event")
			log.Println(err)
			result = append(result, SaleItemResult{
				ItemID:    itemID,
				Error:     err.Error(),
				ErrorCode: InternalError,
			})
			continue
		}

		updateArgs := map[string]interface{}{
			"soldWeight": totalSoldWeight,
		}
		updateResult, err := collection.UpdateMany(findArgs, updateArgs)
		if err != nil {
			err = errors.Wrap(err, "SaleCreated-Event: Error writing new weight to database")
			log.Println(err)
			result = append(result, SaleItemResult{
				ItemID:    itemID,
				Error:     err.Error(),
				ErrorCode: DatabaseError,
			})
			continue
		}
		err = mx.Unlock(context.Background())
		if err != nil {
			err = errors.Wrapf(
				err,
				"SaleCreatedEvent: Failed to unlock ItemID: %s", itemIDStr,
			)
			log.Println(err)
		}
		err = lockSession.Close()
		if err != nil {
			err = errors.Wrapf(
				err,
				"SaleCreatedEvent: Failed to close lock-session for ItemID: %s", itemIDStr,
			)
			log.Println(err)
		}
		itemMap["weight"] = soldWeight

		if updateResult.ModifiedCount < 1 {
			err = errors.New("no items updated")
			err = errors.Wrap(err, "SaleCreated-Event")
			log.Println(err)
			result = append(result, SaleItemResult{
				ItemID:    itemID,
				Error:     err.Error(),
				ErrorCode: UserError,
			})
			continue
		}
		if updateResult.MatchedCount < 1 {
			err = errors.New("no items matched")
			err = errors.Wrap(err, "SaleCreated-Event")
			log.Println(err)
			result = append(result, SaleItemResult{
				ItemID:    itemID,
				Error:     err.Error(),
				ErrorCode: UserError,
			})
			continue
		}

		result = append(result, SaleItemResult{
			ItemID:          itemID,
			TotalSoldWeight: totalSoldWeight,
			TotalWeight:     inv.TotalWeight,
		})
	}

	return result
}
