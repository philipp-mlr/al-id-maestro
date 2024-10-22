tailwind:
	@npx tailwindcss -i ./input.css -o ./public/css/style.css --watch

tailwind-minify:
	@npx tailwindcss -i ./input.css -o ./public/css/style.css --minify

templ:
	@templ generate -watch -proxy=http://localhost:1323
