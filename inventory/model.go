package inventory

import (
	"encoding/json"

	util "github.com/TerrexTech/go-commonutils/commonutil"

	"github.com/TerrexTech/uuuid"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/pkg/errors"
)

// Inventory defines the Inventory Aggregate.
type Inventory struct {
	ID           objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ItemID       uuuid.UUID        `bson:"itemID,omitempty" json:"itemID,omitempty"`
	Barcode      string            `bson:"barcode,omitempty" json:"barcode,omitempty"`
	DateArrived  int64             `bson:"dateArrived,omitempty" json:"dateArrived,omitempty"`
	DateSold     int64             `bson:"dateSold,omitempty" json:"dateSold,omitempty"`
	DeviceID     uuuid.UUID        `bson:"deviceID,omitempty" json:"deviceID,omitempty"`
	DonateWeight float64           `bson:"donateWeight,omitempty" json:"donateWeight,omitempty"`
	ExpiryDate   int64             `bson:"expiryDate,omitempty" json:"expiryDate,omitempty"`
	Lot          string            `bson:"lot,omitempty" json:"lot,omitempty"`
	Name         string            `bson:"name,omitempty" json:"name,omitempty"`
	Origin       string            `bson:"origin,omitempty" json:"origin,omitempty"`
	Price        float64           `bson:"price,omitempty" json:"price,omitempty"`
	Quantity     int64             `bson:"quantity,omitempty" json:"quantity,omitempty"`
	RSCustomerID uuuid.UUID        `bson:"rsCustomerID,omitempty" json:"rsCustomerID,omitempty"`
	SalePrice    float64           `bson:"salePrice,omitempty" json:"salePrice,omitempty"`
	SKU          string            `bson:"sku,omitempty" json:"sku,omitempty"`
	SoldWeight   float64           `bson:"soldWeight,omitempty" json:"soldWeight,omitempty"`
	Timestamp    int64             `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
	TotalWeight  float64           `bson:"totalWeight,omitempty" json:"totalWeight,omitempty"`
	UPC          int64             `bson:"upc,omitempty" json:"upc,omitempty"`
	WasteWeight  float64           `bson:"wasteWeight,omitempty" json:"wasteWeight,omitempty"`
}

// marshalInventory is simplified version of Inventory, for convenience
// in Marshalling and Unmarshalling operations.
type marshalInventory struct {
	ID           objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ItemID       string            `bson:"itemID,omitempty" json:"itemID,omitempty"`
	Barcode      string            `bson:"barcode,omitempty" json:"barcode,omitempty"`
	DateArrived  int64             `bson:"dateArrived,omitempty" json:"dateArrived,omitempty"`
	DateSold     int64             `bson:"dateSold,omitempty" json:"dateSold,omitempty"`
	DeviceID     string            `bson:"deviceID,omitempty" json:"deviceID,omitempty"`
	DonateWeight float64           `bson:"donateWeight,omitempty" json:"donateWeight,omitempty"`
	ExpiryDate   int64             `bson:"expiryDate,omitempty" json:"expiryDate,omitempty"`
	Lot          string            `bson:"lot,omitempty" json:"lot,omitempty"`
	Name         string            `bson:"name,omitempty" json:"name,omitempty"`
	Origin       string            `bson:"origin,omitempty" json:"origin,omitempty"`
	Price        float64           `bson:"price,omitempty" json:"price,omitempty"`
	Quantity     int64             `bson:"quantity,omitempty" json:"quantity,omitempty"`
	RSCustomerID string            `bson:"rsCustomerID,omitempty" json:"rsCustomerID,omitempty"`
	SalePrice    float64           `bson:"salePrice,omitempty" json:"salePrice,omitempty"`
	SKU          string            `bson:"sku,omitempty" json:"sku,omitempty"`
	SoldWeight   float64           `bson:"soldWeight,omitempty" json:"soldWeight,omitempty"`
	Timestamp    int64             `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
	TotalWeight  float64           `bson:"totalWeight,omitempty" json:"totalWeight,omitempty"`
	UPC          int64             `bson:"upc,omitempty" json:"upc,omitempty"`
	WasteWeight  float64           `bson:"wasteWeight,omitempty" json:"wasteWeight,omitempty"`
}

// MarshalBSON returns bytes of BSON-type.
func (i Inventory) MarshalBSON() ([]byte, error) {
	in := &marshalInventory{
		ID:           i.ID,
		ItemID:       i.ItemID.String(),
		Barcode:      i.Barcode,
		DateArrived:  i.DateArrived,
		DateSold:     i.DateSold,
		DeviceID:     i.DeviceID.String(),
		DonateWeight: i.DonateWeight,
		ExpiryDate:   i.ExpiryDate,
		Lot:          i.Lot,
		Name:         i.Name,
		Origin:       i.Origin,
		Price:        i.Price,
		Quantity:     i.Quantity,
		RSCustomerID: i.RSCustomerID.String(),
		SalePrice:    i.SalePrice,
		SKU:          i.SKU,
		SoldWeight:   i.SoldWeight,
		Timestamp:    i.Timestamp,
		TotalWeight:  i.TotalWeight,
		UPC:          i.UPC,
		WasteWeight:  i.WasteWeight,
	}

	return bson.Marshal(in)
}

// MarshalJSON returns bytes of JSON-type.
func (i *Inventory) MarshalJSON() ([]byte, error) {
	in := map[string]interface{}{
		"itemID":       i.ItemID.String(),
		"barcode":      i.Barcode,
		"dateArrived":  i.DateArrived,
		"dateSold":     i.DateSold,
		"deviceID":     i.DeviceID.String(),
		"donateWeight": i.DonateWeight,
		"expiryDate":   i.ExpiryDate,
		"lot":          i.Lot,
		"name":         i.Name,
		"origin":       i.Origin,
		"price":        i.Price,
		"quantity":     i.Quantity,
		"rsCustomerID": i.RSCustomerID.String(),
		"salePrice":    i.SalePrice,
		"sku":          i.SKU,
		"soldWeight":   i.SoldWeight,
		"timestamp":    i.Timestamp,
		"totalWeight":  i.TotalWeight,
		"upc":          i.UPC,
		"wasteWeight":  i.WasteWeight,
	}

	if i.ID != objectid.NilObjectID {
		in["_id"] = i.ID.Hex()
	}
	return json.Marshal(in)
}

// UnmarshalBSON returns BSON-type from bytes.
func (i *Inventory) UnmarshalBSON(in []byte) error {
	m := make(map[string]interface{})
	err := bson.Unmarshal(in, m)
	if err != nil {
		err = errors.Wrap(err, "Unmarshal Error")
		return err
	}

	err = i.unmarshalFromMap(m)
	return err
}

// UnmarshalJSON returns JSON-type from bytes.
func (i *Inventory) UnmarshalJSON(in []byte) error {
	m := make(map[string]interface{})
	err := json.Unmarshal(in, &m)
	if err != nil {
		err = errors.Wrap(err, "Unmarshal Error")
		return err
	}

	err = i.unmarshalFromMap(m)
	return err
}

// unmarshalFromMap unmarshals Map into Inventory.
func (i *Inventory) unmarshalFromMap(m map[string]interface{}) error {
	var err error
	var assertOK bool

	// Hoping to discover a better way to do this someday
	if m["_id"] != nil {
		i.ID, assertOK = m["_id"].(objectid.ObjectID)
		if !assertOK {
			i.ID, err = objectid.FromHex(m["_id"].(string))
			if err != nil {
				err = errors.Wrap(err, "Error while asserting ObjectID")
				return err
			}
		}
	}

	if m["itemID"] != nil {
		i.ItemID, err = uuuid.FromString(m["itemID"].(string))
		if err != nil {
			err = errors.Wrap(err, "Error while asserting ItemID")
			return err
		}
	}

	if m["deviceID"] != nil {
		i.DeviceID, err = uuuid.FromString(m["deviceID"].(string))
		if err != nil {
			err = errors.Wrap(err, "Error while asserting DeviceID")
			return err
		}
	}

	if m["rsCustomerID"] != nil {
		i.RSCustomerID, err = uuuid.FromString(m["rsCustomerID"].(string))
		if err != nil {
			err = errors.Wrap(err, "Error while asserting RSCustomerID")
			return err
		}
	}

	if m["barcode"] != nil {
		i.Barcode, assertOK = m["barcode"].(string)
		if !assertOK {
			return errors.New("Error while asserting Barcode")
		}
	}
	if m["dateArrived"] != nil {
		i.DateArrived, err = util.AssertInt64(m["dateArrived"])
		if err != nil {
			err = errors.Wrap(err, "Error while asserting DateArrived")
			return err
		}
	}
	if m["dateSold"] != nil {
		i.DateSold, err = util.AssertInt64(m["dateSold"])
		if err != nil {
			err = errors.Wrap(err, "Error while asserting DateSold")
			return err
		}
	}
	if m["donateWeight"] != nil {
		i.DonateWeight, err = util.AssertFloat64(m["donateWeight"])
		if err != nil {
			err = errors.Wrap(err, "Error while asserting DonateWeight")
			return err
		}
	}
	if m["expiryDate"] != nil {
		i.ExpiryDate, err = util.AssertInt64(m["expiryDate"])
		if err != nil {
			err = errors.Wrap(err, "Error while asserting ExpiryDate")
			return err
		}
	}
	if m["lot"] != nil {
		i.Lot, assertOK = m["lot"].(string)
		if !assertOK {
			return errors.New("Error while asserting Lot")
		}
	}
	if m["name"] != nil {
		i.Name, assertOK = m["name"].(string)
		if !assertOK {
			return errors.New("Error while asserting Name")
		}
	}
	if m["origin"] != nil {
		i.Origin, assertOK = m["origin"].(string)
		if !assertOK {
			return errors.New("Error while asserting Origin")
		}
	}
	if m["price"] != nil {
		i.Price, err = util.AssertFloat64(m["price"])
		if err != nil {
			err = errors.Wrap(err, "Error while asserting Price")
			return err
		}
	}
	if m["quantity"] != nil {
		i.Quantity, err = util.AssertInt64(m["quantity"])
		if err != nil {
			err = errors.Wrap(err, "Error while asserting Quantity")
			return err
		}
	}
	if m["salePrice"] != nil {
		i.SalePrice, err = util.AssertFloat64(m["salePrice"])
		if err != nil {
			err = errors.Wrap(err, "Error while asserting SalePrice")
			return err
		}
	}
	if m["sku"] != nil {
		i.SKU, assertOK = m["sku"].(string)
		if !assertOK {
			return errors.New("Error while asserting Sku")
		}
	}
	if m["soldWeight"] != nil {
		i.SoldWeight, err = util.AssertFloat64(m["soldWeight"])
		if err != nil {
			err = errors.Wrap(err, "Error while asserting SoldWeight")
			return err
		}
	}
	if m["timestamp"] != nil {
		i.Timestamp, err = util.AssertInt64(m["timestamp"])
		if err != nil {
			err = errors.Wrap(err, "Error while asserting Timestamp")
			return err
		}
	}
	if m["totalWeight"] != nil {
		i.TotalWeight, err = util.AssertFloat64(m["totalWeight"])
		if err != nil {
			err = errors.Wrap(err, "Error while asserting TotalWeight")
			return err
		}
	}
	if m["upc"] != nil {
		i.UPC, err = util.AssertInt64(m["upc"])
		if err != nil {
			err = errors.Wrap(err, "Error while asserting UPC")
			return err
		}
	}
	if m["wasteWeight"] != nil {
		i.WasteWeight, err = util.AssertFloat64(m["wasteWeight"])
		if err != nil {
			err = errors.Wrap(err, "Error while asserting WasteWeight")
			return err
		}
	}

	return nil
}
