PLAYER_IMAGES_BINARY=playerImagesApp

player_images:
	@echo "Building binary..."
	go build -o ${PLAYER_IMAGES_BINARY} ./
	@echo "Done!"

player_images_amd:
	@echo "Building binary..."
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${PLAYER_IMAGES_BINARY} ./
	@echo "Done!"
