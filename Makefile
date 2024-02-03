watch-css:
	@echo "Watching and compiling Tailwind CSS..."
	npx tailwindcss -i ./frontend/public/styles/global.css -o ./frontend/public/styles/tailwind.css --watch

dev:
	./bin/air