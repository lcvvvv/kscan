module kscan

go 1.18

require (
	github.com/atotto/clipboard v0.1.4
	//database
	github.com/denisenkom/go-mssqldb v0.10.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/huin/asn1ber v0.0.0-20120622192748-af09f62e6358
	github.com/icodeface/tls v0.0.0-20190904083142-17aec93c60e5
	github.com/jlaffaye/ftp v0.0.0-20220630165035-11536801d1ff
	github.com/lcvvvv/appfinger v0.1.1

	//gonmap
	github.com/lcvvvv/gonmap v1.3.4
	github.com/lcvvvv/pool v0.0.0-00010101000000-000000000000
	github.com/lcvvvv/simplehttp v0.1.1
	github.com/lcvvvv/stdio v0.1.2
	github.com/lib/pq v1.10.2

	//grdp
	github.com/lunixbochs/struc v0.0.0-20200707160740-784aaebc1d40
	github.com/miekg/dns v1.1.50
	github.com/sijms/go-ora/v2 v2.2.15

	//protocol
	github.com/stacktitan/smb v0.0.0-20190531122847-da9a425dceb8
	go.mongodb.org/mongo-driver v1.7.1
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5

	//chinese
	golang.org/x/text v0.3.7
)

require (
	github.com/PuerkitoBio/goquery v1.8.0 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/klauspost/compress v1.9.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/twmb/murmur3 v1.1.6 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.0.2 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/net v0.0.0-20210916014120-12bc252f5db8 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/tools v0.1.6-0.20210726203631-07bc1bf47fb2 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)

replace github.com/lcvvvv/pool => ./lib/pool

//replace github.com/lcvvvv/gonmap => ../go-github/gonmap
