package models

type RefStripeType struct {
	ID     int    `json:"id" example:"1"`
	Name   string `json:"name" example:"Single Solid"`
	NameTH string `json:"name_th" example:"เส้นทึบเดี่ยว"`
}

func (t *RefStripeType) TableName() string {
	return "ref_stripe_type"
}

type RefStripeColor struct {
	ID     int    `json:"id" example:"1"`
	Name   string `json:"name" example:"White"`
	NameTH string `json:"name_th" example:"เส้นสีขาว"`
}

func (t *RefStripeColor) TableName() string {
	return "ref_stripe_color"
}
