package req

type CreateRoleRequest struct {
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Desc   string `json:"desc"`
	Status int    `json:"status"`
	MenuIDs []uint `json:"menu_ids"`
}

type UpdateRoleRequest struct {
	ID     uint   `json:"id" binding:"required"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Desc   string `json:"desc"`
	Status int    `json:"status"`
	MenuIDs []uint `json:"menu_ids"`
}

type RoleListRequest struct {
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"page_size" binding:"required"`
	Name     string `json:"name"`
	Status   int    `json:"status"`
}

type AssignRoleRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	RoleIDs []uint `json:"role_ids" binding:"required"`
}
