tailwind:
	@npx tailwindcss -i ./input.css -o ./public/css/style.css --watch

tailwind-minify:
	@npx tailwindcss -i ./input.css -o ./public/css/style.css --minify

templ:
	@templ generate -watch -proxy=http://localhost:1323

docker-build:
	docker build -t philippmue/philipp.software-website .

docker-run:
	docker run -d -p 1323:1323 philippmue/philipp.software-website

docker-push:
	docker push philippmue/philipp.software-website:latest