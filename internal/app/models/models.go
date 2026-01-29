package models

type OrderReq struct {
	CouponCode string      `json:"couponCode" description:"Optional promo code applied to the order"`
	Items      []OrderItem `json:"items" binding:"required"`
}

type ApiResponse struct {
	Code    int    `json:"code" format:"int32"`
	Type    string `json:"type"`
	Message string `json:"message"`
}
