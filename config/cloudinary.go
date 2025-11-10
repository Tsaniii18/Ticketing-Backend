package config

import (
	"context"	
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/joho/godotenv"
)

var Cld *cloudinary.Cloudinary

func InitCloudinary() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system environment")
	}

	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	log.Println("Cloudinary Config Check:")
	log.Println("Cloud Name:", cloudName)
	log.Println("API Key:", apiKey)
	log.Println("API Secret:", len(apiSecret), "chars")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		log.Fatal("Cloudinary credentials are missing")
	}

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Fatal("Failed to initialize Cloudinary:", err)
	}

	Cld = cld
	log.Println("Cloudinary initialized successfully")
}

func UploadImage(ctx context.Context, file interface{}, folder string) (string, error) {
	uploadResult, err := Cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return "", err
	}
	return uploadResult.SecureURL, nil
}

func UploadImageFromPath(ctx context.Context, filePath string, folder string) (string, error) {
	uploadResult, err := Cld.Upload.Upload(ctx, filePath, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return "", err
	}
	return uploadResult.SecureURL, nil
}

func DeleteImage(ctx context.Context, publicID string) error {
	_, err := Cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}