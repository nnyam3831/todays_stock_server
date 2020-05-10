# Today`s Stock Server

Golang 웹 크롤링 + Rest API

## Stock Scrapper with Naver, Echo

## Project Stack

- Golang
- Echo
- Rest API
- GoQuery
- Heroku

## API


* /golden
  골든 크로스 데이터 파싱

- /rise
  상승 주식 파싱

* /kos
  코스피, 코스닥 추출

* /search
  인기 검색어 추출

## error

- JSON 형태로 return 할 때 struct 변수명이 Capital로 시작해야 인식할 수 있다.

- 한글 깨질 경우 utf-8인지 아닌지 확인할 것

- Heroku 배포 go.mod로 의존성 관리
