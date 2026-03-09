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
	userDAO      dao.IUserDAO
	redisService *RedisService
}

func NewUserService(userDAO dao.IUserDAO, redisService *RedisService) *UserService {
	return &UserService{
		userDAO:      userDAO,
		redisService: redisService,
	}
}

func (s *UserService) Login(loginReq req.LoginRequest) (string, *db.User, error) {
	user, err := s.userDAO.GetByUsername(loginReq.Username)
	if err != nil {
		return "", nil, errors.New("Invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		return "", nil, errors.New("Invalid username or password")
	}

	if user.Status != 1 {
		return "", nil, errors.New("Account has been disabled")
	}

	// 检查用户是否被禁用（Redis 标记）
	isDisabled, err := s.redisService.IsUserDisabled(user.ID)
	if err != nil {
		utils.LogErrorf("检查用户状态失败：%v", err)
	}
	if isDisabled {
		return "", nil, errors.New("User has been disabled or deleted")
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
		return "", nil, err
	}

	return token, user, nil
}

func (s *UserService) CreateUser(createReq req.CreateUserRequest) error {
	_, err := s.userDAO.GetByUsername(createReq.Username)
	if err == nil {
		return errors.New("Username already exists")
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

	if err := s.userDAO.Create(user); err != nil {
		return err
	}

	// 清除用户名缓存（防止缓存穿透）
	s.redisService.DeleteUserCacheByUsername(user.Username)

	return nil
}

func (s *UserService) Register(registerReq req.RegisterRequest) error {
	_, err := s.userDAO.GetByUsername(registerReq.Username)
	if err == nil {
		return errors.New("Username already exists")
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

	if err := s.userDAO.Create(user); err != nil {
		return err
	}

	// 清除用户名缓存
	s.redisService.DeleteUserCacheByUsername(user.Username)

	return nil
}

func (s *UserService) GetUserByID(id uint) (*db.User, error) {
	// 先从缓存中读取
	cacheKey := s.redisService.GetUserCacheKey(id)
	user, err := s.redisService.GetUserFromCache(id)
	if err == nil && user != nil {
		return user, nil
	}

	// 缓存未命中，从数据库读取
	user, err = s.userDAO.GetByID(id)
	if err != nil {
		return nil, err
	}

	err = s.userDAO.PreloadRoles(user)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	s.redisService.SetUserCache(user, cacheKey)

	return user, nil
}

// GetUserByIDWithCache 带缓存的用户查询（推荐方法）
func (s *UserService) GetUserByIDWithCache(id uint) (*db.User, error) {
	// 1. 先尝试从缓存获取
	cachedUser, err := s.redisService.GetCachedUserInfo(id)
	if err == nil && cachedUser != nil {
		return cachedUser, nil
	}

	// 2. 缓存未命中，从数据库获取
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// 3. 存入缓存（5 分钟）
	s.redisService.CacheUserInfo(id, user, 5*time.Minute)

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

	// 更新数据库
	if err := s.userDAO.Update(user); err != nil {
		return err
	}

	// 清除缓存
	s.redisService.DeleteUserCache(user.ID)
	s.redisService.DeleteUserCacheByUsername(user.Username)

	return nil
}

func (s *UserService) ChangePassword(userID uint, changeReq req.ChangePasswordRequest) error {
	user, err := s.userDAO.GetByID(userID)
	if err != nil {
		return errors.New("User does not exist")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changeReq.OldPassword)); err != nil {
		return errors.New("Old password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(changeReq.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	// 更新数据库
	if err := s.userDAO.Update(user); err != nil {
		return err
	}

	// 清除缓存（密码修改后需要清除）
	s.redisService.DeleteUserCache(user.ID)
	s.redisService.DeleteUserCacheByUsername(user.Username)

	return nil
}

func (s *UserService) GetUserList(page, pageSize int, username string, status int) ([]db.User, int64, error) {
	return s.userDAO.GetList(page, pageSize, username, status)
}

func (s *UserService) DeleteUser(id uint) error {
	// 先获取用户信息
	user, err := s.userDAO.GetByID(id)
	if err != nil {
		return errors.New("User does not exist")
	}

	// 删除用户（软删除）
	if err := s.userDAO.Delete(id); err != nil {
		return err
	}

	// 清除用户缓存（包括用户名缓存）
	s.redisService.DeleteUserCache(id)
	s.redisService.DeleteUserCacheByUsername(user.Username)

	return nil
}

func (s *UserService) AssignRoles(userID uint, roleIDs []uint) error {
	return s.userDAO.AssignRoles(userID, roleIDs)
}

func (s *UserService) UpdateAvatar(userID uint, avatarURL string) error {
	user, err := s.userDAO.GetByID(userID)
	if err != nil {
		return errors.New("User does not exist")
	}

	user.Avatar = avatarURL

	// 更新数据库
	if err := s.userDAO.Update(user); err != nil {
		return err
	}

	// 清除缓存
	s.redisService.DeleteUserCache(user.ID)
	s.redisService.DeleteUserCacheByUsername(user.Username)

	return nil
}
