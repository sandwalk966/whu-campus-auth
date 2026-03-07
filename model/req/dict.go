package req

type CreateDictRequest struct {
	Name   string         `json:"name" binding:"required"`
	Code   string         `json:"code" binding:"required"`
	Desc   string         `json:"desc"`
	Status int            `json:"status"`
	Items  []DictItemReq  `json:"items"`
}

type DictItemReq struct {
	Label  string `json:"label"`
	Value  string `json:"value"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
}

type UpdateDictRequest struct {
	ID     uint           `json:"id" binding:"required"`
	Name   string         `json:"name"`
	Code   string         `json:"code"`
	Desc   string         `json:"desc"`
	Status int            `json:"status"`
	Items  []DictItemReq  `json:"items"`
}

type DictListRequest struct {
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"page_size" binding:"required"`
	Name     string `json:"name"`
	Status   int    `json:"status"`
}
