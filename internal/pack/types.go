package pack

type Solution struct {
	Amount       int         `json:"amount"`
	ItemsShipped int         `json:"items_shipped"`
	Overage      int         `json:"overage"`
	PackCount    int         `json:"pack_count"`
	Packs        map[int]int `json:"packs"`
}
