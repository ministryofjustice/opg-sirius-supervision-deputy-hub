package model

type OrderStatus struct {
	Handle      string `json:"handle"`
	Incomplete  string `json:"incomplete"`
	Category    string `json:"category"`
	Complete    string `json:"complete"`
	StatusCount int
}

func (os OrderStatus) IsSelected(selectedOrderStatuses []string) bool {
	for _, selectedOrderStatus := range selectedOrderStatuses {
		if os.Handle == selectedOrderStatus {
			return true
		}
	}
	return false
}
