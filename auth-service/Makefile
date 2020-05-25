build: 
	nmake login
	nmake register
	nmake token
	go build main.go

run-user:
	cd ../user-service/
	go run main.go

run-auth:
	go run main.go

login:
	protoc -I. --go_out=plugins=grpc:./services/login --grpc-gateway_out=logtostderr=true:./services/login --swagger_out=logtostderr=true:./services/login --proto_path proto login.proto 
	mv ./services/login/egnite.app/microservices/authentication/proto/login/* ./services/login/
	rm -r ./services/login/egnite.app
	cp ./services/login/*.json ./www/
register:
	protoc -I. --go_out=plugins=grpc:./services/register --grpc-gateway_out=logtostderr=true:./services/register --swagger_out=logtostderr=true:./services/register --proto_path proto register.proto
	mv ./services/register/egnite.app/microservices/authentication/proto/register/* ./services/register/
	rm -r ./services/register/egnite.app
	cp ./services/register/*.json ./www/

token:
	protoc -I. --go_out=plugins=grpc:./services/token --grpc-gateway_out=logtostderr=true:./services/token --swagger_out=logtostderr=true:./services/token --proto_path proto token.proto
	mv ./services/token/egnite.app/microservices/authentication/proto/token/* ./services/token/
	rm -r ./services/token/egnite.app
	cp ./services/token/*.json ./www/
