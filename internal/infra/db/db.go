package db

import "go.uber.org/fx"

var Module = fx.Option(
	fx.Provide(
		NewMongoDB,
		NewMongoPaymentRepository,
		NewMongoMerchantRepository,
	),
)
