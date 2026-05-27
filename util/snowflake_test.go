package util

import (
	"fmt"
	"sync"
	"testing"
)

func TestSnowflakeGenerate(t *testing.T) {
	sf, err := NewSnowflake(1)
	if err != nil {
		t.Fatalf("创建雪花ID失败: %v", err)
	}

	id := sf.Generate()
	if id <= 0 {
		t.Errorf("生成的ID无效: %d", id)
	}

	fmt.Printf("生成的雪花ID: %d\n", id)
}

func TestSnowflakeUniqueness(t *testing.T) {
	sf, _ := NewSnowflake(1)

	idMap := make(map[int64]bool)
	count := 10000

	for i := 0; i < count; i++ {
		id := sf.Generate()
		if idMap[id] {
			t.Errorf("发现重复ID: %d", id)
		}
		idMap[id] = true
	}

	fmt.Printf("成功生成 %d 个唯一ID\n", count)
}

func TestSnowflakeConcurrency(t *testing.T) {
	sf, _ := NewSnowflake(1)

	var wg sync.WaitGroup
	idChan := make(chan int64, 10000)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				id := sf.Generate()
				idChan <- id
			}
		}()
	}

	wg.Wait()
	close(idChan)

	idMap := make(map[int64]bool)
	count := 0
	for id := range idChan {
		count++
		if idMap[id] {
			t.Errorf("并发测试发现重复ID: %d", id)
		}
		idMap[id] = true
	}

	fmt.Printf("并发测试: 成功生成 %d 个唯一ID\n", count)
}

func TestGenerateOrderNo(t *testing.T) {
	orderNo := GenerateOrderNo()
	if orderNo == "" {
		t.Error("订单号为空")
	}

	fmt.Printf("生成的订单号: %s\n", orderNo)
}

func TestOrderNoUniqueness(t *testing.T) {
	orderNoMap := make(map[string]bool)
	count := 1000

	for i := 0; i < count; i++ {
		orderNo := GenerateOrderNo()
		if orderNoMap[orderNo] {
			t.Errorf("发现重复订单号: %s", orderNo)
		}
		orderNoMap[orderNo] = true
	}

	fmt.Printf("成功生成 %d 个唯一订单号\n", count)
}
