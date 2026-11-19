package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type FakeCollection struct {
	EstimatedDocumentCountFn func(ctx context.Context, opts ...options.Lister[options.EstimatedDocumentCountOptions]) (int64, error)
	CountDocumentsFn         func(ctx context.Context, filter any, opts ...options.Lister[options.CountOptions]) (int64, error)
	InsertOneFn              func(ctx context.Context, doc any, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error)
	FindFn                   func(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error)
	FindOneFn                func(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult
	UpdateOneFn              func(ctx context.Context, filter any, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error)
	DeleteOneFn              func(ctx context.Context, filter any, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error)
}

func (f *FakeCollection) EstimatedDocumentCount(ctx context.Context, opts ...options.Lister[options.EstimatedDocumentCountOptions]) (int64, error) {
	return f.EstimatedDocumentCountFn(ctx, opts...)
}

func (f *FakeCollection) CountDocuments(ctx context.Context, filter any, opts ...options.Lister[options.CountOptions]) (int64, error) {
	return f.CountDocumentsFn(ctx, filter, opts...)
}

func (f *FakeCollection) InsertOne(ctx context.Context, doc any, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error) {
	return f.InsertOneFn(ctx, doc, opts...)
}

func (f *FakeCollection) Find(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error) {
	return f.FindFn(ctx, filter, opts...)
}

func (f *FakeCollection) FindOne(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
	return f.FindOneFn(ctx, filter, opts...)
}

func (f *FakeCollection) UpdateOne(ctx context.Context, filter any, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
	return f.UpdateOneFn(ctx, filter, update, opts...)
}

func (f *FakeCollection) DeleteOne(ctx context.Context, filter any, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error) {
	return f.DeleteOneFn(ctx, filter, opts...)
}
