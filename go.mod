module safer.place/realtime

go 1.20

require (
	api.safer.place v0.0.14
	connectrpc.com/connect v1.11.0
	github.com/bwmarrin/discordgo v0.27.1
	github.com/google/uuid v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/mattn/go-sqlite3 v1.14.17
	go.uber.org/zap v1.24.0
	golang.org/x/exp v0.0.0-20230725093048-515e97ebf090
	golang.org/x/sync v0.3.0
	google.golang.org/protobuf v1.31.0
	safer.place/webserver v0.0.3
)

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/rs/cors v1.9.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
)

replace (
	api.safer.place v0.0.14 => github.com/saferplace/api v0.0.14
	safer.place/webserver v0.0.3 => github.com/saferplace/webserver-go v0.0.3
)
