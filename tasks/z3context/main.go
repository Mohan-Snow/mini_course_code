package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// don't update TestRandFunc, it's for tests
var TestRandFunc func(n int64) int64 = rand.Int63n

var ErrService = errors.New("service error")
var ErrTimeout = errors.New("timeout ended")

type order struct {
	value   int64
	id      int64
	fee     int64
	itemIDs []int64
}

func (o order) String() string {
	return fmt.Sprintf(" id: %d, value: %d, fee: %d, itemIDs: %v", o.id, o.value, o.fee, o.itemIDs)
}

type Service struct {
	idsToFail map[int64]struct{}
}

func New(idsThatReturnError map[int64]struct{}) *Service {
	if idsThatReturnError == nil {
		idsThatReturnError = make(map[int64]struct{})
	}
	return &Service{
		idsToFail: idsThatReturnError,
	}
}

// Возвращает заказ и возможную ошибку. Внутри Sleep эмитирует выполнение полезной нагрузки.
func (s *Service) getOrderByID(id int64) (*order, error) {
	randMultiplier := TestRandFunc(10)
	fmt.Printf("Sleep=%d ", randMultiplier)
	time.Sleep(time.Duration(randMultiplier) * time.Second)
	//time.Sleep(3 * time.Second)

	// some user ids fails
	if _, ok := s.idsToFail[id]; ok {
		return nil, ErrService
	}

	return &order{
		value: int64(rand.Intn(1000)),
		id:    id,
		fee:   int64(rand.Intn(200)),
		itemIDs: func() []int64 {
			res := make([]int64, 10)
			for i := 0; i < 10; i++ {
				res[i] = int64(i * rand.Intn(30))
			}
			return res
		}(),
	}, nil
}

type orderResponse struct {
	order *order
	err   error
}

// Ожидает ответа или возвращение ошибки, если ответ не был получен вовремя.
func (s *Service) getOrderByIDWrapper(contextWithTimeout context.Context, id int64) (*order, error) {
	//contextWithTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//defer cancel()
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//}()
	//wg.Wait()
	//orderChannel := make(chan *order)
	//errChannel := make(chan error)
	orderRespChannel := make(chan orderResponse, 1)

	go func() {
		order, err := s.getOrderByID(id)
		orderRespChannel <- orderResponse{order: order, err: err}
	}()

	select {
	case <-contextWithTimeout.Done():
		return nil, ErrTimeout
	case resp := <-orderRespChannel:
		return resp.order, resp.err
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	usersWithErr := map[int64]struct{}{
		1: {},
		2: {},
		3: {},
	}

	service := New(usersWithErr)

	// user with err
	//fmt.Println(service.getOrderByID(1))
	// user with err
	//fmt.Println(service.getOrderByID(2))
	// user with err
	//fmt.Println(service.getOrderByID(3))

	//fmt.Println(service.getOrderByID(4))
	//fmt.Println(service.getOrderByID(5))

	for i := 0; i < 10; i++ {
		i := i // "shadowing"
		// Pass a context with a timeout to tell a blocking function that it
		// should abandon its work after the timeout elapses.
		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
		result, err := service.getOrderByIDWrapper(ctx, int64(i))
		fmt.Println(result, err)
		cancel()
	}

}
