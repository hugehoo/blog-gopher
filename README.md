# blog-gopher

### 개요
- 다양한 회사에서 자체적으로 운영하는 기술블로그를 모아볼 수 있는 서비스를 제공합니다.

### Service Architecture
![img.png](img.png)

### Tech Stacks
- Go v1.21
- Gin v1.10
- Goquery v1.9.2
- MongoDB
- AWS Lambda

### API
- [X] GET `/posts` : 블로그 전체 조회
- [X] GET `/search` : 블로그 검색

### Scrapping Blogs
- [X] 당근
- [X] 토스
- [X] 버즈빌
- [X] 무신사
- [X] 29cm
- [X] 올리브영
- [X] 카카오페이
- [X] 컬리
- [X] 카카오뱅크
- [X] 데브시스터즈
- [X] 오늘의집
- [X] 라인
- [X] 쏘카
- [ ] 화해
- [ ] 크몽
- [ ] AWS KR
- [ ] 리디

### Main Feature
- [X] Tech blog scrapper
- [X] Categorizing by corp
- [ ] Check read article(blog)
- [ ] Bookmark article

### Todo 
- [X] goroutine
- [X] Date Parsing
- [X] repository
- [X] search
- [X] web framework
- [X] default toggle
- [ ] Test code
- [ ] scheduler per day
  - 크론잡이나 별도 스케줄러 추가 -> 매일 새벽 세시
  - 배치 돌리기 전에 헬스체크 3번 정도 날리자(10초 주기)
  - gocron 라이브러리 사용
  - 기존 post-scrapper 사용
  - 스크래핑 시 기존에 저장된 값이 나오면 break 사용해서 해당 고루틴 탈출
  - 기존 모든 블로그 스크래핑 하는 함수와 별개로 partial update 할 함수를 새로 만들어야겠다.
  - 우선 몽고에서 각 회사별 가장 최신의 블로그 가져와 날짜 확인. 
  - 스크래핑 기준 날짜와 최신 날짜 비교하여 해당 날짜 사이의 글만 긁기. -> 근데 이건 동일한 날짜를 못가져올 수도.
  - 

### 어떤 문제를 해결하려고 하는건지, 혹은 어떤 기술적 어려움이 있는지.
- 읽은 블로그를 체크하고 싶다. -> 인증/인가 필요 (github)
- 자주 보는 아티클을 별도로 저장하고 싶다.
- 블로그 스크래핑을 10초 내로 완료하고 싶다.
- 나만 사용할 수 있는 기능 -> 블로그 저장/읽음 표시
- 텍스트 전문 검색 기능

### 이슈 -> 몽고디비 기본조회가 느림
- date 필드 문자열로 돼있어 이슈이슈 핫이슈 -> 문자열을 모두 그걸로 date type 으로 바꿈
- VPC peering -> 이거 시도해봐야징.
