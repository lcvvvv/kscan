module kscan

go 1.16

require (
	github.com/atotto/clipboard v0.1.4
	//database
	github.com/denisenkom/go-mssqldb v0.10.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/huin/asn1ber v0.0.0-20120622192748-af09f62e6358
	github.com/icodeface/tls v0.0.0-20190904083142-17aec93c60e5
	github.com/jlaffaye/ftp v0.0.0-20220630165035-11536801d1ff
	github.com/lcvvvv/appfinger v0.0.0-00010101000000-000000000000

	//gonmap
	github.com/lcvvvv/gonmap v1.3.4
	github.com/lcvvvv/pool v0.0.0-00010101000000-000000000000
	github.com/lcvvvv/simplehttp v0.0.0-00010101000000-000000000000
	github.com/lcvvvv/stdio v0.0.0-00010101000000-000000000000
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

replace github.com/lcvvvv/appfinger => ./lib/appfinger

replace github.com/lcvvvv/simplehttp => ./lib/simplehttp

replace github.com/lcvvvv/stdio => ./lib/stdio

replace github.com/lcvvvv/pool => ./lib/pool

//replace github.com/lcvvvv/gonmap => ../go-github/gonmap
