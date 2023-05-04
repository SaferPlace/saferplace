module safer.place/realtime

go 1.20

require (
	api.safer.place v0.0.11
	github.com/kelseyhightower/envconfig v1.4.0
	go.uber.org/zap v1.24.0
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/rs/cors v1.8.2 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
)

require (
	github.com/bufbuild/connect-go v1.6.0
	github.com/bwmarrin/discordgo v0.27.1
	github.com/google/uuid v1.3.0
	github.com/mattn/go-sqlite3 v1.14.16
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/sync v0.1.0
	google.golang.org/protobuf v1.30.0
	safer.place/webserver v0.0.2
)

replace api.safer.place v0.0.11 => github.com/saferplace/api v0.0.11

// REMOVE ME
replace safer.place/webserver => ../webserver-go
