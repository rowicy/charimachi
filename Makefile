.PHONY: run

run_api_dev:
	cd api && go mod tidy && \
	( \
		kill -9 $$(lsof -ti :8080 2>/dev/null) 2>/dev/null; \
		trap 'echo "SIGINT received, killing server..."; kill -9 $$(lsof -ti tcp:8080 2>/dev/null) 2>/dev/null; exit' INT; \
		go run . & \
		sleep 2; \
		(open http://localhost:8080/swagger/index.html || xdg-open http://localhost:8080/swagger/index.html); \
		wait \
	)


run_mobile_app:
	cd mobile-app && npm install && npm start