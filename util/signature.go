package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SignatureValidator struct {
	appID     string
	secretKey string
}

func NewSignatureValidator(appID, secretKey string) *SignatureValidator {
	return &SignatureValidator{
		appID:     appID,
		secretKey: secretKey,
	}
}

type SignatureRequest struct {
	Timestamp int64  `json:"timestamp"`
	Sign      string `json:"sign"`
	Nonce     string `json:"nonce"`
}

func (sv *SignatureValidator) GenerateSignature(params map[string]interface{}) string {
	timestamp := time.Now().UnixMilli()
	nonce := generateNonce(16)

	params["timestamp"] = timestamp
	params["nonce"] = nonce
	params["app_id"] = sv.appID

	sortedKeys := make([]string, 0, len(params))
	for k := range params {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	signStr := ""
	for _, k := range sortedKeys {
		signStr += fmt.Sprintf("%s=%v&", k, params[k])
	}
	signStr += "key=" + sv.secretKey

	hash := sha256.Sum256([]byte(signStr))
	signature := hex.EncodeToString(hash[:])

	return signature
}

func (sv *SignatureValidator) ValidateSignature(params map[string]interface{}) bool {
	timestamp, ok := params["timestamp"].(int64)
	if !ok {
		return false
	}

	if math.Abs(float64(time.Now().UnixMilli()-timestamp)) > 5*60*1000 {
		return false
	}

	sign, ok := params["sign"].(string)
	if !ok {
		return false
	}

	paramsCopy := make(map[string]interface{})
	for k, v := range params {
		if k != "sign" {
			paramsCopy[k] = v
		}
	}

	sortedKeys := make([]string, 0, len(paramsCopy))
	for k := range paramsCopy {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	signStr := ""
	for _, k := range sortedKeys {
		signStr += fmt.Sprintf("%s=%v&", k, paramsCopy[k])
	}
	signStr += "key=" + sv.secretKey

	hash := sha256.Sum256([]byte(signStr))
	calculatedSign := hex.EncodeToString(hash[:])

	return hmac.Equal([]byte(sign), []byte(calculatedSign))
}

func generateNonce(length int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	timestamp := time.Now().UnixNano()
	for i := range result {
		result[i] = chars[(int(timestamp)+i*17)%len(chars)]
	}
	return string(result)
}

type SignatureContext struct{}

func ValidateRequestSignature(c *gin.Context) bool {
	timestamp := c.GetHeader("X-Timestamp")
	sign := c.GetHeader("X-Sign")
	nonce := c.GetHeader("X-Nonce")

	if timestamp == "" || sign == "" || nonce == "" {
		return false
	}

	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}

	if math.Abs(float64(time.Now().UnixMilli()-timestampInt)) > 5*60*1000 {
		return false
	}

	params := map[string]interface{}{
		"timestamp": timestampInt,
		"nonce":     nonce,
	}

	c.Request.ParseForm()
	for key, values := range c.Request.Form {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	sortedKeys := make([]string, 0, len(params))
	for k := range params {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	signStr := ""
	for _, k := range sortedKeys {
		signStr += fmt.Sprintf("%s=%v&", k, params[k])
	}

	secretKey := "your-secret-key"
	signStr += "key=" + secretKey

	hash := sha256.Sum256([]byte(signStr))
	calculatedSign := hex.EncodeToString(hash[:])

	return hmac.Equal([]byte(sign), []byte(calculatedSign))
}
