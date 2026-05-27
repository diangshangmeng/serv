package main

import (
	"fmt"

	"voucher-platform/config"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		return
	}

	err = config.InitDB()
	if err != nil {
		fmt.Printf("初始化数据库失败: %v\n", err)
		return
	}
	defer config.DB.Close()

	fmt.Println("开始数据库迁移...")

	// 尝试删除shipping_address_id字段（不管是否存在）
	err = config.DB.Exec("ALTER TABLE orders DROP COLUMN shipping_address_id").Error
	if err != nil {
		fmt.Printf("⚠ 删除 shipping_address_id 字段: %v\n", err)
	} else {
		fmt.Println("✓ 成功删除 orders 表中的 shipping_address_id 字段")
	}

	// 尝试删除shipping_address字段（不管是否存在）
	err = config.DB.Exec("ALTER TABLE orders DROP COLUMN shipping_address").Error
	if err != nil {
		fmt.Printf("⚠ 删除 shipping_address 字段: %v\n", err)
	} else {
		fmt.Println("✓ 成功删除 orders 表中的 shipping_address 字段")
	}

	fmt.Println("\n数据库迁移完成！代金券交易不再需要地址信息。")
	fmt.Println("现在可以重新测试创建订单功能了。")
}
