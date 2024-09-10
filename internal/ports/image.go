package ports

type ImagePort interface {
	Upload(imageName, imageUrl string) error
	CheckImageAlreadyUploaded(imageName string) bool
}
