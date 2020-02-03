package handlers

import (
	"aws-go-helper/config"
	"aws-go-helper/models"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/echo"
)

// ResizeImageHandler resize image handler
func ResizeImageHandler(c echo.Context) error {
	var name = c.Param("name")
	var crop = c.QueryParam("crop")
	var size = c.QueryParam("size")

	var img = &models.Image{
		FileName:  name,
		Width:     100,
		Height:    100,
		Dimension: "100x100",
	}

	if crop != "" {
		img.Crop = models.CropOption(crop)
	}

	var sizeParams = strings.Split(size, "x")
	if len(sizeParams) > 2 {
		if w, err := strconv.ParseInt(sizeParams[0], 10, 64); err == nil {
			img.Width = int(w)
		}

		if h, err := strconv.ParseInt(sizeParams[1], 10, 64); err == nil {
			img.Height = int(h)
		}

	}

	var cfg = config.Instance

	log.Printf("Img: %v | %v | %v | %v | %v | %v\n", img.FileName, img.Crop, img.Extension, img.Dimension, img.Height, img.Width)

	// Init session
	sess := session.Must(session.NewSession())
	sess.Config.Region = aws.String(cfg.AWSS3Region)
	srv := s3.New(sess)
	ctx := c.Request().Context()

	sessOfResize := session.Must(session.NewSession())
	sessOfResize.Config.Region = aws.String(cfg.AWSS3ResizedBucket)
	srvOfResize := s3.New(sessOfResize)

	// Check if resized already exits, just return
	targetKey := img.GetS3Key("", img.GetOutputFileName())
	_, err := srvOfResize.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(cfg.AWSS3ResizedBucket),
		Key:    aws.String(targetKey),
	})
	if err == nil {
		log.Printf("%s already exits", targetKey)
		var resizedImageLocation = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", cfg.AWSS3ResizedBucket, targetKey)
		return c.Redirect(http.StatusTemporaryRedirect, resizedImageLocation)
	}

	originalKey := img.GetS3Key("", img.FileName)

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
		Bucket:      aws.String(cfg.AWSS3ResizedBucket),
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
