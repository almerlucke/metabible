package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (store *Store) createBooksIndex(ctx context.Context) error {
	col := store.Collection("books")

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys: bson.D{
			{"book_name", 1},
		},
		Options: options.Index().SetUnique(true),
	}, {
		Keys: bson.D{
			{"book_div", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"testament", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"book_order", 1},
		},
		Options: nil,
	}})

	return err
}

func (store *Store) createChaptersIndex(ctx context.Context) error {
	col := store.Collection("chapters")

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys: bson.D{
			{"book", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"chapter_num", 1},
		},
		Options: nil,
	}})

	return err
}

func (store *Store) createEastonIndex(ctx context.Context) error {
	col := store.Collection("easton")

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys: bson.D{
			{"dict_lookup", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"dict_text", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"match_slugs", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"match_type", 1},
		},
		Options: nil,
	}})

	return err
}

func (store *Store) createEventsIndex(ctx context.Context) error {
	col := store.Collection("events")

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys: bson.D{
			{"title", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"sort_key", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"start_date", 1},
		},
		Options: nil,
	}})

	return err
}

func (store *Store) createPeopleIndex(ctx context.Context) error {
	col := store.Collection("people")

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys: bson.D{
			{"dict_text", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"dictionary_text", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"display_title", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"gender", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"max_year", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"min_year", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"name", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"person_lookup", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"slug", 1},
		},
		Options: nil,
	}})

	return err
}

func (store *Store) createPeopleGroupsIndex(ctx context.Context) error {
	col := store.Collection("people_groups")

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys: bson.D{
			{"group_name", 1},
		},
		Options: nil,
	}})

	return err
}

func (store *Store) createPeriodsIndex(ctx context.Context) error {
	col := store.Collection("periods")

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys: bson.D{
			{"bc_ad", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"formatted_year", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"iso_year", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"year_num", 1},
		},
		Options: nil,
	}})

	return err
}

func (store *Store) createPlacesIndex(ctx context.Context) error {
	col := store.Collection("places")

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys: bson.D{
			{"aliases", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"comment", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"dict_text", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"display_title", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"esv_name", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"feature_sub_type", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"feature_type", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"kjv_name", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"latitude", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"longitude", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"place_lookup", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"recogito_label", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"slug", 1},
		},
		Options: nil,
	}})

	return err
}

func (store *Store) createVersesIndex(ctx context.Context) error {
	col := store.Collection("verses")

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys: bson.D{
			{"md_text", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"rich_text", 1},
		},
		Options: nil,
	}, {
		Keys: bson.D{
			{"verse_text", 1},
		},
		Options: nil,
	}})

	return err
}

func (store *Store) createIndexes(ctx context.Context) error {
	for _, f := range []func(context.Context) error{
		store.createBooksIndex,
		store.createChaptersIndex,
		store.createEastonIndex,
		store.createEventsIndex,
		store.createPeopleIndex,
		store.createPeopleGroupsIndex,
		store.createPeriodsIndex,
		store.createPlacesIndex,
		store.createVersesIndex,
	} {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
