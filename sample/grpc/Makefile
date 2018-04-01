all:
	@rm route_guide.pb.go || true
	protoc -I . --go_out=plugins=grpc:. ./route_guide.proto
	@rm -fr ./mock || true
	@mkdir -p ./mock
	mockgen github.com/banaio/mock/sample/grpc RouteGuideClient,RouteGuide_ListFeaturesClient,RouteGuide_RecordRouteClient,RouteGuide_RouteChatClient > mock/route_guide.go
	@echo
	go test ./...
