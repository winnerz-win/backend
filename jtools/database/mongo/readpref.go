package mongo

import "go.mongodb.org/mongo-driver/mongo/readpref"

func _Primary() *readpref.ReadPref {
	return readpref.Primary()
}

type ReadPref readpref.Mode

const (
	Primary            = ReadPref(readpref.PrimaryMode)
	PrimaryPreferred   = ReadPref(readpref.PrimaryPreferredMode)
	Secondary          = ReadPref(readpref.SecondaryMode)
	SecondaryPreferred = ReadPref(readpref.SecondaryPreferredMode)
	Nearest            = ReadPref(readpref.NearestMode)
)

func (my ReadPref) Mode() *readpref.ReadPref {

	switch my {
	case Primary:
		return readpref.Primary()
	case PrimaryPreferred:
		return readpref.PrimaryPreferred()

	case Secondary:
		return readpref.Secondary()
	case SecondaryPreferred:
		return readpref.SecondaryPreferred()

	case Nearest:
		return readpref.Nearest()
	}

	return readpref.Primary()
}
