protoc -I=. --gogofaster_out=plugins=grpc:. --gogofaster_opt=paths=source_relative ./proto/user.proto

echo done gen protofile!