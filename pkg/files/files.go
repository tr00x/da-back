package files

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	// _ "github.com/chai2010/webp" // for decoding webp images
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
)

const maxFileSize = 10 * 1024 * 1024 // 10MB

var extensions map[string]bool = map[string]bool{
	"jpg":  true,
	"jpeg": true,
	"png":  true,
	// "webp": true,
	// "mp4":  true,
	// "heic": true,
}

// func WriteImage(c *fiber.Ctx, upload_path, folder, id string) string {
// 	base := upload_path + folder + id
// 	image, header, _ := c.Request.FormFile("image")

// 	if image == nil {
// 		return ""
// 	}
// 	splitedFileName := strings.Split(header.Filename, ".")
// 	extension := splitedFileName[len(splitedFileName)-1]
// 	filename := fmt.Sprintf("%v.", uuid.NewString()) + extension

// 	if extension == "webp" || extension == "svg" || extension == "jpeg" ||
// 		extension == "jpg" || extension == "png" {

// 		buf := bytes.NewBuffer(nil)
// 		io.Copy(buf, image)
// 		err := os.WriteFile(
// 			base+"/"+filename,
// 			buf.Bytes(), os.ModePerm,
// 		)

// 		if err != nil {
// 			return ""
// 		}
// 		return filename
// 	}
// 	return ""
// }

func SaveFiles(files []*multipart.FileHeader, base string, widths []uint) ([]string, int, error) {
	err := CreateFolderIfNotExists("." + base)

	if err != nil {
		return nil, 500, err
	}

	var filePaths []string
	var fileNames []string
	var video = 0
	var images = 0

	for index := range files {

		if files[index].Size > maxFileSize {
			return nil, 400, fmt.Errorf("file %s is too large", files[index].Filename)
		}

		splitedFileName := strings.Split(files[index].Filename, ".")
		extension := splitedFileName[len(splitedFileName)-1]
		extensionExists := extensions[extension]

		if !extensionExists {
			return nil, 400, fmt.Errorf("this file (extension) is forbidden: .%s", extension)
		}

		if extension == "mp4" {
			video += 1
		} else {
			images += 1
		}

		if video > 1 || images > 10 {
			return nil, 400, fmt.Errorf("trying to upload %v video and %v images", video, images)
		}
		fileNames = append(fileNames, uuid.NewString())
	}

	for index := range files {
		readerFile, _ := files[index].Open()
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, readerFile)

		err := os.WriteFile(
			"."+base+"/"+fileNames[index],
			buf.Bytes(),
			os.ModePerm,
		)

		if err != nil {
			return nil, 500, err
		}

		go func() {
			for _, width := range widths {

				err = ResizeImage("."+base+"/"+fileNames[index], width)
				if err != nil {
					return
				}

			}
			err = os.Remove("." + base + "/" + fileNames[index])

			if err != nil {
				return
			}
		}()
		filePaths = append(filePaths, base+"/"+fileNames[index])
	}
	return filePaths, 0, nil
}

// func ResizeImage(imagePath string, width uint) error {
// 	file, err := os.Open(imagePath)
// 	if err != nil {
// 		return fmt.Errorf("failed to open image: %w", err)
// 	}
// 	defer file.Close()

// 	img, _, err := image.Decode(file)
// 	if err != nil {
// 		return fmt.Errorf("failed to decode image: %w", err)
// 	}

// 	resizedImg := resize.Resize(width, 0, img, resize.Lanczos3)
// 	size := "l"

// 	if width == 320 {
// 		size = "m"
// 	}

// 	outputPath := strings.TrimSuffix(imagePath, filepath.Ext(imagePath)) + "_" + size + ".jpg"
// 	outFile, err := os.Create(outputPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to create output file: %w", err)
// 	}
// 	defer outFile.Close()

// 	options := jpeg.Options{Quality: 90}
// 	err = jpeg.Encode(outFile, resizedImg, &options)
// 	return err
// }

func ResizeImage(imagePath string, width uint) error {
	file, err := os.Open(imagePath)

	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)

	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// Reopen file to read EXIF
	file.Seek(0, 0)
	x, err := exif.Decode(file)

	if err == nil {
		orientTag, err := x.Get(exif.Orientation)
		if err == nil {
			orientation, _ := orientTag.Int(0)
			img = applyOrientation(img, orientation)
		}
	}

	// Resize the image
	resizedImg := resize.Resize(width, 0, img, resize.Lanczos3)
	size := "l"

	if width == 320 {
		size = "m"
	}

	outputPath := strings.TrimSuffix(imagePath, filepath.Ext(imagePath)) + "_" + size + ".jpg"
	outFile, err := os.Create(outputPath)

	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	options := jpeg.Options{Quality: 90}
	return jpeg.Encode(outFile, resizedImg, &options)
}

func applyOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 3:
		return imaging.Rotate180(img)
	case 6:
		return imaging.Rotate270(img)
	case 8:
		return imaging.Rotate90(img)
	default:
		return img
	}
}

func CreateFolderIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return nil
}

func RemoveFile(path string) error {
	dir := filepath.Dir("." + path)
	baseName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	entries, err := os.ReadDir(dir)

	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), baseName) {
			err := os.Remove(filepath.Join(dir, entry.Name()))

			if err != nil {
				return err
			}
		}
	}
	return nil
}

func RemoveFolder(path string) error {
	err := os.RemoveAll("." + path)

	return err
}

func SaveOriginal(file *multipart.FileHeader, folder string) (string, error) {
	err := CreateFolderIfNotExists("." + folder)

	if err != nil {
		return "", fmt.Errorf("failed to create output folder: %v", err)
	}

	// Open the uploaded file
	src, err := file.Open()

	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}

	defer src.Close()
	// Generate a UUID for the output file name
	outputID := uuid.New().String()
	// Save the uploaded file temporarily
	tempVideoPath := filepath.Join(folder, outputID+filepath.Ext(file.Filename))
	dst, err := os.Create("." + tempVideoPath)

	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}

	defer dst.Close()
	_, err = io.Copy(dst, src)

	if err != nil {
		return "", fmt.Errorf("failed to copy uploaded file: %v", err)
	}

	// Define output path for the HLS
	// outputPath := filepath.Join(folder, outputID+".m3u8")
	// go VideoToHLS(tempVideoPath, outputPath)
	return tempVideoPath, nil
}

func SaveVideos(file *multipart.FileHeader, folder string) (string, error) {
	err := CreateFolderIfNotExists("." + folder)

	if err != nil {
		return "", fmt.Errorf("failed to create output folder: %v", err)
	}

	// Open the uploaded file
	src, err := file.Open()

	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Generate a UUID for the output file name
	outputID := uuid.New().String()

	// Save the uploaded file temporarily
	tempVideoPath := filepath.Join("."+folder, outputID+"_upload"+filepath.Ext(file.Filename))
	dst, err := os.Create(tempVideoPath)

	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)

	if err != nil {
		return "", fmt.Errorf("failed to copy uploaded file: %v", err)
	}
	// Define output path for the HLS
	outputPath := filepath.Join(folder, outputID+".m3u8")
	go VideoToHLS(tempVideoPath, outputPath)
	return outputPath, nil // Return the UUID of the saved video
}

func VideoToHLS(tempVideoPath string, outputPath string) {
	// Run FFmpeg to convert video into HLS
	cmd := exec.Command(
		"ffmpeg",
		"-i", tempVideoPath,
		"-codec:v", "libx264",
		"-codec:a", "aac",
		"-strict", "-2",
		"-start_number", "0",
		"-hls_time", "5",
		"-f", "hls", "."+outputPath,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("failed to run ffmpeg: %v\nOutput: %s", err, string(output))
		os.Remove(tempVideoPath)
	}

	// Remove the temp uploaded video
	err = os.Remove(tempVideoPath)

	if err != nil {
		fmt.Printf("failed to remove temp file: %v", err)
	}
}
