module github.com/pangum/migration

go 1.16

require (
	github.com/go-sql-driver/mysql v1.5.0
	github.com/pangum/logging v0.0.3
	github.com/pangum/pangu v0.0.1
	github.com/rubenv/sql-migrate v0.0.0-20210408115534-a32ed26c37ea
	github.com/storezhang/gox v1.7.9
	github.com/storezhang/simaqian v0.0.3
	github.com/ziutek/mymysql v1.5.4 // indirect
)

// replace github.com/pangum/pangu => ../../storezhang/pangu
// replace github.com/storezhang/gox => ../../storezhang/gox
