package repository

import (
	"blog-gopher/common/types"
	"blog-gopher/config"
	"context"
	"go.mongodb.org/mongo-driver/bson"
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
		doc := bson.D{
			{"title", post.Title},
			{"url", post.Url},
			{"summary", post.Summary},
			{"date", post.Date},
			{"corp", post.Corp},
		}
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
	//.         // 날짜를 기준으로 내림차순 정렬
	//	SetSkip(int64((page - 1) * pageSize)). // 페이지 건너뛰기
	//	SetLimit(int64(pageSize))              // 페이지 크기 설정

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

func (r *Repository) SearchBlogs(searchWord string, page int, size int) []types.Post {

	//filter := bson.M{
	//	"$text": bson.M{
	//		"$search": searchWord,
	//	},
	//}
	pagingOptions := options.Find().
		SetSort(bson.D{{"date", -1}})
	// title 필드에서 검색하는 쿼리 생성 - 정규 표현식 사용
	filter := bson.M{
		"title": bson.M{
			"$regex":   searchWord,
			"$options": "i", // 대소문자 구분 없음
		},
	}
	cursor, err := r.collection.Find(context.TODO(), filter, pagingOptions)
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
