module github.com/jamespfennell/transiter

go 1.23.0

toolchain go1.24.2

require (
	connectrpc.com/connect v1.18.1
	github.com/benbjohnson/clock v1.3.5
	github.com/google/go-cmp v0.7.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3
	github.com/jackc/pgx/v5 v5.7.5
	github.com/jackc/tern/v2 v2.3.3
	github.com/jamespfennell/gtfs v0.1.24
	github.com/prometheus/client_golang v1.22.0
	github.com/pseudomuto/protoc-gen-doc v1.5.1
	github.com/urfave/cli/v2 v2.27.6
	golang.org/x/exp v0.0.0-20250531010427-b6e5de432a8b
	google.golang.org/genproto/googleapis/api v0.0.0-20250528174236-200df99c418a
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.6
)

require (
	dario.cat/mergo v1.0.2 // indirect
	github.com/Masterminds/semver/v3 v3.3.1 // indirect
	github.com/Masterminds/sprig/v3 v3.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/spf13/cast v1.9.2 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	golang.org/x/sync v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
)

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/pseudomuto/protokit v0.2.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/net v0.41.0
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/jamespfennell/gtfs => github.com/milas/gtfs-go v0.0.0-20250531194311-3306bf8cb3c3

//replace github.com/jamespfennell/gtfs => ../gtfs-go
