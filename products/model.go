package products

type Product struct {
	ProductDetails ProductDetails `json:"product"`
}

type ProductDetails struct {
	Key     string            `json:"key"`
	Version int               `json:"version"`
	Options map[string]string `json:"options"`
}

type ProductDbItem struct {
	Product
}
