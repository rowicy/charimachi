.PHONY: run

run_api_dev:
	cd api && go mod tidy && \
	( \
		trap "exit" INT; \
		go run . & \
		sleep 2; \
		(open http://localhost:8080/swagger/index.html || xdg-open http://localhost:8080/swagger/index.html); \
		wait \
	)


run_mobile_app:
	cd mobile-app && npm install && npm start