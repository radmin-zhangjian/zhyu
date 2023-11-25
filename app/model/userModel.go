package model

type User struct {
	Id       int64  `gorm:"primaryKey;column:id;" json:"id"`
	Username string `gorm:"column:user_name;type:varchar(30);default:(-)" json:"user_name"`
	Password string `gorm:"column:password;type:varchar(100);default:(-)" json:"password"`
	Phone    string `gorm:"column:phone;type:varchar(20);default:(-)" json:"phone"`
	Name     string `gorm:"column:name;type:varchar(30);default:(-)" json:"name"`
	Age      uint8  `gorm:"column:age;type:tinyint(1);default:(0)" json:"age"`
	Address  string `gorm:"column:address;type:varchar(100);default:(-)" json:"address"`
	Photo    string `gorm:"column:photo;type:varchar(100);default:(-)" json:"photo"`
	Status   uint8  `gorm:"column:status;type:tinyint(1);default:(1)" json:"status"`
	//Deleted    gorm.DeletedAt `gorm:"column:deleted;type:timestamp;default:(-)" json:"deleted"`
	CreatedAt int64 `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt int64 `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}

// TableName 自定义表名
func (*User) TableName() string {
	return "zhyu_user"
}
