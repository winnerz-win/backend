package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type Iterator interface {
	RemainingBatchLength() int
	Close() error
	Next(r interface{}) bool

	One(r interface{}) error
	All(slicePointer interface{}) error

	Error() error
}

type cIterator struct {
	ctx    context.Context
	err    error
	cursor *mongo.Cursor
}

func (my cIterator) Error() error { return my.err }

func (my *cIterator) RemainingBatchLength() int {
	if my.err != nil {
		return 0
	}
	return my.cursor.RemainingBatchLength()
}

func (my *cIterator) Close() error {
	if my.err != nil {
		return my.err
	}
	return my.cursor.Close(my.ctx)
}

func (my *cIterator) Next(r interface{}) bool {
	if my.err != nil {
		return false
	}
	next := my.cursor.Next(my.ctx)
	if next {
		my.cursor.Decode(r)
	}
	if !next {
		my.cursor.Close(my.ctx)
	}
	return next
}

func (my *cIterator) One(r interface{}) error {
	if my.err != nil {
		return my.err
	}
	defer func() {
		my.cursor.Close(my.ctx)
	}()
	if my.cursor.Next(my.ctx) {
		return my.cursor.Decode(r)
	}
	return errors.New("empty one data")
}

func (my *cIterator) All(slicePointer interface{}) error {
	if my.err != nil {
		return my.err
	}
	defer func() {
		my.cursor.Close(my.ctx)
	}()
	return my.cursor.All(
		my.ctx,
		slicePointer,
	)
}
