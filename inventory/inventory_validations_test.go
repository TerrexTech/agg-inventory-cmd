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
	// var etcd *clientv3.Client

	Describe("delete", func() {
		It("should return error if filter is empty", func() {
			uuid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())
			cid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())
			uid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())

			mockEvent := &model.Event{
				EventAction:   "delete",
				CorrelationID: cid,
				AggregateID:   2,
				Data:          []byte("{}"),
				NanoTime:      time.Now().UnixNano(),
				UserUUID:      uid,
				UUID:          uuid,
				Version:       3,
				YearBucket:    2018,
			}
			kr := Delete(nil, mockEvent)
			Expect(kr.AggregateID).To(Equal(mockEvent.AggregateID))
			Expect(kr.CorrelationID).To(Equal(mockEvent.CorrelationID))
			Expect(kr.Error).ToNot(BeEmpty())
			Expect(kr.ErrorCode).To(Equal(int16(InternalError)))
			Expect(kr.UUID).To(Equal(mockEvent.UUID))
		})
	})

	Describe("insert", func() {
		It("should return error if itemID is empty", func() {
			uuid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())
			cid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())
			uid, err := uuuid.NewV4()
			Expect(err).ToNot(HaveOccurred())

			mockEvent := &model.Event{
				EventAction:   "insert",
				CorrelationID: cid,
				AggregateID:   2,
				Data:          []byte("{}"),
				NanoTime:      time.Now().UnixNano(),
				UserUUID:      uid,
				UUID:          uuid,
				Version:       3,
				YearBucket:    2018,
			}
			kr := Insert(nil, mockEvent)
			Expect(kr.AggregateID).To(Equal(mockEvent.AggregateID))
			Expect(kr.CorrelationID).To(Equal(mockEvent.CorrelationID))
			Expect(kr.Error).ToNot(BeEmpty())
			Expect(kr.ErrorCode).To(Equal(int16(InternalError)))
			Expect(kr.UUID).To(Equal(mockEvent.UUID))
		})
	})

	Describe("update", func() {
		It("should return error if filter is empty", func() {
			uuid, err := uuuid.NewV4()
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
				EventAction:   "delete",
				CorrelationID: cid,
				AggregateID:   2,
				Data:          marshalArgs,
				NanoTime:      time.Now().UnixNano(),
				UserUUID:      uid,
				UUID:          uuid,
				Version:       3,
				YearBucket:    2018,
			}
			kr := Update(nil, nil, mockEvent)
			Expect(kr.AggregateID).To(Equal(mockEvent.AggregateID))
			Expect(kr.CorrelationID).To(Equal(mockEvent.CorrelationID))
			Expect(kr.Error).ToNot(BeEmpty())
			Expect(kr.ErrorCode).To(Equal(int16(InternalError)))
			Expect(kr.UUID).To(Equal(mockEvent.UUID))
		})

		It("should return error if update is empty", func() {
			uuid, err := uuuid.NewV4()
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
				EventAction:   "delete",
				CorrelationID: cid,
				AggregateID:   2,
				Data:          marshalArgs,
				NanoTime:      time.Now().UnixNano(),
				UserUUID:      uid,
				UUID:          uuid,
				Version:       3,
				YearBucket:    2018,
			}
			kr := Update(nil, nil, mockEvent)
			Expect(kr.AggregateID).To(Equal(mockEvent.AggregateID))
			Expect(kr.CorrelationID).To(Equal(mockEvent.CorrelationID))
			Expect(kr.Error).ToNot(BeEmpty())
			Expect(kr.ErrorCode).To(Equal(int16(InternalError)))
			Expect(kr.UUID).To(Equal(mockEvent.UUID))
		})

		It("should return error if itemID in update is empty", func() {
			uuid, err := uuuid.NewV4()
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
				EventAction:   "delete",
				CorrelationID: cid,
				AggregateID:   2,
				Data:          marshalArgs,
				NanoTime:      time.Now().UnixNano(),
				UserUUID:      uid,
				UUID:          uuid,
				Version:       3,
				YearBucket:    2018,
			}
			kr := Update(nil, nil, mockEvent)
			Expect(kr.AggregateID).To(Equal(mockEvent.AggregateID))
			Expect(kr.CorrelationID).To(Equal(mockEvent.CorrelationID))
			Expect(kr.Error).ToNot(BeEmpty())
			Expect(kr.ErrorCode).To(Equal(int16(InternalError)))
			Expect(kr.UUID).To(Equal(mockEvent.UUID))
		})
	})
})
