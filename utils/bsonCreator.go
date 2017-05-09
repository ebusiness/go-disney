package utils

import (
	"gopkg.in/mgo.v2/bson"
)

// BsonCreator is a bson creator (without pointer)
type BsonCreator struct {
	Pipeline []bson.M
}

// Append the slice of bson.M
func (bc BsonCreator) Append(ms ...bson.M) BsonCreator {
	if bc.Pipeline == nil {
		bc.Pipeline = []bson.M{}
	}
	bc.Pipeline = append(bc.Pipeline, ms...)
	return bc
}

// Lookup without unwind will find a slice document (array)
func (bc BsonCreator) Lookup(from, localField, foreignField, as string) BsonCreator {
	return bc.Append(bson.M{
		"$lookup": bson.M{
			"from":         from,
			"localField":   localField,
			"foreignField": foreignField,
			"as":           as,
		},
	})
}

// LookupWithUnwind will find single document, not a slice
func (bc BsonCreator) LookupWithUnwind(from, localField, foreignField, as, lang string) BsonCreator {
	if len(lang) == 0 {
		return bc.Lookup(from, localField, foreignField, as).
			Append(bson.M{"$unwind": "$" + as})
	}

	return bc.Append(bson.M{"$addFields": bson.M{"old": "$$ROOT"}}).
		Lookup(from, localField, foreignField, as).
		Append(bson.M{"$unwind": "$" + as}).
		Append(bson.M{"$addFields": bson.M{"old." + as: "$" + as + "." + lang}}).
		Append(bson.M{"$replaceRoot": bson.M{"newRoot": "$old"}})
}

// GraphLookup Performs a recursive search on a collection,
// with options for restricting the search by recursion depth and query filter.
func (bc BsonCreator) GraphLookup(from, startWith, connectFromField, connectToField, as, lang string) BsonCreator {
	graphLookup := bson.M{
		"$graphLookup": bson.M{
			"from":             from,
			"startWith":        startWith,
			"connectFromField": connectFromField,
			"connectToField":   connectToField,
			"as":               as,
		},
	}
	if len(lang) == 0 {
		return bc.Append(graphLookup)
	}
	return bc.Append(bson.M{"$addFields": bson.M{"old": "$$ROOT"}}).
		Append(graphLookup).
		Append(bson.M{"$addFields": bson.M{"old." + as: "$" + as + "." + lang}}).
		Append(bson.M{"$replaceRoot": bson.M{"newRoot": "$old"}})
}
