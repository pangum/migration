module github.com/storezhang/pangu-migration

go 1.16

require (
	github.com/go-sql-driver/mysql v1.5.0
	github.com/rubenv/sql-migrate v0.0.0-20210408115534-a32ed26c37ea
	github.com/storezhang/glog v1.0.8
	github.com/storezhang/gox v1.5.2
	github.com/storezhang/pangu v1.2.4
	github.com/storezhang/pangu-logging v1.0.0
)

replace github.com/storezhang/pangu => ../../storezhang/pangu
replace github.com/storezhang/gox => ../../storezhang/gox
