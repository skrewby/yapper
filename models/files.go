package models

import (
	"image"
	"image/png"
	"log/slog"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Files struct {
	db *pgxpool.Pool
}

func NewFilesModel(db *pgxpool.Pool) *Files {
	return &Files{
		db,
	}
}

func (m Files) UploadImage(file multipart.File) (string, error) {
	img, header, err := image.Decode(file)
	if err != nil {
		slog.Error("Could not decode image binary:" + err.Error())
		return "", err
	}

	fileName := "/files/" + uuid.NewString() + "." + header
	path := ".dev" + fileName
	f, err := os.Create(path)
	if err != nil {
		slog.Error("Could not create file" + path + ": " + err.Error())
		return "", err
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		slog.Error("Could not encode image as png: " + err.Error())
		return "", err
	}

	if err := f.Close(); err != nil {
		slog.Error("Could not close file: " + err.Error())
		return "", err
	}

	return fileName, nil
}
