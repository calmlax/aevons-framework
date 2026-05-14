package auth

import "encoding/json"

// AuthCodeInfo 授权码在 Redis 中存储的结构，关联 userId、clientId 和授权范围。
type AuthCodeInfo struct {
	UserId   int64    `json:"user_id,string"`
	ClientId string   `json:"client_id"`
	Scopes   []string `json:"scopes"`
}

// LoginRequest 所有登录模式的请求 DTO。
// client_id 和 client_secret 通过 Authorization: Basic base64(client_id:client_secret) 传递，不在请求体中。
type LoginRequest struct {
	GrantType    string `json:"grant_type" binding:"required"` // password | authorization_code | email | email | client_credentials | refresh_token
	ClientId     string `json:"-"`                             // 从 Basic Auth 头解析，不参与 JSON 绑定
	ClientSecret string `json:"-"`                             // 从 Basic Auth 头解析，不参与 JSON 绑定
	Username     string `json:"username"`
	Password     string `json:"password"`
	KeyId        string `json:"key_id"` // RSA 私钥标识，关联加密传输的密码
	Email        string `json:"email"`
	Code         string `json:"code"`          // 邮箱验证码或授权码
	RefreshToken string `json:"refresh_token"` // refresh_token 模式专用
	RedirectURI  string `json:"redirect_uri"`  // 授权码模式可选，用于校验回调地址

	// 网络原数据（在 Handler 内部获取填充）
	ClientIP  string `json:"-"`
	UserAgent string `json:"-"`
}

// UpdatePasswordRequest 用户自行修改密码请求
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NextPassword    string `json:"nextPassword" binding:"required"`
}

// PublicKeyResponse RSA 公钥获取响应
type PublicKeyResponse struct {
	KeyId     string `json:"key_id"`     // 私钥唯一标识，登录时需带回
	PublicKey string `json:"public_key"` // RSA公钥（PEM格式），前端用于加密密码
}

// TokenPair 颁发给客户端的令牌对响应（符合 RFC 6749）。
type TokenPair struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`         // "Bearer"
	ExpiresIn        int64  `json:"expires_in"`         // access token 有效期（秒），RFC 6749 标准字段
	RefreshExpiresIn int64  `json:"refresh_expires_in"` // refresh token 有效期（秒）
	Scope            string `json:"scope,omitempty"`    // 实际授权的 scope，空格分隔
}

// LoginUser 存储于 Redis 中的用户会话对象，以 Access Token 为键。
type LoginUser struct {
	UserId       int64    `json:"user_id,string"`
	Username     string   `json:"username"`
	Nickname     string   `json:"nickname"`
	Email        string   `json:"email"`
	Avatar       string   `json:"avatar"`
	Roles        []Role   `json:"roles"`
	Depts        []Dept   `json:"depts"`
	Permissions  []string `json:"permissions"`
	RefreshToken string   `json:"refresh_token,omitempty"`
	ClientId     string   `json:"client_id"`
}

type Role struct {
	Id        int64   `json:"id,string"`
	RoleKey   string  `json:"role_key"`
	DataScope int16   `json:"data_scope"`
	DeptIds   []int64 `json:"dept_ids,string"`
	RoleName  string  `json:"role_name"`
}

type Dept struct {
	DeptId int64 `json:"dept_id,string"`
	PostId int64 `json:"post_id,string"`
}

type Meta struct {
	Title       string   `json:"title,omitempty"`
	Icon        string   `json:"icon,omitempty"`
	TitleKey    string   `json:"titleKey,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Hidden      bool     `json:"hidden,omitempty"`
	ActiveMenu  string   `json:"activeMenu,omitempty"`
	IsFrame     bool     `json:"isFrame,omitempty"`
}

type Menu struct {
	Key       string `json:"key,omitempty"`
	Path      string `json:"path,omitempty"`
	Component string `json:"component,omitempty"`
	Meta      *Meta  `json:"meta,omitempty"`
	Query     string `json:"query,omitempty"`
	Children  []Menu `json:"children,omitempty"`
}

// MarshalJSON 将 LoginUser 序列化为 JSON 字节。
// 确保 Roles 和 Permissions 为 nil 时序列化为空数组而非 null。
func (u *LoginUser) MarshalJSON() ([]byte, error) {
	type Alias LoginUser
	alias := (*Alias)(u)
	// 归一化 nil 切片为空切片，保证 JSON 输出为 [] 而非 null
	if alias.Roles == nil {
		alias.Roles = []Role{}
	}
	if alias.Depts == nil {
		alias.Depts = []Dept{}
	}
	if alias.Permissions == nil {
		alias.Permissions = []string{}
	}
	return json.Marshal(alias)
}

// UnmarshalJSON 将 JSON 字节反序列化为 LoginUser。
func (u *LoginUser) UnmarshalJSON(data []byte) error {
	type Alias LoginUser
	aux := (*Alias)(u)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	// 归一化 nil 切片为空切片
	if aux.Roles == nil {
		aux.Roles = []Role{}
	}
	if aux.Depts == nil {
		aux.Depts = []Dept{}
	}
	if aux.Permissions == nil {
		aux.Permissions = []string{}
	}
	return nil
}
