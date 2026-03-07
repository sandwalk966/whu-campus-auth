package req

type CreateMenuRequest struct {
	Name      string `json:"name" binding:"required"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Icon      string `json:"icon"`
	Sort      int    `json:"sort"`
	ParentID  uint   `json:"parent_id"`
	Type      int    `json:"type"`
	Status    int    `json:"status"`
}

type UpdateMenuRequest struct {
	ID        uint   `json:"id" binding:"required"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Icon      string `json:"icon"`
	Sort      int    `json:"sort"`
	ParentID  uint   `json:"parent_id"`
	Type      int    `json:"type"`
	Status    int    `json:"status"`
}

type MenuListRequest struct {
	Page     int `json:"page" binding:"required"`
	PageSize int `json:"page_size" binding:"required"`
}
