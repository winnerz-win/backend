package mongo

import (
	"fmt"
	"jtools/dbg"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ObjectIDFromHex(hexID string) primitive.ObjectID {

	id, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		fmt.Println(err)
	}
	return id
}
func ObjectIdHex(hexID string) primitive.ObjectID {
	return ObjectIDFromHex(hexID)
}

func HexToObjectID(idString string) interface{} {
	return ObjectIDFromHex(idString)
}
func BSONID(hexID string) primitive.ObjectID {
	return ObjectIDFromHex(hexID)
}

func SelectorID(id string) Bson { return Bson{"_id": ObjectIDFromHex(id)} }

func IsObjectIdHex(id string) bool {
	return primitive.IsValidObjectID(id)
}

func InterfaceIdToObjectIDHex(id interface{}) string {
	switch v := id.(type) {
	case primitive.ObjectID:
		return v.Hex()
	case string:
		return v
	}
	return ""
}

/////////////////////////////////////////////////////////////////////////

type OBJECTID primitive.ObjectID

func (my OBJECTID) String() string { return primitive.ObjectID(my).Hex() }

/////////////////////////////////////////////////////////////////////////

// IDS : `bson:",inline" json:",inline"`   [ID]
type IDS struct {
	ID interface{} `bson:"_id" json:"_id"`
}

func (my IDS) Valid() bool    { return my.ID != nil }
func (my IDS) ValidIDS() bool { return my.Valid() }

func (my IDS) Selector() Bson { return Bson{"_id": my.ID} }

func (my IDS) IDString() string {
	return my.ID.(primitive.ObjectID).Hex()
}

func (my *IDS) MakeIDS() {
	if my.ID == nil {
		my.ID = primitive.NewObjectID()
	}
}

func NewIDS() IDS {
	return IDS{
		ID: primitive.NewObjectID(),
	}
}

// MapToBson :
func MapToBson(c map[string]interface{}) Bson {
	m := Bson{}
	for key, v := range c {
		if key == "_id" {
			switch val := v.(type) {
			case string:
				m[key] = ObjectIdHex(val)
			case primitive.ObjectID:
				m[key] = val
			default:
				m[key] = ObjectIdHex(dbg.Cat(val))
			}

		} else {
			m[key] = v
		}
	}
	return m
}
