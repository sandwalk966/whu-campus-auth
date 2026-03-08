package service

import (
	"errors"
	"time"
	"whu-campus-auth/dao"
	"whu-campus-auth/model/db"
	"whu-campus-auth/model/req"
	"whu-campus-auth/utils"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userDAO dao.IUserDAO
}

func NewUserService(userDAO dao.IUserDAO) *UserService {
	return &UserService{userDAO: userDAO}
}

func (s *UserService) Login(loginReq req.LoginRequest) (string, error) {
	user, err := s.userDAO.GetByUsername(loginReq.Username)
	if err != nil {
		return "", errors.New("用户名或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		return "", errors.New("用户名或密码错误")
	}

	if user.Status != 1 {
		return "", errors.New("账号已被禁用")
	}

	j := utils.NewJWT()
	claims := utils.Claims{
		ID:       user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpiresTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "whu-campus-auth",
		},
	}

	token, err := j.CreateToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) CreateUser(createReq req.CreateUserRequest) error {
	_, err := s.userDAO.GetByUsername(createReq.Username)
	if err == nil {
		return errors.New("用户名已存在")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &db.User{
		Username: createReq.Username,
		Password: string(hashedPassword),
		Email:    createReq.Email,
		Phone:    createReq.Phone,
		Status:   createReq.Status,
	}

	return s.userDAO.Create(user)
}

func (s *UserService) Register(registerReq req.RegisterRequest) error {
	_, err := s.userDAO.GetByUsername(registerReq.Username)
	if err == nil {
		return errors.New("用户名已存在")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &db.User{
		Username: registerReq.Username,
		Password: string(hashedPassword),
		Nickname: registerReq.Nickname,
		Email:    registerReq.Email,
		Phone:    registerReq.Phone,
		Status:   1,
	}

	return s.userDAO.Create(user)
}

func (s *UserService) GetUserByID(id uint) (*db.User, error) {
	user, err := s.userDAO.GetByID(id)
	if err != nil {
		return nil, err
	}

	err = s.userDAO.PreloadRoles(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(updateReq req.UpdateUserRequest) error {
	user, err := s.userDAO.GetByID(updateReq.ID)
	if err != nil {
		return errors.New("用户不存在")
	}

	user.Nickname = updateReq.Nickname
	user.Avatar = updateReq.Avatar
	user.Email = updateReq.Email
	user.Phone = updateReq.Phone
	user.Gender = updateReq.Gender
	user.Status = updateReq.Status

	return s.userDAO.Update(user)
}

func (s *UserService) ChangePassword(userID uint, changeReq req.ChangePasswordRequest) error {
	user, err := s.userDAO.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changeReq.OldPassword)); err != nil {
		return errors.New("原密码错误")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(changeReq.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userDAO.Update(user)
}

func (s *UserService) GetUserList(page, pageSize int, username string, status int) ([]db.User, int64, error) {
	return s.userDAO.GetList(page, pageSize, username, status)
}

func (s *UserService) DeleteUser(id uint) error {
	return s.userDAO.Delete(id)
}

func (s *UserService) AssignRoles(userID uint, roleIDs []uint) error {
	return s.userDAO.AssignRoles(userID, roleIDs)
}

func (s *UserService) UpdateAvatar(userID uint, avatarURL string) error {
	user, err := s.userDAO.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	user.Avatar = avatarURL
	return s.userDAO.Update(user)
}
