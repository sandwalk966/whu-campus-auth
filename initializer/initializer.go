package initializer

import (
	"whu-campus-auth/dao"
	dbModel "whu-campus-auth/model/db"
	"whu-campus-auth/model/req"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitDictData 初始化字典数据
// 在项目启动时调用，自动创建常用字典
func InitDictData(db *gorm.DB) {
	dictDAO := dao.NewDictDAO(db)
	
	dicts := []req.CreateDictRequest{
		{
			Name: "Gender Dictionary",
			Code: "gender",
			Desc: "System gender options",
			Status: 1,
			Items: []req.DictItemReq{
				{Label: "Male", Value: "1", Sort: 1, Status: 1},
				{Label: "Female", Value: "2", Sort: 2, Status: 1},
				{Label: "Unknown", Value: "0", Sort: 3, Status: 1},
			},
		},
		{
			Name: "User Status",
			Code: "user_status",
			Desc: "User account status",
			Status: 1,
			Items: []req.DictItemReq{
				{Label: "Active", Value: "1", Sort: 1, Status: 1},
				{Label: "Disabled", Value: "0", Sort: 2, Status: 1},
			},
		},
		{
			Name: "Menu Type",
			Code: "menu_type",
			Desc: "System menu type",
			Status: 1,
			Items: []req.DictItemReq{
				{Label: "Directory", Value: "1", Sort: 1, Status: 1},
				{Label: "Menu", Value: "2", Sort: 2, Status: 1},
				{Label: "Button", Value: "3", Sort: 3, Status: 1},
			},
		},
		{
			Name: "Role Status",
			Code: "role_status",
			Desc: "Role status",
			Status: 1,
			Items: []req.DictItemReq{
				{Label: "Active", Value: "1", Sort: 1, Status: 1},
				{Label: "Disabled", Value: "0", Sort: 2, Status: 1},
			},
		},
	}

	for _, createReq := range dicts {
		// 检查字典是否已存在
		_, err := dictDAO.GetByCode(createReq.Code)
		if err == nil {
			// 字典已存在，跳过
			continue
		}

		// 创建字典
		dict := &dbModel.Dict{
			Name:   createReq.Name,
			Code:   createReq.Code,
			Desc:   createReq.Desc,
			Status: createReq.Status,
		}

		if err := dictDAO.Create(dict); err != nil {
			zap.L().Error("创建字典失败", zap.String("code", createReq.Code), zap.Error(err))
			continue
		}

		// 创建字典项
		for _, itemReq := range createReq.Items {
			item := dbModel.DictItem{
				DictID: dict.ID,
				Label:  itemReq.Label,
				Value:  itemReq.Value,
				Sort:   itemReq.Sort,
				Status: itemReq.Status,
			}
			if err := dictDAO.GetDB().Create(&item).Error; err != nil {
				zap.L().Error("创建字典项失败", zap.String("dict", createReq.Code), zap.Error(err))
			}
		}

		zap.L().Info("字典初始化成功", zap.String("code", createReq.Code))
	}
}
