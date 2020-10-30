module github.com/MaibornWolff/maDocK8s/core/utils

go 1.15

replace github.com/MaibornWolff/maDocK8s/core/types => ../types

require (
	github.com/MaibornWolff/maDocK8s/core/types v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.33.1
)
