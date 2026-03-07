package service

import (
	"errors"
	"whu-campus-auth/dao"
	"whu-campus-auth/model/db"
	"whu-campus-auth/model/req"
)

type RoleService struct {
	roleDAO dao.IRoleDAO
}

func NewRoleService(roleDAO dao.IRoleDAO) *RoleService {
	return &RoleService{roleDAO: roleDAO}
}

func (s *RoleService) CreateRole(createReq req.CreateRoleRequest) error {
	_, err := s.roleDAO.GetByCode(createReq.Code)
	if err == nil {
		return errors.New("角色编码已存在")
	}

	role := &db.Role{
		Name:   createReq.Name,
		Code:   createReq.Code,
		Desc:   createReq.Desc,
		Status: createReq.Status,
	}

	if err := s.roleDAO.Create(role); err != nil {
		return err
	}

	if len(createReq.MenuIDs) > 0 {
		return s.roleDAO.AssignMenus(role.ID, createReq.MenuIDs)
	}

	return nil
}

func (s *RoleService) UpdateRole(updateReq req.UpdateRoleRequest) error {
	role, err := s.roleDAO.GetByID(updateReq.ID)
	if err != nil {
		return errors.New("角色不存在")
	}

	if updateReq.Name != "" {
		role.Name = updateReq.Name
	}
	if updateReq.Code != "" {
		role.Code = updateReq.Code
	}
	role.Desc = updateReq.Desc
	role.Status = updateReq.Status

	if err := s.roleDAO.Update(role); err != nil {
		return err
	}

	if updateReq.MenuIDs != nil {
		return s.roleDAO.AssignMenus(role.ID, updateReq.MenuIDs)
	}

	return nil
}

func (s *RoleService) GetRoleByID(id uint) (*db.Role, error) {
	role, err := s.roleDAO.GetByID(id)
	if err != nil {
		return nil, err
	}

	err = s.roleDAO.PreloadMenus(role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) GetRoleList(page, pageSize int, name string, status int) ([]db.Role, int64, error) {
	return s.roleDAO.GetList(page, pageSize, name, status)
}

func (s *RoleService) DeleteRole(id uint) error {
	return s.roleDAO.Delete(id)
}

func (s *RoleService) GetAllRoles() ([]db.Role, error) {
	return s.roleDAO.GetAll()
}
