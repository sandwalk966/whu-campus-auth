package service

import (
	"whu-campus-auth/dao"
	"whu-campus-auth/model/db"
	"whu-campus-auth/model/req"
)

type MenuService struct {
	menuDAO dao.IMenuDAO
}

func NewMenuService(menuDAO dao.IMenuDAO) *MenuService {
	return &MenuService{menuDAO: menuDAO}
}

func (s *MenuService) CreateMenu(createReq req.CreateMenuRequest) error {
	menu := &db.Menu{
		Name:      createReq.Name,
		Path:      createReq.Path,
		Component: createReq.Component,
		Icon:      createReq.Icon,
		Sort:      createReq.Sort,
		ParentID:  createReq.ParentID,
		Type:      createReq.Type,
		Status:    createReq.Status,
	}

	return s.menuDAO.Create(menu)
}

func (s *MenuService) UpdateMenu(updateReq req.UpdateMenuRequest) error {
	menu, err := s.menuDAO.GetByID(updateReq.ID)
	if err != nil {
		return nil
	}

	if updateReq.Name != "" {
		menu.Name = updateReq.Name
	}
	if updateReq.Path != "" {
		menu.Path = updateReq.Path
	}
	if updateReq.Component != "" {
		menu.Component = updateReq.Component
	}
	if updateReq.Icon != "" {
		menu.Icon = updateReq.Icon
	}
	menu.Sort = updateReq.Sort
	menu.ParentID = updateReq.ParentID
	menu.Type = updateReq.Type
	menu.Status = updateReq.Status

	return s.menuDAO.Update(menu)
}

func (s *MenuService) GetMenuByID(id uint) (*db.Menu, error) {
	return s.menuDAO.GetByID(id)
}

func (s *MenuService) GetMenuList(page, pageSize int) ([]db.Menu, int64, error) {
	return s.menuDAO.GetList(page, pageSize, "", 0)
}

func (s *MenuService) GetMenuTree() ([]db.Menu, error) {
	return s.menuDAO.GetTree()
}

func (s *MenuService) DeleteMenu(id uint) error {
	return s.menuDAO.Delete(id)
}

func (s *MenuService) GetMenusByRoleID(roleID uint) ([]db.Menu, error) {
	return s.menuDAO.GetByRoleID(roleID)
}
