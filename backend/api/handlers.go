package api

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"strings"

	"github.com/007Secret/007Secret/utils"
	"github.com/gin-gonic/gin"
)

const maxAttachmentSize = 10 * 1024 * 1024

type AttachmentMetadata struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Type string `json:"type"`
}

type SecretRequest struct {
	Content    string              `json:"content" binding:"required"`
	Attachment *AttachmentMetadata `json:"attachment,omitempty"`
}

type secretPayload struct {
	Content    string              `json:"content"`
	Attachment *AttachmentMetadata `json:"attachment,omitempty"`
}

type SecretResponse struct {
	Key      string `json:"key"`
	Password string `json:"password"`
}

type RetrievedSecretResponse struct {
	Content    string              `json:"content"`
	Attachment *AttachmentMetadata `json:"attachment,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:length]
}

func generateAlphanumericString(length int) (string, error) {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	result := make([]byte, length)
	alphabetLength := big.NewInt(int64(len(alphabet)))
	for i := range result {
		index, err := rand.Int(rand.Reader, alphabetLength)
		if err != nil {
			return "", err
		}
		result[i] = alphabet[index.Int64()]
	}

	return string(result), nil
}

// 生成加密密钥
func generateEncryptionKey(key, password string) []byte {
	// 使用SHA-256生成固定长度的密钥
	h := sha256.New()
	h.Write([]byte(key + password))
	return h.Sum(nil)
}

func compressData(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func decompressData(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var out bytes.Buffer
	if _, err := out.ReadFrom(r); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// decodeSecretPayload decodes stored bytes without reading or consuming Redis data.
func decodeSecretPayload(data, encryptionKey []byte) (secretPayload, error) {
	decompressed, err := decompressData(data)
	if err != nil {
		return secretPayload{}, errors.New("Failed to decompress data")
	}
	decrypted, err := utils.DecryptData(decompressed, encryptionKey)
	if err != nil {
		return secretPayload{}, errors.New("Failed to decrypt data")
	}

	var payload secretPayload
	if err := json.Unmarshal(decrypted, &payload); err != nil {
		// Secrets created before attachment support stored raw text.
		return secretPayload{Content: string(decrypted)}, nil
	}
	return payload, nil
}

func validateAttachmentMetadata(attachment *AttachmentMetadata) error {
	if attachment == nil {
		return nil
	}
	if strings.TrimSpace(attachment.Name) == "" || strings.ContainsAny(attachment.Name, "/\\") {
		return http.ErrNotSupported
	}
	if attachment.Size < 0 || attachment.Size > maxAttachmentSize {
		return http.ErrContentLength
	}
	return nil
}

func CreateSecretHandler(c *gin.Context) {
	var req SecretRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}
	if err := validateAttachmentMetadata(req.Attachment); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid attachment metadata"})
		return
	}

	// 生成随机密钥和密码
	key := generateRandomString(8)
	password, err := generateAlphanumericString(4)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to generate password",
		})
		return
	}

	// 生成加密密钥
	encryptionKey := generateEncryptionKey(key, password)

	payload, err := json.Marshal(secretPayload{Content: req.Content, Attachment: req.Attachment})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to encode secret"})
		return
	}

	// 先加密数据
	encrypted, err := utils.EncryptData(payload, encryptionKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to encrypt data",
		})
		return
	}

	// 再压缩数据
	compressed, err := compressData(encrypted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to compress data",
		})
		return
	}

	// 存储数据
	if err := utils.StoreSecret(key+":"+password, compressed); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to store secret",
		})
		return
	}

	c.JSON(http.StatusOK, SecretResponse{
		Key:      key,
		Password: password,
	})
}

func GetSecretHandler(c *gin.Context) {
	key := c.Query("key")
	password := c.Query("password")

	if key == "" || password == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Missing key or password",
		})
		return
	}

	// 获取并删除数据
	data, err := utils.GetSecret(key + ":" + password)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Secret not found or already accessed",
		})
		return
	}

	payload, err := decodeSecretPayload(data, generateEncryptionKey(key, password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, RetrievedSecretResponse{Content: payload.Content, Attachment: payload.Attachment})
}

func GetAttachmentMetadataHandler(c *gin.Context) {
	key, password := c.Query("key"), c.Query("password")
	if key == "" || password == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Missing key or password"})
		return
	}
	data, err := utils.PeekSecret(key + ":" + password)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Secret not found or already accessed"})
		return
	}
	payload, err := decodeSecretPayload(data, generateEncryptionKey(key, password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"attachment": payload.Attachment})
}
