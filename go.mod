module Dp218GO

go 1.17

replace problem.micro => ./microservice/ProblemMicro

replace supplier.micro => ./microservice/SupplierMicro

replace scooter_client.micro => ./microservice/ScooterClient

replace order.micro => ./microservice/OrderService

replace scooter_server.micro => ./microservice/ScooterServer

require (
	github.com/gorilla/sessions v1.2.1
	golang.org/x/crypto v0.0.0-20211209193657-4570a0811e8b
)

require (
	github.com/Shopify/sarama v1.30.1
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/golang-migrate/migrate/v4 v4.15.1
	github.com/golang/mock v1.6.0
	github.com/gorilla/mux v1.8.0
	github.com/jackc/pgconn v1.10.1
	github.com/jackc/pgx/v4 v4.14.0
	github.com/jszwec/csvutil v1.6.0
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
	problem.micro v0.0.0-00010101000000-000000000000
	supplier.micro v0.0.0-00010101000000-000000000000
)

require (
	github.com/asaskevich/govalidator v0.0.0-20200108200545-475eaeb16496 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.2.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.9.0 // indirect
	github.com/jackc/puddle v1.2.0 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.2 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/lib/pq v1.10.4 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f // indirect
	golang.org/x/sys v0.0.0-20211210111614-af8b64212486 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211013025323-ce878158c4d4 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
