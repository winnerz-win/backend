package mongo

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StringToDocument : mongo.VOID , mongo.MAP
func StringToDocument(void any) error {
	var item MAP
	switch v := void.(type) {
	case VOID:
		item = MAP(v)
	case MAP:
		item = v
	case map[string]interface{}:
		item = MAP(v)

	default:
		return errors.New("[StringToDocument] Type Error")
	}

	var re_err error
	item.Get("_id", func(val interface{}) {
		switch id := val.(type) {
		case string:
			_id, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				re_err = err
				return
			}
			item["_id"] = _id

		case primitive.ObjectID:

		default:
			re_err = errors.New("Invalid _id format.")
		} //switch
	})
	_string_to_datetime(item)
	_string_to_decimal128(item)

	return re_err
}

func DocumentToString(item MAP, isContainsID ...bool) MAP {
	if !isTrue(isContainsID) {
		delete(item, "_id")
	}
	_datetime_to_string(item)
	_decimal128_to_string(item)
	return item
}
