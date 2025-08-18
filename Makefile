.PHONY: run

run_api_dev:
	cd api && go mod tidy && go run . & \
	sleep 2 && \
	(open http://localhost:8080/swagger/index.html || xdg-open http://localhost:8080/swagger/index.html)

