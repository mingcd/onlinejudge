package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
)

var (
	passwordError       error
	UsernameError       error
	emailError          error
	nameError           error
	collegeError        error
	uniqueUsernameError error
	uniqueEmailError    error
	cost                int
)

func init() {
	cost = 10
	passwordError = errors.New("Invalid password")
	UsernameError = errors.New("Invalid username. Only letters, numbers, underscores, and period are allowed")
	emailError = errors.New("Invalid email")
	nameError = errors.New("Invalid name")
	collegeError = errors.New("Invalid college name")
	uniqueUsernameError = errors.New("Username already in use, please choose another")
	uniqueEmailError = errors.New("Email already is use")
}

func (user *User) IsUsernameUnique() bool {
	o := orm.NewOrm()
	o.Using("default")
	return !o.QueryTable("user").Filter("username", user.Username).Exist()
}

func (user *User) IsEmailUnique() bool {
	o := orm.NewOrm()
	o.Using("default")
	return !o.QueryTable("user").Filter("email", user.Password).Exist()
}

func (user *User) Create() (int, bool) {
	o := orm.NewOrm()
	o.Using("default")
	password := []byte(user.Password)
	hashedPassword, _ := bcrypt.GenerateFromPassword(password, cost)
	user.Password = string(hashedPassword)
	uid, err := o.Insert(user)
	if err == nil {
		user.Password = ""
		return int(uid), true
	}
	return 0, false
}

func (user *User) Login() bool {
	o := orm.NewOrm()
	o.Using("default")
	password := []byte(user.Password)
	err := o.QueryTable("user").Filter("username", user.Username).One(user, "uid", "password")
	if err == nil {
		if e := bcrypt.CompareHashAndPassword([]byte(user.Password), password); e == nil {
			return true
		}
	}
	return false
}

func (user *User) Update()bool {
	o := orm.NewOrm()
	o.Using("default")
	_,err := o.Update(user)
	if nil == err{
		return true
	}
	return false
}
/*
func (user *User) MakeEditor() bool {
	o := orm.NewOrm()
	o.Using("default")
	_, err := o.QueryTable("user").Filter("uid", user.Uid).Update(orm.Params{"is_editor": 1})
	if err == nil {
		return true
	}
	return false
}

func (user *User) RevokeEditor() bool {
	o := orm.NewOrm()
	o.Using("default")
	_, err := o.QueryTable("user").Filter("uid", user.Uid).Update(orm.Params{"is_editor": 0})
	if err == nil {
		return true
	}
	return false
}
*/

func (user *User) ChangePassword(password string) bool {
	o := orm.NewOrm()
	o.Using("default")
	pass := []byte(password)
	hashedPassword, _ := bcrypt.GenerateFromPassword(pass, 2)
	password = string(hashedPassword)
	_, err := o.QueryTable("user").Filter("uid", user.Uid).Update(orm.Params{"password": password})
	if err == nil {
		return true
	}
	return false
}

/*
func (user *User) GetUserInfo() bool {
	o := orm.NewOrm()
	o.Using("default")
	err := o.QueryTable("user").Filter("uid", user.Uid).One(user, "uid", "username", "nickname",
		"college", "email")
	if err == nil {
		return true
	}
	return false
}

func (user *User) Delete() bool {
	o := orm.NewOrm()
	o.Using("default")
	_, err := o.Delete(user)
	if err == nil {
		return true
	}
	return false
}

func (user *User) IsEditor() bool {
	o := orm.NewOrm()
	o.Using("default")
	err := o.Read(user)
	user.Password = ""
	if err == nil {
		if user.IsAdmin == 1 {
			return true
		}
	}
	return false
}

func (user *User) SearchByName() ([]User, int64) {
	var users []User
	o := orm.NewOrm()
	o.Using("default")
	count, err := o.Raw("select * from user where nickname like ?", "%"+user.Nickname+"%").QueryRows(&users)
	if err != nil {
		return nil, 0
	}
	return users, count
}

func (user *User) Get() bool {
	o := orm.NewOrm()
	o.Using("default")
	if err := o.Read(user); err == nil {
		return true
	}
	return false
}


func (user *User) GetByUsername() bool {
	o := orm.NewOrm()
	o.Using("default")
	if err := o.Read(user, "Username"); err == nil {
		return true
	}
	return false
}
*/
func (user *User) GetByUid() bool {
	o := orm.NewOrm()
	o.Using("default")
	if err := o.Read(user); err == nil {
		return true
	}
	return false
}

func (user *User) Sort() ([]User){
	var users []User
	o := orm.NewOrm()
	o.Using("default")
	o.QueryTable("user").All(&users)


	for _, person := range users {
		acnum,_ := o.QueryTable("solution").Filter("Uid",person.Uid).Filter("Result",1).GroupBy("Pid").Count()
		totnum,_ := o.QueryTable("solution").Filter("Uid",person.Uid).Count()
		person.AcceptsCount = int(acnum)
		person.SolutionsCount = int(totnum)
		if totnum > 0{
			person.Radio = float64(acnum * 100.0 / totnum)
		}
		o.Update(&person)
	}

	o.QueryTable("user").OrderBy("-AcceptsCount","-SolutionsCount").All(&users)
	return users
}
func (user *User) LoginVerify() error {
	var status bool
	status = CheckUserName(user.Username)
	if !status {
		return UsernameError
	}
	status = CheckPassword(user.Password)
	if !status {
		return passwordError
	}
	return nil
}

/*
func (user *User) GetEditors() []User {
	var users []User
	o := orm.NewOrm()
	o.Using("default")
	_, err := o.QueryTable("user").Filter("is_editor", 1).All(&users)
	if err != nil {
		return nil
	}
	return users
}
*/

func (user *User) SignupVerify() error {
	var status bool
	status = CheckUserName(user.Username)
	if !status {
		return UsernameError
	}
	status = CheckPassword(user.Password)
	if !status {
		return passwordError
	}
	status = CheckEmail(user.Email)
	if !status {
		return emailError
	}
	status = CheckCollege(user.College)
	if !status {
		return collegeError
	}
	status = CheckName(user.Username)
	if !status {
		return UsernameError
	}
	status = user.IsUsernameUnique()
	if !status {
		return uniqueUsernameError
	}
	status = user.IsEmailUnique()
	if !status {
		return uniqueEmailError
	}
	return nil
}

