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
	"sync"
	"time"
)

type Repository struct {
	cfg        *config.Config
	mongo      *mongo.Client
	collection *mongo.Collection
}

func NewRepository() (*Repository, error) {
	client := config.ConnectMongoDB(config.MongoUri)
	collection := config.GetCollection(config.DB, config.COLLECTION)

	models := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "corp", Value: 1},
				{Key: "date", Value: -1},
			},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys: bson.D{
				{Key: "date", Value: -1},
			},
			Options: options.Index().SetUnique(false),
		},
	}
	_, err := collection.Indexes().CreateMany(context.TODO(), models)
	if err != nil {
		return nil, fmt.Errorf("failed to create index: %w", err)
	}

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

	oneMonthAgo := time.Now().AddDate(0, -1, -10)
	filter := bson.M{
		"date": bson.M{"$gte": oneMonthAgo},
	}

	if page < 1 {
		page = 1
	}
	findOptions := options.Find().
		SetSort(bson.D{{"date", -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))
	cursor, err := r.collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Printf("Error in Find: %v", err)
		return nil
	}
	defer cursor.Close(context.TODO())

	var results []types.Post
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Printf("Error decoding results: %v", err)
		return nil
	}
	return results
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
	}
	return result
}

func (r *Repository) GetLatestPost() time.Time {
	opts := options.FindOne().SetSort(bson.D{{"date", -1}}) // date 필드를 기준으로 내림차순 정렬
	var result types.Post
	if err := r.collection.FindOne(context.TODO(), bson.M{}, opts).Decode(&result); err != nil {
		log.Println("can't find latest post")
	}
	return result.Date
}

func (r *Repository) AutoSearchQuery(searchWord string) []types.Post {
	//if len(searchWord) > 5 { // 이 코드를 사용하면 한글이 안먹힘, 예를 들어 카프카는 ㅋ,ㅏ,ㅍ,ㅡ,ㅋ,ㅏ 로 인식하는가봄.
	//	searchWord = searchWord[:5]
	//}
	pipeline := mongo.Pipeline{
		{{"$search", bson.D{
			{"index", "title-auto-search"},
			{"autocomplete", bson.D{
				{"path", "title"},
				{"query", fmt.Sprintf("\"%s\"", searchWord)},
			}},
		}}},
		{{"$limit", 20}},
		bson.D{{Key: "$project", Value: bson.M{
			"title":   1,
			"date":    1,
			"url":     1,
			"corp":    1,
			"summary": 1,
		}}},
	}
	log.Println(searchWord)
	var results []types.Post
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Here Error occurs", err) // 왜 자동완성 후보를 클릭만 하면 터지지? pipeline 을 통해 보내는 단어길이도 제한했는데;;
		log.Println("pipeline", pipeline)
		return results
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	return results
}

func (r *Repository) SearchBlogsByCrop(corp string) []types.Post {

	page := 1
	pageSize := 30
	skip := (page - 1) * pageSize

	// 필터 조건 생성
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{
			"$and": []bson.M{
				{"corp": corp},
			},
		}}},

		// 정렬: 먼저 스코어로, 그 다음 날짜로
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "score", Value: -1},
			{Key: "date", Value: -1},
		}}},

		{{"$skip", skip}},
		{{"$limit", pageSize}},

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

	//var results []types.Post
	results := make([]types.Post, pageSize)
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	return results
}

func (r *Repository) UpdateDate() {

	// 날짜 포맷을 맞추기 위한 파싱 레이아웃 (예: "2024-04-30 00:00:00 +0000 UTC")
	layout := "2006-01-02 15:04:05 -0700 UTC"

	// 예시로 업데이트할 문서의 filter 조건
	filter := bson.M{"date": bson.M{"$type": "string"}} // date가 string인 문서들만 필터링

	// 카운트 (업데이트할 문서가 몇개인지 확인)
	count, err := r.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Fatalf("Failed to count documents: %v", err)
	}
	fmt.Printf("Found %d documents to update\n", count)

	// 커서로 업데이트할 모든 문서를 가져오기
	cursor, err := r.collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatalf("Failed to find documents: %v", err)
	}
	defer cursor.Close(context.TODO())

	// 업데이트 모델들을 저장할 채널
	updateChannel := make(chan mongo.WriteModel, count)

	// 동기화 용 WaitGroup
	var wg sync.WaitGroup

	// 한 번에 처리할 최대 goroutine 수 (이 값을 조정할 수 있습니다)
	maxGoroutines := 10
	sem := make(chan struct{}, maxGoroutines)

	// 각 문서를 업데이트하는 goroutine을 시작
	for cursor.Next(context.TODO()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatalf("Failed to decode document: %v", err)
		}

		// "date" 필드를 string에서 time.Time으로 변환
		dateStr, ok := result["date"].(string)
		if !ok {
			log.Printf("Error: date field is not a string for document %v", result["_id"])
			continue // 날짜 필드가 string이 아닌 경우 건너뛰기
		}

		// 날짜 문자열을 time.Time으로 변환
		date, err := time.Parse(layout, dateStr)
		if err != nil {
			log.Printf("Error parsing date string %s: %v", dateStr, err)
			continue // 날짜 파싱에 실패한 문서는 건너뛰기
		}

		// time.Time을 ISODate로 변환
		isoDate := primitive.NewDateTimeFromTime(date)

		// goroutine 실행
		wg.Add(1)
		go func(doc bson.M) {
			defer wg.Done()

			// 세마포어를 통해 goroutine 수를 제한
			sem <- struct{}{}
			defer func() { <-sem }()

			// 업데이트 모델 생성
			update := mongo.NewUpdateOneModel().
				SetFilter(bson.M{"_id": doc["_id"]}).
				SetUpdate(bson.M{
					"$set": bson.M{
						"date": isoDate, // 변환된 ISODate로 업데이트
					},
				})

			// 채널에 업데이트 모델 추가
			updateChannel <- update
		}(result)
	}

	// cursor에서 에러가 발생한 경우 처리
	if err := cursor.Err(); err != nil {
		log.Fatalf("Cursor error: %v", err)
	}

	// 모든 goroutine이 끝날 때까지 대기
	wg.Wait()

	// 채널에 저장된 업데이트 모델들을 bulkWrite로 한 번에 업데이트
	var models []mongo.WriteModel
	for update := range updateChannel {
		models = append(models, update)
	}

	// bulkWrite로 업데이트
	if len(models) > 0 {
		_, err = r.collection.BulkWrite(context.TODO(), models)
		if err != nil {
			log.Fatalf("Failed to bulk update documents: %v", err)
		}
	}
}
