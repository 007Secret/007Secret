package api

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/007Secret/007Secret/utils"
	"github.com/gin-gonic/gin"
)

type SecretRequest struct {
	Content string `json:"content" binding:"required"`
}

type SecretResponse struct {
	Key      string `json:"key"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:length]
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

func CreateSecretHandler(c *gin.Context) {
	var req SecretRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// 生成随机密钥和密码
	key := generateRandomString(8)
	password := generateRandomString(4)

	// 压缩数据
	compressed, err := compressData([]byte(req.Content))
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

	// 解压数据
	decompressed, err := decompressData(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to decompress data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"content": string(decompressed),
	})
}
