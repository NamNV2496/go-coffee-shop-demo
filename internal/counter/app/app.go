package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/domain"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/handler/router"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/service"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/mq"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/mq/producer"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/utils"
)

type AppInterface interface {
	Start() error
	CreateOrder(ctx context.Context, req *gin.Context)
	SubmitOrder(ctx context.Context, req *gin.Context)
	// UpdateOrderStatus(ctx context.Context, req *gin.Context)
	GetItem(ctx context.Context, req *gin.Context)
	GetOrders(ctx context.Context, req *gin.Context)
}

type App struct {
	orderService service.OrderService
	grpcClient   router.ProductGRPCClient
	producer     producer.Client
}

func NewApp(
	orderService service.OrderService,
	grpcClient router.ProductGRPCClient,
	producer producer.Client,
) *App {
	return &App{
		orderService: orderService,
		grpcClient:   grpcClient,
		producer:     producer,
	}
}

func (app App) Start() error {

	// app.ConsumerHandler.StartConsumerUp(ctx)
	return nil
}

// func (app App) UpdateOrderStatus(ctx context.Context, req *gin.Context) {
// 	orderId, err := strconv.Atoi(req.Query("orderId"))
// 	if err != nil {
// 		orderId = 0
// 	}
// 	status, err := strconv.Atoi(req.Query("status"))
// 	if err != nil || status != int(mq.Canceled) {
// 		panic("Status is invalide")
// 	}
// 	app.OrderService.UpdateStatusOrder(context.Background(), int32(orderId), int32(status))
// 	req.JSON(http.StatusCreated, gin.H{"message": "Thành công"})
// }

func (app App) CreateOrder(ctx context.Context, req *gin.Context) {

	var orderList domain.OrderItemListDto

	if err := req.BindJSON(&orderList); err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
		return
	}

	if err := app.orderService.CreateOrder(
		context.Background(),
		orderList.OrderItems,
		orderList.CustomerId,
	); err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
		return
	}
	// trigger to kitchen
	data, err := json.Marshal(orderList)
	if err != nil {
		fmt.Println("Error marshall data to send to kitchen")
	}
	go func() {
		fmt.Println("Call trigger to kitchen: ", orderList)
		if err := app.producer.Produce(context.Background(), mq.TOPIC_PROCESS_COOK, data); err != nil {
			return
		}
	}()
	utils.WrapperResponse(req, http.StatusCreated, "")
}

func (app App) SubmitOrder(ctx context.Context, req *gin.Context) {
	queryId := req.Query("customerId")
	customerId, err := strconv.Atoi(queryId)
	if err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
		return
	}
	id, err := app.orderService.SubmitOrder(context.Background(), int32(customerId))
	if err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
		return
	}
	utils.WrapperResponse(req, http.StatusOK, id)
}

func (app App) GetItem(ctx context.Context, req *gin.Context) {
	// get from parameters
	id := req.Query("id")
	name := req.Query("name")
	itemType := req.Query("itemType")
	itype, _ := strconv.Atoi(itemType)
	num, _ := strconv.Atoi(id)
	pageQuery := req.Query("page")
	page, _ := strconv.Atoi(pageQuery)
	sizeQuery := req.Query("size")
	size, _ := strconv.Atoi(sizeQuery)
	items, err := app.grpcClient.GetProductByIdOrNameOrType(int32(num), name, int32(itype), int32(page), int32(size))
	if err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
		return
	}
	utils.WrapperResponse(req, http.StatusOK, items)
}

func (app App) GetOrders(ctx context.Context, req *gin.Context) {

	orderId, err := strconv.Atoi(req.Query("orderId"))
	if err != nil {
		orderId = 0
	}
	customerId, err := strconv.Atoi(req.Query("customerId"))
	if err != nil {
		customerId = 0
	}

	orderDtoRes, err := app.orderService.GetOrder(context.Background(), int32(orderId), int32(customerId))
	if err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
		return
	}
	utils.WrapperResponse(req, http.StatusOK, orderDtoRes)
}
