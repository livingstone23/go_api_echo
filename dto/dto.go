package dto

type CategoryDto struct {
	Name string `json:"name"`
}

type GenericDto struct {
	State string `json:"state"`
	Message string `json:"message"`
}

type ProductDto struct {
	Name string `json:"name"`
	Price float64 `json:"price"`
	Stock int `json:"stock"`
	Description string `json:"description"`
	CategoryID string `json:"category_id"`

}

type UserDto struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Telephone string `json:"telephone"`
	Password string `json:"password"`
}

type LoginDto struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginAnswerDto struct {
	Name string `json:"name"`
	Token string `json:"token"`
}

type ProductPictureDto struct {
	name string `json:"name"`
	ProductID string `json:"product_id"`
}


