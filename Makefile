tailwind:
	@npx tailwindcss -i ./input.css -o ./website/public/css/style.css --watch

tailwind-minify:
	@npx tailwindcss -i ./input.css -o ./website/public/css/style.css --minify

templ:
	@templ generate -watch -proxy=http://localhost:5000
