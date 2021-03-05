package mongo

import (
	"context"
	"reflect"

	log "github.com/sirupsen/logrus"
	"github.com/skiprco/go-utils/v2/converters"
	"github.com/skiprco/go-utils/v2/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IMongoRepository interface {
	GetOne(ctx context.Context, collectionName string, query map[string]interface{}, acceptsEmptyResult bool, response interface{}, methodName string) *errors.GenericError
	GetMultiple(ctx context.Context, collectionName string, query map[string]interface{}, responses interface{}, methodName string, opts ...*GetMultipleOption) *errors.GenericError
	Save(ctx context.Context, collectionName string, entity interface{}, entityId interface{}, methodName string) *errors.GenericError
	Count(ctx context.Context, collectionName string, query map[string]interface{}, methodName string) (int64, *errors.GenericError)
	Delete(ctx context.Context, collectionName string, entityId string, methodName string) *errors.GenericError
}

type mongoRepository struct {
	client      *mongo.Client
	dbName      string
	collections map[string]*mongo.Collection
	domain      string
}

type GetMultipleOption struct {
	Sort  *GetMultipleOptionSort
	Limit int64
	Skip  int64
}

type GetMultipleOptionSort struct {
	FieldName string
	Ascending bool
}

func NewMongoRepository(ctx context.Context, mongoURL string, dbName string, collectionNames []string) (IMongoRepository, *errors.GenericError) {
	domain := "go-util"
	client, err := createClient(ctx, domain, mongoURL)
	if err != nil {
		return nil, err
	}
	database := client.Database(dbName)
	collections := map[string]*mongo.Collection{}
	for _, collectionName := range collectionNames {
		collections[collectionName] = database.Collection(collectionName)
	}

	return &mongoRepository{
		domain:      domain,
		client:      client,
		dbName:      dbName,
		collections: collections,
	}, nil
}

func createClient(ctx context.Context, domain, address string) (*mongo.Client, *errors.GenericError) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(address))
	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"address": address,
		}).Error("Failed to create mongo client")
		return nil, errors.NewGenericError(500, domain, "create_mongo_client", "error_connection", map[string]string{"error": err.Error()})
	}
	return client, nil
}

func (r *mongoRepository) getCollection(collectionName string) (*mongo.Collection, *errors.GenericError) {
	collection, isPresent := r.collections[collectionName]
	if !isPresent {
		return nil, errors.NewGenericError(500, r.domain, "mongoRepository", "can_t_find_collection", map[string]string{"collectionName": collectionName})
	}
	return collection, nil
}

// Save the entity. Try to find the entity by entityId on the field _id of the mongo collections.
// If there is a match, the entity is updated. If not, the entity is create
// the methodName parameter is used for logging / error
//
// Raises
//
// - 500/panic_during_sanitize_object: A panic occured during sanitation
//
// - 500/can_t_create_entity: Mongo library returned an error while doing an upsert
func (r *mongoRepository) Save(ctx context.Context, collectionName string, entity interface{}, entityId interface{}, methodName string) *errors.GenericError {
	// Sanitize entity
	var genErr *errors.GenericError = nil
	if reflect.TypeOf(entity).Kind() == reflect.Ptr {
		genErr = converters.SanitizeObject(entity)
	} else {
		genErr = converters.SanitizeObject(&entity)
	}
	if genErr != nil {
		return genErr
	}

	collection, genErr := r.getCollection(collectionName)
	if genErr != nil {
		return genErr
	}
	opts := options.Update().SetUpsert(true)
	value := bson.M{"$set": entity}
	query := bson.M{"_id": entityId}
	_, err := collection.UpdateOne(ctx, query, value, opts)
	if err != nil {
		log.WithFields(log.Fields{
			"error":       err,
			"method_name": methodName,
			"entity":      entity,
			"entity_id":   entityId,
		}).Error("can't create entity")
		return errors.NewGenericError(500, r.domain, methodName, "can_t_create_entity", nil)
	}
	return nil
}

// Count the number of entities found by the query
// the methodName parameter is used for logging / error
func (r *mongoRepository) Count(ctx context.Context, collectionName string, query map[string]interface{}, methodName string) (int64, *errors.GenericError) {
	collection, genErr := r.getCollection(collectionName)
	if genErr != nil {
		return 0, genErr
	}
	count, err := collection.CountDocuments(ctx, convertToBson(query))
	if err != nil {
		return 0, errors.NewGenericError(500, r.domain, methodName, "can_count_entities", nil)
	}
	return count, nil
}

// GetOne populates the response with the found entity based on the query
// the parameter acceptsEmptyResult defines the behavior when there is no result :
//		if acceptsEmptyResult == true, the function returns a nil error and response is not populate
//		if acceptsEmptyResult == false, the function returns an error
// the methodName parameter is used for logging / error
func (r *mongoRepository) GetOne(ctx context.Context, collectionName string, query map[string]interface{}, acceptsEmptyResult bool, response interface{}, methodName string) *errors.GenericError {
	collection, genErr := r.getCollection(collectionName)
	if genErr != nil {
		return genErr
	}
	result := collection.FindOne(ctx, convertToBson(query))

	if err := result.Err(); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			if acceptsEmptyResult {
				return nil
			}
			log.WithField(
				"error", err,
			).Warn("No entity found")
			return errors.NewGenericError(404, r.domain, methodName, "no_entity", nil)
		default:
			log.WithField(
				"error", err,
			).Error("can't fetch entity")
			return errors.NewGenericError(500, r.domain, methodName, "can_t_fetch_entity", nil)
		}
	}

	err := result.Decode(response)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"query": query,
		}).Error("Failed to decode response")
		return errors.NewGenericError(500, r.domain, methodName, "decode_error", nil)
	}
	return nil
}

// GetMultiple populates the response with all the entities found base on the query and the opts
// the responses must be a list of pointer (like []MyEntity ). Can be define like that : var responses []MyEntity
// opts contains a list of limits / sorting options
// the methodName parameter is used for logging / error
func (r *mongoRepository) GetMultiple(ctx context.Context, collectionName string, query map[string]interface{}, responses interface{}, methodName string, opts ...*GetMultipleOption) *errors.GenericError {
	collection, genErr := r.getCollection(collectionName)
	if genErr != nil {
		return genErr
	}
	mongoOpts := &options.FindOptions{}
	if len(opts) > 1 {
		return errors.NewGenericError(500, r.domain, methodName, "only_one_opts_take_in_care", nil)
	} else if len(opts) > 0 {
		opt := opts[0]
		if opt.Sort != nil {
			if opt.Sort.Ascending {
				mongoOpts.Sort = bson.D{primitive.E{Key: opt.Sort.FieldName, Value: 1}}
			} else {
				mongoOpts.Sort = bson.D{primitive.E{Key: opt.Sort.FieldName, Value: -1}}
			}
		}
		if opt.Limit != 0 {
			mongoOpts.Limit = &opt.Limit
		}
		if opt.Skip != 0 {
			mongoOpts.Skip = &opt.Skip
		}
	}

	cur, err := collection.Find(ctx, convertToBson(query), mongoOpts)
	if err != nil {
		return errors.NewGenericError(500, r.domain, methodName, "can_t_fetch_entity", nil)
	}
	err = cur.All(ctx, responses)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"query": query,
		}).Error("Failed to decode response")
		return errors.NewGenericError(500, r.domain, methodName, "decode_error", nil)
	}
	return nil
}

// Delete removes entity by entityId
// the methodName parameter is used for logging / error
func (r *mongoRepository) Delete(ctx context.Context, collectionName string, entityId string, methodName string) *errors.GenericError {
	collection, genErr := r.getCollection(collectionName)
	if genErr != nil {
		return genErr
	}
	// Call database
	pipeline := bson.M{"_id": entityId}
	_, err := collection.DeleteOne(ctx, pipeline)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"entityId": entityId,
		}).Error("Failed to delete entity")
		return errors.NewGenericError(500, r.domain, methodName, "delete_entity", nil)
	}
	return nil
}

func convertToBson(query map[string]interface{}) bson.M {
	bs := bson.M{}
	for k, v := range query {
		if v != nil && v != "" {
			bs[k] = v
		}
	}
	return bs
}
