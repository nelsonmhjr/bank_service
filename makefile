default:
	@echo "=============building Local API============="
	docker build -f Dockerfile -t api .

up: default
	@echo "=============starting api locally============="
	docker-compose up -d

logs:
	docker-compose logs -f --tail=100 -t api

down:
	docker-compose down

test:
	docker-compose run api go test -v -cover -coverprofile coverage.out ./...

test_doc: 
	docker-compose run api ginkgo -r --v --reportPassed -cover -coverprofile coverage.out

test_local:
	godotenv go test -v -cover -coverprofile coverage.out ./...

test_local_debug:
	DEBUG_TEST=true godotenv ginkgo -r --v --reportPassed --trace -cover -coverprofile coverage.out

cover_report:
	go tool cover -html=coverage.out

clean: down
	@echo "=============cleaning up============="
	rm -f api
	docker system prune -f
	docker volume prune -f