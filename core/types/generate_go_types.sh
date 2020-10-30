cd $GOPATH/src/github.com/MaibornWolff/maDocK8s/core/types
rm  `find . -name "*.pb.go"`
protowrap -I. --go_out=plugins=grpc:`go env GOPATH`/src ./**/*.proto
