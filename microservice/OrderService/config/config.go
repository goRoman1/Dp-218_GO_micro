package config

import (
	"os"
)

var PG_HOST = os.Getenv("PG_HOST")
var PG_PORT = os.Getenv("PG_PORT")
var POSTGRES_DB = os.Getenv("POSTGRES_DB")
var POSTGRES_USER = os.Getenv("POSTGRES_USER")
var POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
var TEMPLATES_PATH = os.Getenv("TEMPLATES_PATH")
var GRPC_PORT = os.Getenv("GRPC_PORT")
var ORDER_GRPC_PORT = os.Getenv("ORDER_GRPC_PORT")
var MONO_TEMPLATES_PATH = os.Getenv("MONO_TEMPLATES_PATH")
