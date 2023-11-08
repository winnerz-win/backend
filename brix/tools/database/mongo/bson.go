package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Bson bson.M

type Dson bson.D

func (my *Dson) Append(e Eson) {
	*my = append(*my, e.E())
}

type Eson bson.E

func (my Eson) E() bson.E { return bson.E(my) }
