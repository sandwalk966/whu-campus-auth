package req

type UploadFileRequest struct {
	File interface{} `json:"file" binding:"required"`
}
