package repository

import (
	"blog-gopher/common/dto"
	"blog-gopher/common/types"
	"blog-gopher/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Repository struct {
	cfg        *config.Config
	mongo      *mongo.Client
	collection *mongo.Collection
}

func NewRepository() (*Repository, error) {
	client := config.ConnectMongoDB(config.MongoUri)
	collection := config.GetCollection(config.DB, config.COLLECTION)
	return &Repository{mongo: client, collection: collection},
		nil
}

func (r *Repository) InsertBlogs(posts []types.Post) {
	var documents []interface{}
	for _, post := range posts {
		doc := dto.InsertPost(post)
		documents = append(documents, doc)
	}
	if len(posts) == 0 {
		return
	}

	if _, err := r.collection.InsertMany(context.TODO(), documents); err != nil {
		log.Println("Insert fail", err)
		panic(err)
	}
}

func (r *Repository) FindBlogs(corps []string, page int, pageSize int) []types.Post {
	// 데이터 정렬 및 페이징
	pagingOptions := options.Find().
		SetSort(bson.D{{"date", -1}})

	var filter bson.M
	if len(corps) > 0 { // corps 배열에 값이 있을 때만 $in 필터 적용
		filter = bson.M{
			"corp": bson.M{"$in": corps},
		}
	} else {
		filter = bson.M{}
	}

	cursor, err := r.collection.Find(context.TODO(), filter, pagingOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	var results []types.Post
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
		return nil
	} else {
		return results
	}

}

func (r *Repository) DeleteAll() {
	//r.collection.Drop()

}

func (r *Repository) SearchBlogs(searchWord string) []types.Post {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{
			"$text": bson.M{"$search": fmt.Sprintf("\"%s\"", searchWord)},
		}}},
		// 텍스트 검색 스코어 추가
		bson.D{{Key: "$addFields", Value: bson.M{
			"score": bson.M{"$meta": "textScore"},
		}}},
		bson.D{{Key: "$match", Value: bson.M{
			"score": bson.M{"$gte": 0.75},
		}}},
		// 정렬: 먼저 스코어로, 그 다음 날짜로
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "score", Value: -1},
			{Key: "date", Value: -1},
		}}},
		// 결과 프로젝션
		bson.D{{Key: "$project", Value: bson.M{
			"title":   1,
			"date":    1,
			"url":     1,
			"corp":    1,
			"summary": 1,
		}}},
	}
	cursor, err := r.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	var results []types.Post
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	return results
}

func (r *Repository) UpdatePost(postID string, text string) (interface{}, interface{}) {
	objectID, err := primitive.ObjectIDFromHex(postID) // objectIDFromHex 이걸 해줬어야 했는데, 이거 없이 걍 string 으로 필터치려고 하니 업데이트가 안됐다.
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"content": text}}
	result, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %v", err)
	}
	return result, nil
}

func (r *Repository) SearchPostById(id string) types.Post {
	objectID, _ := primitive.ObjectIDFromHex(id)
	var result types.Post
	if err := r.collection.FindOne(context.TODO(), (bson.M{"_id": objectID})).Decode(&result); err != nil {
		log.Println("Can't find Post")
	} else {
		if err := r.mongo.Disconnect(context.TODO()); err != nil {
			log.Println("Error when disconnect")
		}
	}
	return result
}
