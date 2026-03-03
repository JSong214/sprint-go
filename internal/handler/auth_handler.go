package handler

import (
	"errors"
	"strings"

	"github.com/JSong214/sprint-go/internal/handler/response"
	"github.com/JSong214/sprint-go/internal/service"
	"github.com/JSong214/sprint-go/internal/service/dto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// AuthHandler 认证相关的 HTTP 处理器
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler 创建 AuthHandler 实例
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 用户注册处理器
// @Summary 用户注册
// @Description 创建新用户账户
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "注册信息"
// @Success 200 {object} response.ResSuccess{data=dto.UserResponse} "注册成功"
// @Failure 400 {object} response.ResError "参数验证失败或用户已存在"
// @Failure 500 {object} response.ResError "服务器内部错误"
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	// 1. 解析请求体
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数绑定失败（如 JSON 格式错误）
		h.handleValidationError(c, err)
		return
	}

	// 2. 调用 Service 层业务逻辑
	userResp, err := h.authService.Register(&req)
	if err != nil {
		// 处理业务错误
		h.handleServiceError(c, err)
		return
	}

	// 3. 返回成功响应
	response.SuccessWithMessage(c, "注册成功", userResp)
}

// Login 用户登录处理器
// @Summary 用户登录
// @Description 使用用户名/邮箱和密码登录
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登录信息"
// @Success 200 {object} response.ResSuccess{data=dto.LoginResponse} "登录成功"
// @Failure 400 {object} response.ResError "参数验证失败"
// @Failure 401 {object} response.ResError "用户名或密码错误"
// @Failure 500 {object} response.ResError "服务器内部错误"
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	// 1. 解析请求体
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleValidationError(c, err)
		return
	}

	// 2. 调用 Service 层业务逻辑
	loginResp, err := h.authService.Login(&req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	// 3. 返回成功响应
	response.SuccessWithMessage(c, "登录成功", loginResp)
}

// Logout 用户登出处理器
// @Summary 用户登出
// @Description 退出登录（清除会话/令牌）
// @Tags Auth
// @Security BearerAuth
// @Success 200 {object} response.ResSuccess "登出成功"
// @Failure 401 {object} response.ResError "未认证或令牌无效"
// @Failure 500 {object} response.ResError "服务器内部错误"
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// TODO: 从 JWT token 中解析 userID（需要 JWT 中间件）
	// 示例：userID := c.GetInt64("user_id")

	// MVP 阶段简化实现：前端删除 token 即可
	// 后端暂不需要额外操作（JWT 自动过期）

	// 调用 Service 层（当前为空实现）
	// err := h.authService.Logout(userID)
	// if err != nil {
	// 	h.handleServiceError(c, err)
	// 	return
	// }

	response.SuccessWithMessage(c, "登出成功", nil)
}

// handleValidationError 处理参数验证错误
func (h *AuthHandler) handleValidationError(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		// 提取字段级别的验证错误
		details := make(map[string]string)
		for _, fieldError := range validationErrors {
			field := strings.ToLower(fieldError.Field())
			details[field] = h.getValidationErrorMessage(fieldError)
		}
		response.BadRequest(c, 40002, "参数验证失败", details)
		return
	}

	// JSON 格式错误或其他绑定错误
	response.BadRequest(c, 40002, "参数验证失败", map[string]string{
		"error": err.Error(),
	})
}

// handleServiceError 处理 Service 层返回的业务错误
func (h *AuthHandler) handleServiceError(c *gin.Context, err error) {
	// 根据不同的错误类型返回对应的响应
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		// 用户不存在（登录场景）
		response.Unauthorized(c, "用户名或密码错误")
	case strings.Contains(err.Error(), "already exists"):
		// 用户名或邮箱已存在（注册场景）
		response.BadRequest(c, 40005, "用户名或邮箱已被注册", nil)
	case strings.Contains(err.Error(), "invalid credentials"):
		// 密码错误（登录场景）
		response.Unauthorized(c, "用户名或密码错误")
	default:
		// 其他服务器内部错误
		response.InternalServerError(c, "服务器内部错误")
	}
}

// getValidationErrorMessage 根据验证标签返回友好的错误消息
func (h *AuthHandler) getValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " 不能为空"
	case "email":
		return "邮箱格式不正确"
	case "min":
		return fe.Field() + " 长度至少为 " + fe.Param() + " 位"
	case "max":
		return fe.Field() + " 长度不能超过 " + fe.Param() + " 位"
	default:
		return fe.Field() + " 验证失败"
	}
}
