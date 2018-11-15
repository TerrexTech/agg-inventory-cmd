package test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/TerrexTech/agg-inventory-cmd/inventory"
	"github.com/TerrexTech/go-commonutils/commonutil"
	"github.com/TerrexTech/go-eventstore-models/model"
	"github.com/TerrexTech/go-kafkautils/kafka"
	"github.com/TerrexTech/uuuid"
	"github.com/joho/godotenv"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

func Byf(s string, args ...interface{}) {
	By(fmt.Sprintf(s, args...))
}

func TestInventory(t *testing.T) {
	log.Println("Reading environment file")
	err := godotenv.Load("../.env")
	if err != nil {
		err = errors.Wrap(err,
			".env file not found, env-vars will be read as set in environment",
		)
		log.Println(err)
	}

	missingVar, err := commonutil.ValidateEnv(
		"KAFKA_BROKERS",
		"KAFKA_CONSUMER_EVENT_GROUP",

		"KAFKA_CONSUMER_EVENT_TOPIC",
		"KAFKA_CONSUMER_EVENT_QUERY_GROUP",
		"KAFKA_CONSUMER_EVENT_QUERY_TOPIC",

		"KAFKA_PRODUCER_EVENT_TOPIC",
		"KAFKA_PRODUCER_EVENT_QUERY_TOPIC",
		"KAFKA_PRODUCER_RESPONSE_TOPIC",

		"MONGO_HOSTS",
		"MONGO_USERNAME",
		"MONGO_PASSWORD",
		"MONGO_DATABASE",
		"MONGO_CONNECTION_TIMEOUT_MS",
		"MONGO_RESOURCE_TIMEOUT_MS",
	)

	if err != nil {
		err = errors.Wrapf(err, "Env-var %s is required for testing, but is not set", missingVar)
		log.Fatalln(err)
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "InventoryAggregate Suite")
}

var _ = Describe("InventoryAggregate", func() {
	var (
		kafkaBrokers          []string
		eventsTopic           string
		producerResponseTopic string
		eventRespTopic        string

		producer *kafka.Producer

		mockInv   *inventory.Inventory
		mockEvent *model.Event
	)
	BeforeSuite(func() {
		kafkaBrokers = *commonutil.ParseHosts(
			os.Getenv("KAFKA_BROKERS"),
		)
		eventsTopic = os.Getenv("KAFKA_PRODUCER_EVENT_TOPIC")
		producerResponseTopic = os.Getenv("KAFKA_PRODUCER_RESPONSE_TOPIC")
		eventRespTopic = os.Getenv("KAFKA_CONSUMER_EVENT_TOPIC")

		var err error
		producer, err = kafka.NewProducer(&kafka.ProducerConfig{
			KafkaBrokers: kafkaBrokers,
		})
		Expect(err).ToNot(HaveOccurred())

		itemID, err := uuuid.NewV4()
		Expect(err).ToNot(HaveOccurred())
		deviceID, err := uuuid.NewV4()
		Expect(err).ToNot(HaveOccurred())
		rsCustomerID, err := uuuid.NewV4()
		Expect(err).ToNot(HaveOccurred())

		mockInv = &inventory.Inventory{
			ItemID:       itemID,
			DateArrived:  time.Now().Unix(),
			DeviceID:     deviceID,
			Lot:          "test-lot",
			Name:         "test-name",
			Origin:       "test-origin",
			Price:        13.4,
			RSCustomerID: rsCustomerID,
			SalePrice:    12.23,
			SKU:          "test-sku",
			Timestamp:    time.Now().Unix(),
			TotalWeight:  300,
			UPC:          "test-upc",
			WasteWeight:  12,
		}
		marshalInv, err := json.Marshal(mockInv)
		Expect(err).ToNot(HaveOccurred())

		cid, err := uuuid.NewV4()
		Expect(err).ToNot(HaveOccurred())
		uid, err := uuuid.NewV4()
		Expect(err).ToNot(HaveOccurred())
		uuid, err := uuuid.NewV4()
		Expect(err).ToNot(HaveOccurred())
		mockEvent = &model.Event{
			EventAction:   "insert",
			CorrelationID: cid,
			AggregateID:   inventory.AggregateID,
			Data:          marshalInv,
			NanoTime:      time.Now().UnixNano(),
			UserUUID:      uid,
			UUID:          uuid,
			Version:       0,
			YearBucket:    2018,
		}
	})

	Describe("Inventory Operations", func() {
		It("should insert record", func(done Done) {
			Byf("Producing MockEvent")
			marshalEvent, err := json.Marshal(mockEvent)
			Expect(err).ToNot(HaveOccurred())
			producer.Input() <- kafka.CreateMessage(eventsTopic, marshalEvent)

			// Check if MockEvent was processed correctly
			Byf("Consuming Result")
			c, err := kafka.NewConsumer(&kafka.ConsumerConfig{
				KafkaBrokers: kafkaBrokers,
				GroupName:    "agginv.test.group.1",
				Topics:       []string{producerResponseTopic},
			})
			msgCallback := func(msg *sarama.ConsumerMessage) bool {
				defer GinkgoRecover()
				kr := &model.Document{}
				err := json.Unmarshal(msg.Value, kr)
				Expect(err).ToNot(HaveOccurred())

				if kr.UUID == mockEvent.UUID {
					Expect(kr.Error).To(BeEmpty())
					Expect(kr.ErrorCode).To(BeZero())
					Expect(kr.CorrelationID).To(Equal(mockEvent.CorrelationID))
					Expect(kr.UUID).To(Equal(mockEvent.UUID))

					inv := &inventory.Inventory{}
					err = json.Unmarshal(kr.Result, inv)
					Expect(err).ToNot(HaveOccurred())

					if inv.ItemID == mockInv.ItemID {
						mockInv.ID = inv.ID
						Expect(inv).To(Equal(mockInv))
						return true
					}
				}
				return false
			}

			handler := &msgHandler{msgCallback}
			c.Consume(context.Background(), handler)

			Byf("Checking if record got inserted into Database")
			aggColl, err := loadAggCollection()
			Expect(err).ToNot(HaveOccurred())
			findResult, err := aggColl.FindOne(mockInv)
			Expect(err).ToNot(HaveOccurred())
			findInv, assertOK := findResult.(*inventory.Inventory)
			Expect(assertOK).To(BeTrue())
			Expect(findInv).To(Equal(mockInv))

			close(done)
		}, 20)

		It("should validate new-sale and update sale-weight", func() {
			saleID, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())

			Byf("Creating mock create-sale Event")
			m := map[string]interface{}{
				"items": []map[string]interface{}{
					map[string]interface{}{
						"weight": 12.24,
						"upc":    "test-upc",
						"itemID": mockInv.ItemID,
						"lot":    "test-lot",
						"sku":    "test-sku",
					},
				},
				"saleID":    saleID,
				"timestamp": time.Now().UnixNano(),
			}

			cid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())
			uuid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())

			marshalMap, err := json.Marshal(m)
			Expect(err).ToNot(HaveOccurred())
			e := model.Event{
				AggregateID:   2,
				CorrelationID: cid,
				EventAction:   "update",
				ServiceAction: "createSale",
				Data:          marshalMap,
				NanoTime:      time.Now().UnixNano(),
				UUID:          uuid,
				Version:       0,
				YearBucket:    2018,
			}
			marshalEvent, err := json.Marshal(e)
			Expect(err).ToNot(HaveOccurred())
			msg := kafka.CreateMessage(eventsTopic, marshalEvent)
			Byf("Producing mock create-sale event")
			producer.Input() <- msg

			consTopic := fmt.Sprintf("%s", eventsTopic)
			consumer, err := kafka.NewConsumer(&kafka.ConsumerConfig{
				KafkaBrokers: kafkaBrokers,
				GroupName:    "test-group.1",
				Topics:       []string{consTopic},
			})
			var _ = eventRespTopic

			msgCallback := func(msg *sarama.ConsumerMessage) bool {
				defer GinkgoRecover()
				event := &model.Event{}
				err := json.Unmarshal(msg.Value, event)
				Expect(err).ToNot(HaveOccurred())

				if event.AggregateID == 3 && event.CorrelationID == cid {
					validationResp := &inventory.SaleValidationResp{}
					err = json.Unmarshal(event.Data, validationResp)
					Expect(err).ToNot(HaveOccurred())
					respItemID := validationResp.Result[0].ItemID
					Expect(respItemID).To(Equal(mockInv.ItemID))

					Expect(validationResp.Result[0].Error).To(BeEmpty())
					return true
				}
				return false
			}

			handler := &msgHandler{msgCallback}
			consumer.Consume(context.Background(), handler)
		})

		It("should update record", func(done Done) {
			Byf("Creating update args")
			filterInv := map[string]interface{}{
				"itemID": mockInv.ItemID,
			}
			mockInv.Origin = "new-origin"
			mockInv.Price = 500
			// Remove ObjectID because this is not passed from gateway
			mockID := mockInv.ID
			mockInv.ID = objectid.NilObjectID
			update := map[string]interface{}{
				"filter": filterInv,
				"update": mockInv,
			}
			marshalUpdate, err := json.Marshal(update)
			Expect(err).ToNot(HaveOccurred())
			// Reassign back ID so we can compare easily with database-entry
			mockInv.ID = mockID

			Byf("Creating update MockEvent")
			uuid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())
			mockEvent.EventAction = "update"
			mockEvent.Data = marshalUpdate
			mockEvent.NanoTime = time.Now().UnixNano()
			mockEvent.UUID = uuid

			Byf("Producing MockEvent")
			p, err := kafka.NewProducer(&kafka.ProducerConfig{
				KafkaBrokers: kafkaBrokers,
			})
			Expect(err).ToNot(HaveOccurred())
			marshalEvent, err := json.Marshal(mockEvent)
			Expect(err).ToNot(HaveOccurred())
			p.Input() <- kafka.CreateMessage(eventsTopic, marshalEvent)

			// Check if MockEvent was processed correctly
			Byf("Consuming Result")
			c, err := kafka.NewConsumer(&kafka.ConsumerConfig{
				KafkaBrokers: kafkaBrokers,
				GroupName:    "agginv.test.group.1",
				Topics:       []string{producerResponseTopic},
			})
			msgCallback := func(msg *sarama.ConsumerMessage) bool {
				defer GinkgoRecover()
				kr := &model.Document{}
				err := json.Unmarshal(msg.Value, kr)
				Expect(err).ToNot(HaveOccurred())

				if kr.UUID == mockEvent.UUID {
					Expect(kr.Error).To(BeEmpty())
					Expect(kr.ErrorCode).To(BeZero())
					Expect(kr.CorrelationID).To(Equal(mockEvent.CorrelationID))
					Expect(kr.UUID).To(Equal(mockEvent.UUID))

					result := map[string]int{}
					err = json.Unmarshal(kr.Result, &result)
					Expect(err).ToNot(HaveOccurred())

					if result["matchedCount"] != 0 && result["modifiedCount"] != 0 {
						Expect(result["matchedCount"]).To(Equal(1))
						Expect(result["modifiedCount"]).To(Equal(1))
						return true
					}
				}
				return false
			}

			handler := &msgHandler{msgCallback}
			c.Consume(context.Background(), handler)

			Byf("Checking if record got inserted into Database")
			aggColl, err := loadAggCollection()
			Expect(err).ToNot(HaveOccurred())
			findResult, err := aggColl.FindOne(mockInv)
			Expect(err).ToNot(HaveOccurred())
			findInv, assertOK := findResult.(*inventory.Inventory)
			Expect(assertOK).To(BeTrue())
			Expect(findInv).To(Equal(mockInv))

			close(done)
		}, 20)

		It("should delete record", func(done Done) {
			Byf("Creating delete args")
			deleteArgs := map[string]interface{}{
				"itemID": mockInv.ItemID,
			}
			marshalDelete, err := json.Marshal(deleteArgs)
			Expect(err).ToNot(HaveOccurred())

			Byf("Creating delete MockEvent")
			uuid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())
			mockEvent.EventAction = "delete"
			mockEvent.Data = marshalDelete
			mockEvent.NanoTime = time.Now().UnixNano()
			mockEvent.UUID = uuid

			Byf("Producing MockEvent")
			p, err := kafka.NewProducer(&kafka.ProducerConfig{
				KafkaBrokers: kafkaBrokers,
			})
			Expect(err).ToNot(HaveOccurred())
			marshalEvent, err := json.Marshal(mockEvent)
			Expect(err).ToNot(HaveOccurred())
			p.Input() <- kafka.CreateMessage(eventsTopic, marshalEvent)

			// Check if MockEvent was processed correctly
			Byf("Consuming Result")
			c, err := kafka.NewConsumer(&kafka.ConsumerConfig{
				KafkaBrokers: kafkaBrokers,
				GroupName:    "agginv.test.group.1",
				Topics:       []string{producerResponseTopic},
			})
			msgCallback := func(msg *sarama.ConsumerMessage) bool {
				defer GinkgoRecover()
				kr := &model.Document{}
				err := json.Unmarshal(msg.Value, kr)
				Expect(err).ToNot(HaveOccurred())

				if kr.UUID == mockEvent.UUID {
					Expect(kr.Error).To(BeEmpty())
					Expect(kr.ErrorCode).To(BeZero())
					Expect(kr.CorrelationID).To(Equal(mockEvent.CorrelationID))
					Expect(kr.UUID).To(Equal(mockEvent.UUID))

					result := map[string]int{}
					err = json.Unmarshal(kr.Result, &result)
					Expect(err).ToNot(HaveOccurred())

					if result["deletedCount"] != 0 {
						Expect(result["deletedCount"]).To(Equal(1))
						return true
					}
				}
				return false
			}

			handler := &msgHandler{msgCallback}
			c.Consume(context.Background(), handler)

			Byf("Checking if record got inserted into Database")
			aggColl, err := loadAggCollection()
			Expect(err).ToNot(HaveOccurred())
			_, err = aggColl.FindOne(mockInv)
			Expect(err).To(HaveOccurred())

			close(done)
		}, 20)
	})
})
