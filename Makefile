start:
	go build -a -installsuffix cgo -o ./example-service 
	./example-service  start
createDb:
	go build -o ./example-service  
	./example-service  createDb
scan:
	go build -o ./example-service 
	./example-service  scan