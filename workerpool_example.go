package main

import (
	"fmt"
	"math/rand"
	"time"

	"demo-go/workerpool"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	wp := workerpool.NewWorkerPool(10)
	go wp.Start()

	// test1: 測試是否限制最大工作數量
	orders := getFakeData()
	for _, order := range orders {
		o := order // 避免閉包
		wp.Submit(func() {
			time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
			fmt.Printf("Order ID %d, Customer: %s, Amount: $%.2f\n", o.ID, o.Customer, o.Amount)
		})
	}
	// 等待所有任务完成
	// wp.Wait()

	wp.Stop()

	// test2: 測試 SubmitWait
	// fmt.Println("開始測試 SubmitWait...")
	// wp.SubmitWait(func() {
	// 	fmt.Println("這是 SubmitWait 測試的工作，會等待完成後再繼續執行其他操作。")
	// 	time.Sleep(2 * time.Second)
	// 	fmt.Println("SubmitWait 測試完成。")
	// })

	// test3: 測試 timeout 情況
	// fmt.Println("開始測試 timeout 情況...")
	// for i := 0; i < 5; i++ {
	// 	wp.Submit(func() {
	// 		time.Sleep(1 * time.Second)
	// 		fmt.Println("這個工作應該會被處理，但處理時間比較長。")
	// 	})
	// }

	// test4: 測試取消任務
	// fmt.Println("開始測試取消任務...")
	// go func() {
	// 	time.Sleep(10 * time.Second)
	// 	wp.Stop()
	// }()

	// for i := 0; i < len(orders); i++ {
	// 	o := orders[i]
	// 	wp.Submit(func() {
	// 		fmt.Printf("Processing Order ID %d\n", o.ID)
	// 		time.Sleep(1 * time.Second)
	// 	})

	// 	if i == 50 {
	// 		wp.Stop() // 处理到第51个订单时优雅关闭
	// 		break
	// 	}
	// }

	// 等待所有工作完成

	fmt.Println("所有測試已完成。")
}

// order 結構
type order struct {
	ID       int
	Customer string
	Amount   float64
}

// 模擬生成假數據
func getFakeData() []order {
	// 模拟 xx 个订单处理任务
	mockNum := 100
	orders := make([]order, mockNum)

	// 使用循环生成 xx 个订单
	for i := 0; i < mockNum; i++ {
		orders[i] = order{
			ID:       i + 1,
			Customer: fmt.Sprintf("Customer %d", i+1),
			Amount:   float64(rand.Intn(1000) + 1), // 随机金额在 1 到 1000 之间
		}
	}
	return orders
}
