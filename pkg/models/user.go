package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

const (
	AuthModuleLDAP = "ldap"
)

type User struct {
	ID           int64  `json:"id" gorm:"index"`                          //ID创建时不用传
	AuthModule   string `json:"auth_module"  gorm:"auth_module"`          //认证方式
	SuperAdmin   bool   `json:"super_admin" gorm:"column:super_admin"`    //是否是超级用户
	Name         string `json:"name" gorm:"column:name"`                  //用户名
	DisplayName  string `json:"display_name" gorm:"column:display_name"`  //显示名称
	Role         string `json:"role" gorm:"column:role"`                  //角色
	Group        int    `json:"group" gorm:"column:group"`                //group
	Organization string `json:"organization" gorm:"column:organization"`  //工作组织
	Affiliation  string `json:"affiliation" gorm:"column:affiliation"`    //工作单位
	Position     string `json:"position" gorm:"column:position"`          //职位
	Password     string `json:"password" gorm:"column:password"`          //用户密码不更新密码不用填
	Email        string `json:"email" gorm:"column:email"`                //邮箱地址
	Mobile       string `json:"mobile" gorm:"column:mobile"`              //手机号
	Extend       Extend `json:"extend" gorm:"column:extend;default:'{}'"` //扩展数据
	BaseModel
	//OldPassword string `json:"old_password" gorm:"-" swaggerignore:"true"`
}

type Extend map[string]interface{}

func (t *Extend) Scan(value interface{}) error {
	str, ok := value.([]byte)
	if !ok {
		return errors.New("ErrNotString")
	}
	return json.Unmarshal(str, t)
}

func (t Extend) Value() (driver.Value, error) {
	if t == nil {
		t = map[string]interface{}{}
	}
	if db, err := json.Marshal(t); err != nil {
		return nil, err
	} else {
		return string(db), nil
	}
}

func (t Extend) GormDataType() string {
	return "text"
}
