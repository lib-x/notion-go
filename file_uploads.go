package notion

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
)

// FileUploadsService implements Notion file upload endpoints.
type FileUploadsService struct {
	client *Client
}

// CreateFileUploadRequest creates a Notion file upload object.
type CreateFileUploadRequest struct {
	Mode          string `json:"mode,omitempty"`
	Filename      string `json:"filename,omitempty"`
	ContentType   string `json:"content_type,omitempty"`
	NumberOfParts int    `json:"number_of_parts,omitempty"`
	ExternalURL   string `json:"external_url,omitempty"`
}

// ListFileUploadsParams configures file upload listing.
type ListFileUploadsParams struct {
	Status string
	PaginationParams
}

// UploadFileRequest uploads one file or one multipart part.
type UploadFileRequest struct {
	Filename   string
	Reader     io.Reader
	PartNumber string
}

// Create creates a file upload object.
func (s *FileUploadsService) Create(ctx context.Context, body CreateFileUploadRequest) (FileUpload, error) {
	var out Object
	err := s.client.post(ctx, apiPath("v1", "file_uploads"), body, &out)
	return FileUpload(out), err
}

// List lists file uploads.
func (s *FileUploadsService) List(ctx context.Context, params *ListFileUploadsParams) (*ListResponse, error) {
	q := make(url.Values)
	if params != nil {
		addString(q, "status", params.Status)
		for key, values := range paginationValues(&params.PaginationParams) {
			q[key] = values
		}
	}
	var out ListResponse
	err := s.client.get(ctx, apiPath("v1", "file_uploads"), q, &out)
	return &out, err
}

// Send uploads file bytes for a file upload object.
func (s *FileUploadsService) Send(ctx context.Context, fileUploadID string, req UploadFileRequest) (FileUpload, error) {
	if req.Reader == nil {
		return nil, fmt.Errorf("notion: upload file reader is nil")
	}
	if req.Filename == "" {
		req.Filename = "file"
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writeMultipartFile(writer, "file", req.Filename, req.Reader); err != nil {
		_ = writer.Close()
		return nil, fmt.Errorf("notion: write multipart file: %w", err)
	}
	if req.PartNumber != "" {
		if err := writer.WriteField("part_number", req.PartNumber); err != nil {
			_ = writer.Close()
			return nil, fmt.Errorf("notion: write multipart part_number: %w", err)
		}
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("notion: close multipart writer: %w", err)
	}

	var out Object
	err := s.client.doMultipart(ctx, apiPath("v1", "file_uploads", fileUploadID, "send"), &body, writer.FormDataContentType(), &out)
	return FileUpload(out), err
}

// Complete completes a multi-part file upload.
func (s *FileUploadsService) Complete(ctx context.Context, fileUploadID string) (FileUpload, error) {
	var out Object
	err := s.client.post(ctx, apiPath("v1", "file_uploads", fileUploadID, "complete"), nil, &out)
	return FileUpload(out), err
}

// Get retrieves a file upload object.
func (s *FileUploadsService) Get(ctx context.Context, fileUploadID string) (FileUpload, error) {
	var out Object
	err := s.client.get(ctx, apiPath("v1", "file_uploads", fileUploadID), nil, &out)
	return FileUpload(out), err
}
