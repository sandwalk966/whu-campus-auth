package initializer

import "whu-campus-auth/dao"

// initDAO 初始化所有 DAO
func (d *Dependencies) initDAO() {
	d.UserDAO = dao.NewUserDAO(d.DB)
	d.RoleDAO = dao.NewRoleDAO(d.DB)
	d.MenuDAO = dao.NewMenuDAO(d.DB)
	d.DictDAO = dao.NewDictDAO(d.DB)
}
