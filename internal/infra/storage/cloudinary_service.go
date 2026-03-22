package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	infraconfig "github.com/huy/quizme-backend/internal/infra/config"
)

// FileType represents different types of files with their patterns and folders
type FileType string

const (
	FileTypeProfile         FileType = "PROFILE"
	FileTypeQuizThumbnail   FileType = "QUIZ_THUMBNAIL"
	FileTypeQuestionImage   FileType = "QUESTION_IMAGE"
	FileTypeQuestionAudio   FileType = "QUESTION_AUDIO"
	FileTypeCategoryIcon    FileType = "CATEGORY_ICON"
)

// CloudinaryService handles file uploads and management with Cloudinary
type CloudinaryService struct {
	cld    *cloudinary.Cloudinary
	config *infraconfig.Config
}

// NewCloudinaryService creates a new CloudinaryService
func NewCloudinaryService(cfg *infraconfig.Config) (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromParams(
		cfg.Cloudinary.CloudName,
		cfg.Cloudinary.APIKey,
		cfg.Cloudinary.APISecret,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary: %w", err)
	}

	return &CloudinaryService{
		cld:    cld,
		config: cfg,
	}, nil
}

// GetResourceURL builds the full Cloudinary URL for a file
func (s *CloudinaryService) GetResourceURL(folder, filename string) string {
	if filename == "" {
		return ""
	}
	return fmt.Sprintf("%s%s/image/upload/%s/%s",
		s.config.Cloudinary.BaseURL,
		s.config.Cloudinary.CloudName,
		folder,
		filename,
	)
}

// GetProfileImageURL returns the full URL for a profile avatar
func (s *CloudinaryService) GetProfileImageURL(filename string) string {
	if filename == "" {
		return ""
	}
	return s.GetResourceURL(s.config.Cloudinary.Folder["profile-avatar"], filename)
}

// GetQuizThumbnailURL returns the full URL for a quiz thumbnail
func (s *CloudinaryService) GetQuizThumbnailURL(filename string) string {
	if filename == "" {
		return ""
	}
	return s.GetResourceURL(s.config.Cloudinary.Folder["quiz-thumbnails"], filename)
}

// GetQuestionImageURL returns the full URL for a question image
func (s *CloudinaryService) GetQuestionImageURL(filename string) string {
	if filename == "" {
		return ""
	}
	return s.GetResourceURL(s.config.Cloudinary.Folder["question-images"], filename)
}

// GetQuestionAudioURL returns the full URL for a question audio
func (s *CloudinaryService) GetQuestionAudioURL(filename string) string {
	if filename == "" {
		return ""
	}
	return s.GetResourceURL(s.config.Cloudinary.Folder["question-audios"], filename)
}

// GetCategoryIconURL returns the full URL for a category icon
func (s *CloudinaryService) GetCategoryIconURL(filename string) string {
	if filename == "" {
		return ""
	}
	return s.GetResourceURL(s.config.Cloudinary.Folder["category-icons"], filename)
}

// GenerateProfileImageFilename generates a filename for profile avatar
func (s *CloudinaryService) GenerateProfileImageFilename(profileID uint, extension string) string {
	return fmt.Sprintf("profile_%d_%d%s", profileID, time.Now().UnixMilli(), extension)
}

// GenerateQuizThumbnailFilename generates a filename for quiz thumbnail
func (s *CloudinaryService) GenerateQuizThumbnailFilename(quizID uint, extension string) string {
	return fmt.Sprintf("quiz_thumbnail_%d_%d%s", quizID, time.Now().UnixMilli(), extension)
}

// GenerateQuestionImageFilename generates a filename for question image
func (s *CloudinaryService) GenerateQuestionImageFilename(quizID, questionID uint, extension string) string {
	return fmt.Sprintf("quiz_%d_question_%d_%d%s", quizID, questionID, time.Now().UnixMilli(), extension)
}

// GenerateQuestionAudioFilename generates a filename for question audio
func (s *CloudinaryService) GenerateQuestionAudioFilename(quizID, questionID uint, extension string) string {
	return fmt.Sprintf("quiz_%d_question_%d_%d%s", quizID, questionID, time.Now().UnixMilli(), extension)
}

// GenerateCategoryIconFilename generates a filename for category icon
func (s *CloudinaryService) GenerateCategoryIconFilename(categoryID uint, extension string) string {
	return fmt.Sprintf("category_%d_%d%s", categoryID, time.Now().UnixMilli(), extension)
}

// UploadProfileImage uploads a profile avatar to Cloudinary
func (s *CloudinaryService) UploadProfileImage(ctx context.Context, file multipart.File, originalFilename string, profileID uint) (string, error) {
	extension := getFileExtension(originalFilename)
	filename := s.GenerateProfileImageFilename(profileID, extension)
	folder := s.config.Cloudinary.Folder["profile-avatar"]
	return s.uploadFile(ctx, file, filename, folder)
}

// UploadQuizThumbnail uploads a quiz thumbnail to Cloudinary
func (s *CloudinaryService) UploadQuizThumbnail(ctx context.Context, file multipart.File, originalFilename string, quizID uint) (string, error) {
	extension := getFileExtension(originalFilename)
	filename := s.GenerateQuizThumbnailFilename(quizID, extension)
	folder := s.config.Cloudinary.Folder["quiz-thumbnails"]
	return s.uploadFile(ctx, file, filename, folder)
}

// UploadQuestionImage uploads a question image to Cloudinary
func (s *CloudinaryService) UploadQuestionImage(ctx context.Context, file multipart.File, originalFilename string, quizID, questionID uint) (string, error) {
	extension := getFileExtension(originalFilename)
	filename := s.GenerateQuestionImageFilename(quizID, questionID, extension)
	folder := s.config.Cloudinary.Folder["question-images"]
	return s.uploadFile(ctx, file, filename, folder)
}

// UploadQuestionAudio uploads a question audio to Cloudinary
func (s *CloudinaryService) UploadQuestionAudio(ctx context.Context, file multipart.File, originalFilename string, quizID, questionID uint) (string, error) {
	extension := getFileExtension(originalFilename)
	filename := s.GenerateQuestionAudioFilename(quizID, questionID, extension)
	folder := s.config.Cloudinary.Folder["question-audios"]
	return s.uploadFile(ctx, file, filename, folder)
}

// UploadCategoryIcon uploads a category icon to Cloudinary
func (s *CloudinaryService) UploadCategoryIcon(ctx context.Context, file multipart.File, originalFilename string, categoryID uint) (string, error) {
	extension := getFileExtension(originalFilename)
	filename := s.GenerateCategoryIconFilename(categoryID, extension)
	folder := s.config.Cloudinary.Folder["category-icons"]
	return s.uploadFile(ctx, file, filename, folder)
}

// uploadFile uploads a file to Cloudinary
func (s *CloudinaryService) uploadFile(ctx context.Context, file multipart.File, filename, folder string) (string, error) {
	publicID := strings.TrimSuffix(filename, filepath.Ext(filename))
	resourceType := determineResourceType(filepath.Ext(filename))

	uploadParams := uploader.UploadParams{
		PublicID:     publicID,
		Folder:       folder,
		ResourceType: resourceType,
	}

	_, err := s.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return filename, nil
}

// DeleteFile deletes a file from Cloudinary
func (s *CloudinaryService) DeleteFile(ctx context.Context, filename, folder string) error {
	if filename == "" {
		return nil
	}

	publicID := strings.TrimSuffix(filename, filepath.Ext(filename))
	fullPublicID := folder + "/" + publicID

	destroyParams := uploader.DestroyParams{
		PublicID:     fullPublicID,
		ResourceType: "image",
	}

	_, err := s.cld.Upload.Destroy(ctx, destroyParams)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// DeleteProfileImage deletes a profile avatar from Cloudinary
func (s *CloudinaryService) DeleteProfileImage(ctx context.Context, filename string) error {
	folder := s.config.Cloudinary.Folder["profile-avatar"]
	return s.DeleteFile(ctx, filename, folder)
}

// DeleteQuizThumbnail deletes a quiz thumbnail from Cloudinary
func (s *CloudinaryService) DeleteQuizThumbnail(ctx context.Context, filename string) error {
	folder := s.config.Cloudinary.Folder["quiz-thumbnails"]
	return s.DeleteFile(ctx, filename, folder)
}

// DeleteQuestionImage deletes a question image from Cloudinary
func (s *CloudinaryService) DeleteQuestionImage(ctx context.Context, filename string) error {
	folder := s.config.Cloudinary.Folder["question-images"]
	return s.DeleteFile(ctx, filename, folder)
}

// DeleteQuestionAudio deletes a question audio from Cloudinary
func (s *CloudinaryService) DeleteQuestionAudio(ctx context.Context, filename string) error {
	folder := s.config.Cloudinary.Folder["question-audios"]
	return s.DeleteFile(ctx, filename, folder)
}

// DeleteCategoryIcon deletes a category icon from Cloudinary
func (s *CloudinaryService) DeleteCategoryIcon(ctx context.Context, filename string) error {
	folder := s.config.Cloudinary.Folder["category-icons"]
	return s.DeleteFile(ctx, filename, folder)
}

// Helper functions

func getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return ".jpg"
	}
	return ext
}

func determineResourceType(extension string) string {
	extension = strings.ToLower(extension)

	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	videoExts := []string{".mp4", ".mov", ".avi", ".wmv"}
	audioExts := []string{".mp3", ".wav", ".ogg"}

	for _, ext := range imageExts {
		if extension == ext {
			return "image"
		}
	}

	for _, ext := range videoExts {
		if extension == ext {
			return "video"
		}
	}

	for _, ext := range audioExts {
		if extension == ext {
			return "auto"
		}
	}

	return "raw"
}
