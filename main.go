package main

import (
	. "blog-gopher/common/types"
	"blog-gopher/scrapper/banksalad"
	"blog-gopher/scrapper/daangn"
	"blog-gopher/scrapper/kakaopay"
	"blog-gopher/scrapper/oliveyoung"
	"blog-gopher/scrapper/toss"
	"log"
	"sync"
	"time"
)

//func main() {
//	// 엑셀 파일 생성
//	f := excelize.NewFile()
//
//	// 시트 이름 지정
//	sheetName := "PerformanceResults"
//	f.NewSheet(sheetName)
//
//	// 헤더 작성 (Iteration, Channel, Goroutine, Synchronous, Select)
//	headers := []string{"Iteration", "Channel", "Goroutine", "Synchronous", "Select"}
//	for i, header := range headers {
//		cell, _ := excelize.CoordinatesToCellName(i+1, 1) // A1, B1, C1, D1, E1
//		f.SetCellValue(sheetName, cell, header)
//	}
//
//	// 10회 반복 실행
//	for iteration := 1; iteration <= 10; iteration++ {
//		// Iteration 값 설정
//		cell, _ := excelize.CoordinatesToCellName(1, iteration+1) // A2, A3, ...
//		f.SetCellValue(sheetName, cell, iteration)
//
//		// Channel 방식
//		//start := time.Now()
//		//useChannel() // 실제 Channel 방식 함수 호출
//		//elapsed := time.Since(start)
//		//cell, _ = excelize.CoordinatesToCellName(2, iteration+1) // B2, B3, ...
//		//f.SetCellValue(sheetName, cell, fmt.Sprintf("%.3f", elapsed.Seconds()))
//
//		// Goroutine 방식
//		goroutineStart := time.Now()
//		callGoroutine() // 실제 Goroutine 방식 함수 호출
//		goroutineStartElapsed := time.Since(goroutineStart)
//		cell, _ = excelize.CoordinatesToCellName(3, iteration+1) // C2, C3, ...
//		f.SetCellValue(sheetName, cell, fmt.Sprintf("%.3f", goroutineStartElapsed.Seconds()))
//
//		// Synchronous 방식
//		syncStart := time.Now()
//		callSynchronous() // 실제 Synchronous 방식 함수 호출
//		syncElapsed := time.Since(syncStart)
//		cell, _ = excelize.CoordinatesToCellName(4, iteration+1) // D2, D3, ...
//		f.SetCellValue(sheetName, cell, fmt.Sprintf("%.3f", syncElapsed.Seconds()))
//
//		// Select 방식
//		selectStart := time.Now()
//		useSelect() // 실제 Select 방식 함수 호출
//		selectElapsed := time.Since(selectStart)
//		cell, _ = excelize.CoordinatesToCellName(5, iteration+1) // E2, E3, ...
//		f.SetCellValue(sheetName, cell, fmt.Sprintf("%.3f", selectElapsed.Seconds()))
//	}
//
//	// 엑셀 파일 저장
//	if err := f.SaveAs("performance_results.xlsx"); err != nil {
//		log.Fatal(err)
//	}
//
//	log.Println("엑셀 파일에 성능 결과가 저장되었습니다.")
//}

func main() {
	start := time.Now() // 시작 시간
	var result []Post
	result = callGoroutineChannel(result)
	result = callGoroutine(result)
	result = callSynchronous(result)
	log.Println("Total :", len(result))
	elapsed := time.Since(start) // 경과 시간
	log.Printf("[Goroutine] Execution Time: %s\n", elapsed)
}

func callGoroutine(result []Post) []Post {
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 회사별 스크래핑 함수들을 goroutine으로 실행
	scrapers := []func() []Post{
		kakaopay.CallApi,
		oliveyoung.CallApi,
		daangn.CallApi,
		toss.CallApi,
		banksalad.CallApi,
	}

	for _, scraper := range scrapers {
		wg.Add(1)
		go func(scrapeFunc func() []Post) {
			defer wg.Done()

			// 스크래핑 결과 가져오기
			posts := scrapeFunc()

			// 결과를 안전하게 result에 추가하기 위해 mutex 사용
			mu.Lock()
			result = append(result, posts...)
			mu.Unlock()
		}(scraper)
	}
	// 모든 goroutine이 완료될 때까지 대기
	wg.Wait()
	return result
}

func callGoroutineChannel(result []Post) []Post {
	scrapers := []func() []Post{
		kakaopay.CallApi,
		oliveyoung.CallApi,
		daangn.CallApi,
		toss.CallApi,
		banksalad.CallApi,
	}

	resultChan := make(chan []Post, len(scrapers))

	var wg sync.WaitGroup
	for _, scraper := range scrapers {
		wg.Add(1)
		go func(scrapeFunc func() []Post) {
			defer wg.Done()
			resultChan <- scrapeFunc()
		}(scraper)
	}

	// 결과 수집을 위한 고루틴
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for posts := range resultChan {
		result = append(result, posts...)
	}

	return result
}

//func callGoroutineChannelSelect(result []Post) []Post {
//
//	doneChan := make(chan struct{}) // 모든 작업이 끝났는지 체크하는 채널
//	scrapers := []func() []Post{
//		kakaopay.CallApi,
//		oliveyoung.CallApi,
//		daangn.CallApi,
//		toss.CallApi,
//		banksalad.CallApi,
//	}
//	resultChan := make(chan []Post, len(scrapers))
//
//	// 각 스크래핑 함수를 비동기적으로 실행
//	for _, scraper := range scrapers {
//		go func(scrapeFunc func() []Post) {
//			resultChan <- scrapeFunc()
//		}(scraper)
//	}
//
//	// 별도의 goroutine에서 결과를 수집
//	go func() {
//		for i := 0; i < len(scrapers); i++ {
//			select {
//			case posts := <-resultChan:
//				// 결과가 도착하면 비동기적으로 처리
//				result = append(result, posts...)
//			}
//		}
//		doneChan <- struct{}{} // 결과 수집 완료를 알림
//	}()
//
//	// 모든 작업이 완료될 때까지 대기 (doneChan으로 비동기 완료 확인)
//	<-doneChan
//	return result
//}

func callSynchronous(result []Post) []Post {
	result = append(result, kakaopay.CallApi()...)
	result = append(result, oliveyoung.CallApi()...)
	result = append(result, daangn.CallApi()...)
	result = append(result, toss.CallApi()...)
	result = append(result, banksalad.CallApi()...)
	return result
}
