package videoCompress

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
)

var videoResolutions = []string{"1280x720", "720x480", "480x240", "256x144"}
type ResolutionDoesNotExist struct {
	Message string
}
func (r ResolutionDoesNotExist) Error() string {
	return r.Message
}

func CompressFileToSize(fileHeader *multipart.FileHeader,filename string, requiredSize int64) ( *CompressedFileInfo, error) {
	requiredSizeMB := requiredSize / (1024 * 1024)
	var compressedFile *CompressedFileInfo
	compressedFile ,err := CompressFile(fileHeader, filename, videoResolutions[0])
	if err != nil {
		return nil,err
	}
	compressedFileSizeMB := compressedFile.CompressedFileSize / (1024 * 1024)
	for i:=1;compressedFileSizeMB > requiredSizeMB;i++ {
		err := CleanupCompressedFile(compressedFile.CompressedFilename)
		if err != nil {

			return nil, err
		}
		if i >= len(videoResolutions) {
			return nil, ResolutionDoesNotExist{Message: "a"}
		}
		compressedFile ,err = CompressFile(fileHeader, filename, videoResolutions[i])
		if err != nil {
			return nil, err
		}
		compressedFileSizeMB = compressedFile.CompressedFileSize / (1024 * 1024)
	}
	return compressedFile, err
}

type CompressedFileInfo struct {
	CompressedFilename string
	CompressedFileSize int64
}



func CompressFile(fileHeader *multipart.FileHeader,filename string, resolution string) (*CompressedFileInfo, error) {
	var compressedFileInfo CompressedFileInfo
	file, err := fileHeader.Open()
	if err != nil {
		return &compressedFileInfo,err
	}
	defer file.Close()
	tempFile, err := os.CreateTemp("temp", fmt.Sprintf("*.%s.mp4", filename))
	if err != nil {
		return &compressedFileInfo,err
	}

	_, err = io.Copy(tempFile, file)
	if err != nil {
		return &compressedFileInfo,err
	}
	outputFileName := fmt.Sprintf("%s_compressed.mp4", filename)
	outputFilePath := filepath.Join("/", outputFileName)

	cmd := GetCompressingCommand(resolution, tempFile, outputFilePath)
	err = cmd.Run()
	if err!=nil {
		return &compressedFileInfo,err
	}
	compressedVideo, err := os.Open(outputFilePath)
	defer compressedVideo.Close()
	if err != nil {
		return &compressedFileInfo, err
	}
	tempFile.Close()
	err = os.Remove(filepath.Join("", tempFile.Name()))
	if err != nil {
		return &compressedFileInfo, err
	}
	fileInfo, _ := compressedVideo.Stat()
	//compressedFileInfo.Reader = compressedVideo
	compressedFileInfo.CompressedFilename = outputFileName
	compressedFileInfo.CompressedFileSize = fileInfo.Size()
	fmt.Println("crying here 2")
	return &compressedFileInfo,nil
}

func GetCompressingCommand(res string, inputFile *os.File, outputFilePath string) *exec.Cmd {
	cmd := exec.Command("ffmpeg",
		"-i", inputFile.Name(),
		"-s", res,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-crf", "23",
		outputFilePath,
		)
	return cmd
}

//func GetResolutionFromFileSize(fileSize int64,) string {
//	fileSizeMB := fileSize / (1024 * 1024)
//	if fileSizeMB > 1000 {
//	}
//}

func CleanupCompressedFile(filepath string) error {
	err := os.Remove(filepath)
	if err != nil {
		return err
	}
	return nil
}