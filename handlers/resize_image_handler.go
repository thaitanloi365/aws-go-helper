package handlers

import (
	"aws-go-helper/config"
	"aws-go-helper/models"
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/disintegration/imaging"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/echo"
)

// ResizeImageHandler resize image handler
func ResizeImageHandler(c echo.Context) error {
	var optional = c.Param("optional")
	var img = &models.Image{
		Optional: optional,
	}
	var cfg = config.Instance

	if !img.IsMatchFormat() {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	log.Printf("Img: %v | %v | %v | %v | %v | %v\n", img.FileName, img.Crop, img.Extension, img.Dimension, img.Height, img.Width)

	// Init session
	sess := session.Must(session.NewSession())
	sess.Config.Region = aws.String(cfg.AWSS3Region)
	srv := s3.New(sess)
	ctx := c.Request().Context()

	// Check if resized already exits, just return
	targetKey := img.GetS3Key(cfg.AWSS3ResizedFolder, img.GetOutputFileName())
	_, err := srv.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(cfg.AWSS3Bucket),
		Key:    aws.String(targetKey),
	})
	if err == nil {
		log.Printf("%s already exits", targetKey)
		var resizedImageLocation = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", cfg.AWSS3Bucket, targetKey)
		return c.Redirect(http.StatusTemporaryRedirect, resizedImageLocation)
	}

	originalKey := img.GetS3Key(cfg.AWSS3OriginFolder, img.FileName)

	// Download image
	_, err = srv.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(cfg.AWSS3Bucket),
		Key:    aws.String(originalKey),
	})

	if err != nil {
		log.Printf("%s not found in %s bucket: %v\n", originalKey, cfg.AWSS3Bucket, err)
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s not found in %s bucket", originalKey, cfg.AWSS3Bucket))

	}

	buffer := &aws.WriteAtBuffer{}
	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.DownloadWithContext(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(cfg.AWSS3Bucket),
		Key:    aws.String(originalKey),
	})
	if err != nil {
		log.Printf("could not download image from S3: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not download image from S3")
	}

	//Decode from downloaded data
	originalImage, err := imaging.Decode(bytes.NewReader(buffer.Bytes()))
	if err != nil {
		log.Printf("Decode image error: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not decode downloaded image from S3")
	}

	//Resize image
	resizedImage := img.ResizeOrCrop(originalImage)
	if resizedImage == nil {
		log.Printf("Resized image error: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not resize image")
	}

	//Encode resize image and upload to s3
	var bufferEncode = new(bytes.Buffer)
	errEncode := imaging.Encode(bufferEncode, resizedImage, models.ParseExtension(models.ParseContentType(img.Extension)))
	if errEncode != nil {
		log.Printf("Encode image error: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not encode image")
	}

	//Upload to S3
	uploader := s3manager.NewUploader(sess)
	output, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(cfg.AWSS3Bucket),
		Key:         aws.String(targetKey),
		ContentType: aws.String("image/jpeg"),
		Body:        bytes.NewReader(bufferEncode.Bytes()),
		ACL:         aws.String(cfg.AWSS3ACL),
	})

	if err != nil {
		log.Printf("Upload image error: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not upload resized image")
	}
	//Serve file just uploaded
	log.Printf("New image has uploaded at: %v\n", output.Location)

	return c.Redirect(http.StatusTemporaryRedirect, output.Location)
}
