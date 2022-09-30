module github.com/lcvvvv/appfinger

go 1.17

replace github.com/lcvvvv/simplehttp => ./../simplehttp

replace github.com/lcvvvv/stdio => ./../stdio

require (
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/lcvvvv/simplehttp v0.0.0-00010101000000-000000000000
	github.com/lcvvvv/stdio v0.0.0-00010101000000-000000000000
	github.com/twmb/murmur3 v1.1.6
)

require (
	github.com/andybalholm/cascadia v1.3.1 // indirect
	golang.org/x/net v0.0.0-20210916014120-12bc252f5db8 // indirect
	golang.org/x/text v0.3.7 // indirect
)
