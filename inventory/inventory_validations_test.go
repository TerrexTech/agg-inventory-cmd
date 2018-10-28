package inventory

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/TerrexTech/go-eventstore-models/model"
	"github.com/TerrexTech/uuuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TestInventory only tests basic pre-processing error-checks for Aggregate functions.
func TestInventory(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "InventoryAggregate Suite")
}

var _ = Describe("InventoryAggregate", func() {
	Describe("delete", func() {
		It("should return error if filter is empty", func() {
			timeUUID, err := uuuid.NewV1()
			Expect(err).ToNot(HaveOccurred())
			cid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())
			uid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())

			mockEvent := &model.Event{
				Action:        "delete",
				CorrelationID: cid,
				AggregateID:   2,
				Data:          []byte("{}"),
				Timestamp:     time.Now(),
				UserUUID:      uid,
				TimeUUID:      timeUUID,
				Version:       3,
				YearBucket:    2018,
			}
			kr := Delete(nil, mockEvent)
			Expect(kr.AggregateID).To(Equal(mockEvent.AggregateID))
			Expect(kr.CorrelationID).To(Equal(mockEvent.CorrelationID))
			Expect(kr.Error).ToNot(BeEmpty())
			Expect(kr.ErrorCode).To(Equal(int16(InternalError)))
			Expect(kr.UUID).To(Equal(mockEvent.TimeUUID))
		})
	})

	Describe("insert", func() {
		It("should return error if itemID is empty", func() {
			timeUUID, err := uuuid.NewV1()
			Expect(err).ToNot(HaveOccurred())
			cid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())
			uid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())

			mockEvent := &model.Event{
				Action:        "insert",
				CorrelationID: cid,
				AggregateID:   2,
				Data:          []byte("{}"),
				Timestamp:     time.Now(),
				UserUUID:      uid,
				TimeUUID:      timeUUID,
				Version:       3,
				YearBucket:    2018,
			}
			kr := Insert(nil, mockEvent)
			Expect(kr.AggregateID).To(Equal(mockEvent.AggregateID))
			Expect(kr.CorrelationID).To(Equal(mockEvent.CorrelationID))
			Expect(kr.Error).ToNot(BeEmpty())
			Expect(kr.ErrorCode).To(Equal(int16(InternalError)))
			Expect(kr.UUID).To(Equal(mockEvent.TimeUUID))
		})
	})

	Describe("update", func() {
		It("should return error if filter is empty", func() {
			timeUUID, err := uuuid.NewV1()
			Expect(err).ToNot(HaveOccurred())
			cid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())
			uid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())

			updateArgs := map[string]interface{}{
				"filter": map[string]interface{}{},
				"update": map[string]interface{}{},
			}
			marshalArgs, err := json.Marshal(updateArgs)
			Expect(err).ToNot(HaveOccurred())
			mockEvent := &model.Event{
				Action:        "delete",
				CorrelationID: cid,
				AggregateID:   2,
				Data:          marshalArgs,
				Timestamp:     time.Now(),
				UserUUID:      uid,
				TimeUUID:      timeUUID,
				Version:       3,
				YearBucket:    2018,
			}
			kr := Update(nil, mockEvent)
			Expect(kr.AggregateID).To(Equal(mockEvent.AggregateID))
			Expect(kr.CorrelationID).To(Equal(mockEvent.CorrelationID))
			Expect(kr.Error).ToNot(BeEmpty())
			Expect(kr.ErrorCode).To(Equal(int16(InternalError)))
			Expect(kr.UUID).To(Equal(mockEvent.TimeUUID))
		})

		It("should return error if update is empty", func() {
			timeUUID, err := uuuid.NewV1()
			Expect(err).ToNot(HaveOccurred())
			cid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())
			uid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())

			updateArgs := map[string]interface{}{
				"filter": map[string]interface{}{
					"x": 1,
				},
				"update": map[string]interface{}{},
			}
			marshalArgs, err := json.Marshal(updateArgs)
			Expect(err).ToNot(HaveOccurred())
			mockEvent := &model.Event{
				Action:        "delete",
				CorrelationID: cid,
				AggregateID:   2,
				Data:          marshalArgs,
				Timestamp:     time.Now(),
				UserUUID:      uid,
				TimeUUID:      timeUUID,
				Version:       3,
				YearBucket:    2018,
			}
			kr := Update(nil, mockEvent)
			Expect(kr.AggregateID).To(Equal(mockEvent.AggregateID))
			Expect(kr.CorrelationID).To(Equal(mockEvent.CorrelationID))
			Expect(kr.Error).ToNot(BeEmpty())
			Expect(kr.ErrorCode).To(Equal(int16(InternalError)))
			Expect(kr.UUID).To(Equal(mockEvent.TimeUUID))
		})

		It("should return error if itemID in update is empty", func() {
			timeUUID, err := uuuid.NewV1()
			Expect(err).ToNot(HaveOccurred())
			cid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())
			uid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())

			updateArgs := map[string]interface{}{
				"filter": map[string]interface{}{
					"x": 1,
				},
				"update": map[string]interface{}{
					"itemID": (uuuid.UUID{}).String(),
				},
			}
			marshalArgs, err := json.Marshal(updateArgs)
			Expect(err).ToNot(HaveOccurred())
			mockEvent := &model.Event{
				Action:        "delete",
				CorrelationID: cid,
				AggregateID:   2,
				Data:          marshalArgs,
				Timestamp:     time.Now(),
				UserUUID:      uid,
				TimeUUID:      timeUUID,
				Version:       3,
				YearBucket:    2018,
			}
			kr := Update(nil, mockEvent)
			Expect(kr.AggregateID).To(Equal(mockEvent.AggregateID))
			Expect(kr.CorrelationID).To(Equal(mockEvent.CorrelationID))
			Expect(kr.Error).ToNot(BeEmpty())
			Expect(kr.ErrorCode).To(Equal(int16(InternalError)))
			Expect(kr.UUID).To(Equal(mockEvent.TimeUUID))
		})
	})
})
