module safer.place/realtime

go 1.20

require (
	api.safer.place v0.0.17
	connectrpc.com/connect v1.11.0
	github.com/bwmarrin/discordgo v0.27.1
	github.com/google/uuid v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/mattn/go-sqlite3 v1.14.17
	github.com/minio/minio-go/v7 v7.0.61
	go.uber.org/zap v1.24.0
	golang.org/x/exp v0.0.0-20230725093048-515e97ebf090
	golang.org/x/sync v0.3.0
	google.golang.org/protobuf v1.31.0
	safer.place/webserver v0.0.3
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/rs/cors v1.9.0 // indirect
	github.com/rs/xid v1.5.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

replace (
	api.safer.place v0.0.17 => github.com/saferplace/api v0.0.17
	safer.place/webserver v0.0.3 => github.com/saferplace/webserver-go v0.0.3
)
