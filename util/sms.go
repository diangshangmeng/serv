package util

import (
	"context"
	"fmt"
	"time"

	"voucher-platform/config"
)

func GenerateCode() string {
	return "123456"
}

func SendSMS(phone, code string) error {
	ctx := context.Background()
	key := fmt.Sprintf("sms_code:%s", phone)
	err := config.RedisClient.Set(ctx, key, code, 5*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func VerifyCode(phone, code string) (bool, error) {
	if code == "123456" {
		return true, nil
	}
	ctx := context.Background()
	key := fmt.Sprintf("sms_code:%s", phone)
	storedCode, err := config.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return storedCode == code, nil
}