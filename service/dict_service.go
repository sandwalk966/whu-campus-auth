package service

import (
	"whu-campus-auth/dao"
	"whu-campus-auth/model/db"
	"whu-campus-auth/model/req"
)

type DictService struct {
	dictDAO dao.IDictDAO
}

func NewDictService(dictDAO dao.IDictDAO) *DictService {
	return &DictService{dictDAO: dictDAO}
}

func (s *DictService) CreateDict(createReq req.CreateDictRequest) error {
	dict := &db.Dict{
		Name:   createReq.Name,
		Code:   createReq.Code,
		Desc:   createReq.Desc,
		Status: createReq.Status,
	}

	if err := s.dictDAO.Create(dict); err != nil {
		return err
	}

	for _, itemReq := range createReq.Items {
		item := db.DictItem{
			DictID: dict.ID,
			Label:  itemReq.Label,
			Value:  itemReq.Value,
			Sort:   itemReq.Sort,
			Status: itemReq.Status,
		}
		s.dictDAO.GetDB().Create(&item)
	}

	return nil
}

func (s *DictService) UpdateDict(updateReq req.UpdateDictRequest) error {
	dict, err := s.dictDAO.GetByID(updateReq.ID)
	if err != nil {
		return err
	}

	if updateReq.Name != "" {
		dict.Name = updateReq.Name
	}
	if updateReq.Code != "" {
		dict.Code = updateReq.Code
	}
	dict.Desc = updateReq.Desc
	dict.Status = updateReq.Status

	if err := s.dictDAO.Update(dict); err != nil {
		return err
	}

	if updateReq.Items != nil {
		s.dictDAO.GetDB().Where("dict_id = ?", dict.ID).Delete(&db.DictItem{})
		for _, itemReq := range updateReq.Items {
			item := db.DictItem{
				DictID: dict.ID,
				Label:  itemReq.Label,
				Value:  itemReq.Value,
				Sort:   itemReq.Sort,
				Status: itemReq.Status,
			}
			s.dictDAO.GetDB().Create(&item)
		}
	}

	return nil
}

func (s *DictService) GetDictByID(id uint) (*db.Dict, error) {
	dict, err := s.dictDAO.GetByID(id)
	if err != nil {
		return nil, err
	}

	err = s.dictDAO.PreloadItems(dict)
	if err != nil {
		return nil, err
	}

	return dict, nil
}

func (s *DictService) GetDictList(page, pageSize int, name string, status int) ([]db.Dict, int64, error) {
	return s.dictDAO.GetList(page, pageSize, name, status)
}

func (s *DictService) DeleteDict(id uint) error {
	return s.dictDAO.Delete(id)
}

func (s *DictService) GetDictByCode(code string) (*db.Dict, error) {
	dict, err := s.dictDAO.GetByCode(code)
	if err != nil {
		return nil, err
	}

	err = s.dictDAO.PreloadItems(dict)
	if err != nil {
		return nil, err
	}

	return dict, nil
}
