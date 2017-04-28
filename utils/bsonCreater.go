package utils

import (
	"gopkg.in/mgo.v2/bson"
)

// BsonCreater is a bson creater (without pointer)
type BsonCreater struct {
	Pipeline []bson.M
}

// Append the slice of bson.M
func (bc BsonCreater) Append(ms ...bson.M) BsonCreater {
	if bc.Pipeline == nil {
		bc.Pipeline = []bson.M{}
	}
	bc.Pipeline = append(bc.Pipeline, ms...)
	return bc
}

// Lookup without unwind will find a slice document (array)
func (bc BsonCreater) Lookup(from, localField, foreignField, as string) BsonCreater {
	return bc.lookup(from, localField, foreignField, as, false)
}

// LookupWithUnwind will find single document, not a slice
func (bc BsonCreater) LookupWithUnwind(from, localField, foreignField, as string) BsonCreater {
	return bc.lookup(from, localField, foreignField, as, true)
}

// lookup Performs a left outer join to an unsharded collection in the same database to filter in documents from the “joined” collection for processing.
func (bc BsonCreater) lookup(from, localField, foreignField, as string, withUnwind bool) BsonCreater {
	lookup := bson.M{
		"$lookup": bson.M{
			"from":         from,
			"localField":   localField,
			"foreignField": foreignField,
			"as":           as,
		},
	}

	if withUnwind {
		return bc.Append(lookup, bson.M{"$unwind": "$" + as})
	}
	return bc.Append(lookup)
}

// GraphLookup Performs a recursive search on a collection,
// with options for restricting the search by recursion depth and query filter.
func (bc BsonCreater) GraphLookup(from, startWith, connectFromField, connectToField, as string) BsonCreater {
	graphLookup := bson.M{
		"$graphLookup": bson.M{
			"from":             from,
			"startWith":        startWith,
			"connectFromField": connectFromField,
			"connectToField":   connectToField,
			"as":               as,
		},
	}
	return bc.Append(graphLookup)
}
