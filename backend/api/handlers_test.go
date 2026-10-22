package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/007Secret/007Secret/utils"
	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func TestGenerateAlphanumericString(t *testing.T) {
	validCode := regexp.MustCompile(`^[A-Za-z0-9]{4}$`)

	for range 1_000 {
		code, err := generateAlphanumericString(4)
		if err != nil {
			t.Fatalf("generateAlphanumericString returned an error: %v", err)
		}
		if !validCode.MatchString(code) {
			t.Fatalf("generated extraction code contains invalid characters: %q", code)
		}
	}
}

func TestSecretRequestDeclaresOptionalAttachmentMetadata(t *testing.T) {
	requestType := reflect.TypeOf(SecretRequest{})
	attachmentField, found := requestType.FieldByName("Attachment")
	if !found {
		t.Fatal("SecretRequest must declare optional attachment metadata")
	}

	if attachmentField.Tag.Get("json") != "attachment,omitempty" {
		t.Fatalf("attachment JSON tag = %q, want attachment,omitempty", attachmentField.Tag.Get("json"))
	}
}

func TestCreateSecretRejectsAttachmentLargerThanTenMegabytes(t *testing.T) {
	previousClient := utils.RedisClient
	utils.RedisClient = redis.NewClient(&redis.Options{
		Addr:        "127.0.0.1:0",
		DialTimeout: 10 * time.Millisecond,
	})
	defer func() {
		_ = utils.RedisClient.Close()
		utils.RedisClient = previousClient
	}()

	request := httptest.NewRequest(http.MethodPost, "/api/secret", strings.NewReader(`{
		"content":"report.pdf",
		"attachment":{"name":"report.pdf","size":10485761,"type":"application/pdf"}
	}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Request = request

	CreateSecretHandler(context)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestSecretHandlersRoundTripAttachmentMetadata(t *testing.T) {
	server := miniredis.RunT(t)
	previousClient := utils.RedisClient
	utils.RedisClient = redis.NewClient(&redis.Options{Addr: server.Addr()})
	defer func() {
		_ = utils.RedisClient.Close()
		utils.RedisClient = previousClient
	}()

	createRequest := httptest.NewRequest(http.MethodPost, "/api/secret", strings.NewReader(`{
		"content":"report.pdf",
		"attachment":{"name":"report.pdf","size":10485760,"type":"application/pdf"}
	}`))
	createRequest.Header.Set("Content-Type", "application/json")
	createResponse := httptest.NewRecorder()
	createContext, _ := gin.CreateTestContext(createResponse)
	createContext.Request = createRequest
	CreateSecretHandler(createContext)

	if createResponse.Code != http.StatusOK {
		t.Fatalf("create status = %d, want %d", createResponse.Code, http.StatusOK)
	}

	var created SecretResponse
	if err := json.Unmarshal(createResponse.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/api/secret/get?key="+created.Key+"&password="+created.Password, nil)
	getResponse := httptest.NewRecorder()
	getContext, _ := gin.CreateTestContext(getResponse)
	getContext.Request = getRequest
	GetSecretHandler(getContext)

	if getResponse.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d", getResponse.Code, http.StatusOK)
	}

	var retrieved RetrievedSecretResponse
	if err := json.Unmarshal(getResponse.Body.Bytes(), &retrieved); err != nil {
		t.Fatalf("decode get response: %v", err)
	}
	if retrieved.Content != "report.pdf" || retrieved.Attachment == nil {
		t.Fatalf("retrieved secret = %#v, want content and attachment", retrieved)
	}
	if got := *retrieved.Attachment; got != (AttachmentMetadata{Name: "report.pdf", Size: 10485760, Type: "application/pdf"}) {
		t.Fatalf("attachment = %#v, want report metadata", got)
	}
}

func TestAttachmentMetadataCheckDoesNotConsumeSecret(t *testing.T) {
	server := miniredis.RunT(t)
	previousClient := utils.RedisClient
	utils.RedisClient = redis.NewClient(&redis.Options{Addr: server.Addr()})
	defer func() { _ = utils.RedisClient.Close(); utils.RedisClient = previousClient }()

	createRequest := httptest.NewRequest(http.MethodPost, "/api/secret", strings.NewReader(`{"content":"report.pdf","attachment":{"name":"report.pdf","size":3,"type":"application/pdf"}}`))
	createRequest.Header.Set("Content-Type", "application/json")
	createResponse := httptest.NewRecorder()
	createContext, _ := gin.CreateTestContext(createResponse)
	createContext.Request = createRequest
	CreateSecretHandler(createContext)

	var created SecretResponse
	if err := json.Unmarshal(createResponse.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}

	metadataRequest := httptest.NewRequest(http.MethodGet, "/api/secret/attachment?key="+created.Key+"&password="+created.Password, nil)
	metadataResponse := httptest.NewRecorder()
	metadataContext, _ := gin.CreateTestContext(metadataResponse)
	metadataContext.Request = metadataRequest
	GetAttachmentMetadataHandler(metadataContext)
	if metadataResponse.Code != http.StatusOK {
		t.Fatalf("metadata status = %d, want %d", metadataResponse.Code, http.StatusOK)
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/api/secret/get?key="+created.Key+"&password="+created.Password, nil)
	getResponse := httptest.NewRecorder()
	getContext, _ := gin.CreateTestContext(getResponse)
	getContext.Request = getRequest
	GetSecretHandler(getContext)
	if getResponse.Code != http.StatusOK {
		t.Fatalf("get after metadata status = %d, want %d", getResponse.Code, http.StatusOK)
	}
}
