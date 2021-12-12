package entity

type Person struct {
	Name    string `json:"name" binding:"required"`
	Age     int    `json:"age" binding:"gte=1,lte=130"`
	Company string `json:"company"`
	Email   string `json:"email" binding:"email"` //用validate无效
}

// Order 资源需求
type Order struct {
	Id          int    `json:"id" binding:"required"`
	Partner     Person `json:"partner" binding:"required"`
	Description string `json:"description"`
}
