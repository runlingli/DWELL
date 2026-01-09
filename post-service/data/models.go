package data

import (
	"context"    
	"database/sql" 
	"errors"      
	"log"    
	"time"    
	"golang.org/x/crypto/bcrypt"

const dbTimeout = time.Second * 3 
var db *sql.DB


func New(dbPool *sql.DB) Models {
	db = dbPool 
	return Models{
		User: User{}, 
	}
}

type Models struct {
	Listing Listing
}

type Listing struct {
	ID               int       `json:"id" db:"id"`                     
	Title            string    `json:"title" db:"title"`               
	Price            float64   `json:"price" db:"price"`                
	Location         string    `json:"location" db:"location"`
	Neighborhood     string    `json:"neighborhood" db:"neighborhood"` 
	Coordinates      GeoPoint  `json:"coordinates" db:"coordinates"`   
	Radius           int       `json:"radius" db:"radius"`              
	Type             string    `json:"type" db:"type"`                 
	ImageURL         string    `json:"imageUrl" db:"image_url"`           
	AdditionalImages []string  `json:"additionalImages" db:"additional_images"`
	Description      string    `json:"description" db:"description"`     
	Bedrooms         int       `json:"bedrooms" db:"bedrooms"`          
	Bathrooms        int       `json:"bathrooms" db:"bathrooms"`       
	CreatedAt        time.Time `json:"createdAt" db:"created_at"`       
	AvailableFrom    time.Time `json:"availableFrom" db:"available_from"`
	AvailableTo      time.Time `json:"availableTo" db:"available_to"`   
	AuthorID         int       `json:"authorId" db:"author_id"` 
}

// 坐标结构
type GeoPoint struct {
	Lat float64 `json:"lat" db:"lat"`
	Lng float64 `json:"lng" db:"lng"`
}


// 作者信息
type Author struct {
	Name   string `json:"name" db:"name"`
	Avatar string `json:"avatar" db:"avatar"`
}


// GetAll 返回所有用户，按 LastName 排序
func (u *User) GetAll() ([]*User, error) {
	// 创建一个上下文 ctx，用于控制数据库操作超时
	// context.WithTimeout 会返回 ctx 和 cancel 函数
	// ctx 会在 3 秒后自动取消数据库查询
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel() // 函数结束时释放资源，防止泄漏

	query := `select id, coalesce(google_id, ''), email, first_name, last_name, password, user_active, created_at, updated_at
	from users order by last_name`

	// QueryContext 用 ctx 执行 SQL 查询
	// 返回多个结果 rows
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // 使用完 rows 后关闭

	var users []*User

	// 循环读取每一行结果
	for rows.Next() {
		var user User
		// Scan 将每列值赋给 user 对应字段
		err := rows.Scan(
			&user.ID,
			&user.GoogleID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user) // 添加到结果切片
	}

	return users, nil
}

func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, coalesce(google_id, ''), email, first_name, last_name, password, user_active, created_at, updated_at from users where email = $1`

	var user User
	// QueryRowContext 返回单行数据
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.GoogleID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// =====================
// 根据 ID 获取用户
// =====================
func (u *User) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, coalesce(google_id, ''), email, first_name, last_name, password, user_active, created_at, updated_at from users where id = $1`

	var user User
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.GoogleID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) GetByGoogleID(googleID string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, google_id, email, first_name, last_name, coalesce(password, ''), user_active, created_at, updated_at from users where google_id = $1`
	var user User
	row := db.QueryRowContext(ctx, query, googleID)

	err := row.Scan(
		&user.ID,
		&user.GoogleID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error getting user by Google ID %s: %v", googleID, err)
		return nil, err
	}

	return &user, nil
}

// =====================
// 更新用户
// =====================
func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update users set
		email = $1,
		google_id = $2,
		first_name = $3,
		last_name = $4,
		user_active = $5,
		updated_at = $6
		where id = $7`

	_, err := db.ExecContext(ctx, stmt,
		u.Email,
		u.GoogleID,
		u.FirstName,
		u.LastName,
		u.Active,
		time.Now(),
		u.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// =====================
// 删除用户
// =====================
func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// =====================
// 插入新用户
// =====================
func (u *User) Insert(user User) (int, error) {
	//context.Background()返回一个 空的根 context，没有取消信号和超时。
	//context.WithTimeout(parent, timeout)基于 parent 创建一个新的 context，带有 超时功能。
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// bcrypt 加密用户密码，不可逆的加密。常用12-14，数字越大加密时间越长，越复杂
	var hashedPassword []byte
	var err error

	if user.Password != "" {
		// 仅当密码不为空时才加密
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if err != nil {
			return 0, err
		}
	} else {
		hashedPassword = []byte("") // 或者 []byte("")，数据库字段允许 NULL
	}

	var newID int
	stmt := `insert into users (email, google_id, first_name, last_name, password, user_active, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	// 返回新插入行的 ID
	//queryrowcontext以具体的值替换占位符$1...
	err = db.QueryRowContext(ctx, stmt,
		user.Email,
		user.GoogleID,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.Active,
		time.Now(),
		time.Now(),
	).Scan(&newID)
	//Scan将returning id结果储存到newID

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// =====================
// 重置密码
// =====================
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `update users set password = $1 where id = $2`
	_, err = db.ExecContext(ctx, stmt, hashedPassword, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// =====================
// 校验密码是否匹配
// =====================
func (u *User) PasswordMatches(plainText string) (bool, error) {
	// bcrypt 比较明文密码与数据库存储的哈希密码
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// 密码不匹配
			return false, nil
		default:
			// 其他错误
			return false, err
		}
	}

	return true, nil
}
