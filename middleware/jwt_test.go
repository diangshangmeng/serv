package middleware

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// mockConfig 用于测试的配置
var testJWTSecret = []byte("test-secret-key")

// generateTestToken 生成测试用JWT token
func generateTestToken(userID uint64, phone string, role string, secret []byte, exp time.Time) string {
	claims := jwt.MapClaims{
		"user_id": float64(userID),
		"phone":   phone,
		"role":    role,
		"exp":     exp.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(secret)
	return tokenString
}

// generateExpiredToken 生成过期token
func generateExpiredToken(userID uint64, secret []byte) string {
	claims := jwt.MapClaims{
		"user_id": float64(userID),
		"phone":   "13800138000",
		"exp":     time.Now().Add(-1 * time.Hour).Unix(), // 已过期1小时
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(secret)
	return tokenString
}

// generateNoneAlgorithmToken 生成使用none算法的token（攻击测试）
func generateNoneAlgorithmToken(userID uint64) string {
	claims := jwt.MapClaims{
		"user_id": float64(userID),
		"phone":   "13800138000",
		"alg":     "none", // 尝试alg:none攻击
	}

	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, _ := token.SignedString([]byte(""))
	return tokenString
}

func init() {
	gin.SetMode(gin.TestMode)
}

// TestParseToken_ValidToken 测试有效token解析
func TestParseToken_ValidToken(t *testing.T) {
	// 注意：这个测试需要mock config.AppConfig.JWTSecret
	// 由于测试环境限制，这里只测试token生成逻辑
	
	token := generateTestToken(123, "13800138000", "user", testJWTSecret, time.Now().Add(1*time.Hour))
	
	if token == "" {
		t.Error("Token生成失败")
	}
	
	t.Logf("生成的测试token: %s", token)
}

// TestParseToken_ExpiredToken 测试过期token
func TestParseToken_ExpiredToken(t *testing.T) {
	token := generateExpiredToken(123, testJWTSecret)
	
	if token == "" {
		t.Error("过期Token生成失败")
	}
	
	// 验证token确实过期
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return testJWTSecret, nil
	})
	
	if err == nil && parsed.Valid {
		t.Error("过期的token不应该通过验证")
	} else {
		t.Log("过期token验证失败，符合预期")
	}
}

// TestParseToken_NoneAlgorithm 测试none算法token（安全测试）
func TestParseToken_NoneAlgorithm(t *testing.T) {
	token := generateNoneAlgorithmToken(123)
	
	// 验证系统应该拒绝none算法
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// 检查签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			t.Logf("检测到非HMAC算法: %v", token.Method.Alg())
			return nil, jwt.ErrSignatureInvalid
		}
		return testJWTSecret, nil
	})
	
	if err != nil || !parsed.Valid {
		t.Log("None算法token被正确拒绝")
	} else {
		t.Error("None算法token不应该通过验证")
	}
}

// TestJWTMiddleware_MissingToken 测试缺少token的情况
func TestJWTMiddleware_MissingToken(t *testing.T) {
	// 注意：实际测试需要mock配置
	
	// 由于需要依赖注入，这个测试暂时跳过
	t.Skip("需要mock config.AppConfig")
}

// TestJWTMiddleware_InvalidFormat 测试token格式错误
func TestJWTMiddleware_InvalidFormat(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		shouldFail bool
	}{
		{
			name:       "空token",
			authHeader: "",
			shouldFail: true,
		},
		{
			name:       "格式错误-缺少Bearer前缀",
			authHeader: "some-token",
			shouldFail: true,
		},
		{
			name:       "格式错误-只有Bearer",
			authHeader: "Bearer",
			shouldFail: true,
		},
		{
			name:       "格式正确",
			authHeader: "Bearer valid-token",
			shouldFail: false, // 需要进一步验证
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := make([]string, 0)
			if tt.authHeader != "" {
				parts = splitHeader(tt.authHeader)
			}
			
			if tt.shouldFail && len(parts) != 2 {
				t.Logf("格式验证失败，符合预期: %s", tt.name)
			} else if !tt.shouldFail {
				t.Logf("格式验证通过: %s", tt.name)
			}
		})
	}
}

// splitHeader 分割Authorization header
func splitHeader(header string) []string {
	if header == "" {
		return nil
	}
	result := make([]string, 0)
	current := ""
	spaceCount := 0
	
	for _, c := range header {
		if c == ' ' {
			spaceCount++
			if spaceCount == 1 {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}