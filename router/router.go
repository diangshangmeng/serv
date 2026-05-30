package router

import (
	"voucher-platform/config"
	"voucher-platform/controller"
	"voucher-platform/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	imgURL := config.AppConfig.ImageUploadPath
	r.Use(middleware.CORSMiddleware())

	r.Static("/uploads", imgURL)

	auth := r.Group("/api/auth")
	auth.Use(middleware.RateLimitMiddleware(30, 60))
	{
		auth.POST("/send-code", controller.SendCode)
		auth.POST("/register", controller.Register)
		auth.POST("/login", controller.Login)
	}

	upload := r.Group("/api/upload")
	upload.Use(middleware.RateLimitMiddleware(30, 60))
	{
		upload.POST("/image", controller.UploadImage)
		upload.POST("/auth-image", middleware.JWTMiddleware(), middleware.UserStatusCheckMiddleware(), controller.UploadAuthImage)
		upload.POST("/order-image", middleware.JWTMiddleware(), middleware.UserStatusCheckMiddleware(), controller.UploadOrderImage)
	}

	r.POST("/api/admin/login", middleware.AdminLoginRateLimit(), controller.AdminLogin)

	admin := r.Group("/api/admin")
	admin.Use(middleware.AdminRateLimit())
	admin.Use(middleware.AdminJWTMiddleware())
	{
		admin.POST("/password/update", controller.UpdateAdminPassword)
		admin.GET("/user/pending", controller.GetPendingUsers)
		admin.POST("/user/audit", controller.AuditUser)
		admin.GET("/user/list", controller.GetAllUsers)
		admin.GET("/user/detail", controller.GetUserDetail)
		admin.GET("/user/search", controller.SearchUsers)
		admin.POST("/user/ban", controller.BanUser)
		admin.POST("/user/unban", controller.UnbanUser)
		admin.POST("/user/setStatus", controller.SetUserStatus)
		admin.POST("/user/reset-pwd", controller.ResetUserPassword)
		admin.PUT("/user/update", controller.UpdateUser)
		admin.GET("/order/list", controller.GetAllOrders)
		admin.GET("/order/detail", controller.AdminGetOrderDetail)
		admin.GET("/order/stats", controller.GetOrderStats)
		admin.POST("/order/cancel", controller.AdminCancelOrder)
		admin.POST("/order/seller-cancel", controller.AdminSellerCancelOrder)
		admin.POST("/order/seller-confirm", controller.AdminSellerConfirmOrder)
		admin.GET("/coupon/list", controller.GetCouponList)
		admin.GET("/coupon/detail", controller.GetCouponDetail)

		productImage := admin.Group("/product-image")
		{
			productImage.POST("/batch-upload", controller.BatchUploadProductImages)
			productImage.GET("/list", controller.GetProductImageList)
			productImage.GET("/detail/:id", controller.GetProductImageDetail)
			productImage.PUT("/update/:id", controller.UpdateProductImage)
			productImage.PUT("/mark-used/:id", controller.MarkProductImageAsUsed)
			productImage.DELETE("/delete/:id", controller.DeleteProductImage)
			productImage.POST("/batch-delete", controller.BatchDeleteProductImages)
		}

		product := admin.Group("/product")
		{
			product.POST("", controller.CreateProduct)
			product.GET("/list", controller.GetProductListForDashboard)
			product.GET("/:id", controller.GetProductDetail)
			product.PUT("/:id", controller.UpdateProduct)
			product.DELETE("/:id", controller.DeleteProduct)
			product.POST("/:id/publish", controller.PublishProduct)
			product.POST("/:id/unpublish", controller.UnpublishProduct)
			product.POST("/:id/unlock", controller.UnlockProduct)
		}

		paymentCodeImage := admin.Group("/payment-code-image")
		{
			paymentCodeImage.POST("", controller.UploadPaymentCodeImage)
			paymentCodeImage.GET("/list", controller.GetPaymentCodeImageList)
			paymentCodeImage.GET("/:id", controller.GetPaymentCodeImageDetail)
			paymentCodeImage.PUT("/:id", controller.UpdatePaymentCodeImage)
			paymentCodeImage.DELETE("/:id", controller.DeletePaymentCodeImage)
		}
	}

	user := r.Group("/api/user")
	user.Use(middleware.RateLimitMiddleware(30, 60))
	user.Use(middleware.JWTMiddleware())
	user.Use(middleware.UserStatusCheckMiddleware())
	{
		user.POST("/auth/submit", controller.SubmitAuth)
		user.GET("/auth/status", controller.GetAuthStatus)
	}

	app := r.Group("/api/app")
	app.Use(middleware.RateLimitMiddleware(30, 60))
	app.Use(middleware.JWTMiddleware())
	app.Use(middleware.UserStatusCheckMiddleware())
	{
		app.GET("/product/list", controller.GetAppProductList)
		app.GET("/product/my", controller.GetMyProducts)
		app.GET("/product/:id", controller.GetAppProductDetail)
		app.POST("/product/:id/publish", controller.AppPublishProduct)
		app.POST("/order/create", controller.CreateProductOrder)
		app.POST("/order/place", controller.PlaceOrder)
		app.POST("/order/pay", controller.CompletePayment)
		app.POST("/order/confirm", controller.ConfirmProduct)
		app.POST("/order/cancel", controller.CancelProduct)
		app.GET("/order/status", controller.GetOrderStatus)
		app.GET("/order/buyer", controller.GetAppBuyerOrders)
		app.GET("/order/seller", controller.GetAppSellerOrders)
		app.GET("/order/all", controller.GetAppAllOrders)
	}

	coupon := r.Group("/api/coupon")
	coupon.Use(middleware.RateLimitMiddleware(30, 60))
	coupon.Use(middleware.JWTMiddleware())
	coupon.Use(middleware.UserStatusCheckMiddleware())
	{
		coupon.POST("/add", controller.AddCoupon)
		coupon.GET("/my", controller.GetMyCoupons)
	}

	listing := r.Group("/api/listing")
	listing.Use(middleware.RateLimitMiddleware(30, 60))
	listing.Use(middleware.JWTMiddleware())
	listing.Use(middleware.UserStatusCheckMiddleware())
	{
		listing.POST("/create", controller.CreateListing)
		listing.GET("/list", controller.GetListingList)
		listing.GET("/my", controller.GetMyListings)
	}

	order := r.Group("/api/order")
	order.Use(middleware.RateLimitMiddleware(30, 60))
	order.Use(middleware.JWTMiddleware())
	order.Use(middleware.UserStatusCheckMiddleware())
	{
		order.POST("/create", controller.CreateOrder)
		order.POST("/upload-voucher", controller.UploadVoucher)
		order.POST("/confirm", controller.ConfirmOrder)
		order.POST("/seller-cancel", controller.SellerCancelOrder)
		order.GET("/my-buy", controller.GetMyBuyOrders)
		order.GET("/my-sell", controller.GetMySellOrders)
		order.GET("/detail", controller.GetOrderDetail)
	}

	public := r.Group("/api")
	public.Use(middleware.RateLimitMiddleware(30, 60))
	{
		public.GET("/city/list", controller.GetCityList)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
