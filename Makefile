IMAGE_TAG=golang:1.15
NAME=mdr

.PHONY: build
build: .env 
	@docker run --rm -t --env-file .env -v "${PWD}":/usr/src/app -w /usr/src/app ${IMAGE_TAG} \
		go build -v -o ./bin/${NAME} ./cmd/${NAME}

.PHONY: watch
watch:
	@inotifywait -e close_write,moved_to,create -rmq pkg cmd @bin | \
	while read -r directory events filename; do \
		echo $$'\n'; echo "> $$directory$$filename changed. Recompiling...";\
		make -s build; \
	done
	
.PHONY: run
run: .env build
	./bin/${NAME}

.env: .env.dist
	cp .env.dist .env
