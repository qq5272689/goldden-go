package models

const (
	AuthModuleLDAP = "ldap"
)

type User struct {
	ID           int64  `json:"id" gorm:"index"`                         //ID创建时不用传
	AuthModule   string `json:"auth_module"  gorm:"auth_module"`         //认证方式
	SuperAdmin   bool   `json:"super_admin" gorm:"column:super_admin"`   //是否是超级用户
	Name         string `json:"name" gorm:"column:name"`                 //用户名
	DisplayName  string `json:"display_name" gorm:"column:display_name"` //显示名称
	Role         string `json:"role" gorm:"column:role"`                 //角色
	Group        int    `json:"group" gorm:"column:group"`               //group
	Organization string `json:"organization" gorm:"column:organization"` //工作组织
	Affiliation  string `json:"affiliation" gorm:"column:affiliation"`   //工作单位
	Position     string `json:"position" gorm:"column:position"`         //职位
	Password     string `json:"password" gorm:"column:password"`         //用户密码不更新密码不用填
	Email        string `json:"email" gorm:"column:email"`               //邮箱地址
	Mobile       string `json:"mobile" gorm:"column:mobile"`             //手机号
	BaseModel
	//OldPassword string `json:"old_password" gorm:"-" swaggerignore:"true"`
}
