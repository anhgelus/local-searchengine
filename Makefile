.PHONY: godev
build: 
	pnpm run build
	go build -o local-searchengine .

.PHONY: godev
run:
	./local-searchengine

.PHONY: godev
dev:
	make build
	make run

.PHONY: frontdev
frontdev:
	pnpm run dev

.PHONY: install
install: local-searchengine
	./local-searchengine install
