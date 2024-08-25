package sqlstorage

import (
	"database/sql"
	"fmt"
	"medods/interlan/app/model"
	"medods/storage"
	"net/smtp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// RefreshTokenRepository ...
type RefreshTokenRepository struct {
	storage *Storage
}

type MyClaims struct {
	jwt.RegisteredClaims
	Ip string `json:"ip"`
}

const SecretKey = "525348e77144a9cee9a7471a8b67c50ea85b9e3eb377a3c1a3a23db88f9150eefe76e6a339fdbc62b817595f53d72549d9ebe36438f8c2619846b963e9f43a94"

// CreateTokens ...
func (r *RefreshTokenRepository) CreateTokens(userid int64, ip string, email string) (*model.Auth, error) {

	encRefreshToken, err := encryptString(ip)
	if err != nil {
		return nil, err
	}

	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS512, &MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 10)), // 30d
		},
		Ip: ip,
	})
	accessToken, err := accessClaims.SignedString([]byte(SecretKey))
	if err != nil {
		return nil, err
	}

	tokens := &model.Auth{
		RefreshToken: encRefreshToken,
		AccessToken:  accessToken,
	}

	r.storage.db.QueryRow(
		"INSERT INTO refresh_tokens (userid, email, refresh_token) VALUES ($1, $2, $3)",
		userid,
		email,
		encRefreshToken,
	)

	return tokens, nil
}

// RefreshAccess ...
func (r *RefreshTokenRepository) RefreshAccess(userid int64, ip string, refresh_token string) (*model.Auth, error) {
	user, err := r.FindByUserId(userid)
	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.RefreshToken), []byte(ip)) != nil {
		to := user.Email
		from := "Smolentcev.Aleksandr@yandex.ru"
		subject := "Вход с нового устройства"
		body := fmt.Sprintf("Мы заметили, что в ваш аккаунт --- выполнен вход с ip: %s. Если это вы, ничего делать не нужно. А если нет, то мы поможем вам защитить аккаунт.", ip)

		err := sendEmail(to, from, subject, body)
		if err != nil {
			return nil, err
		}
	}

	encRefreshToken, err := encryptString(ip)
	if err != nil {
		return nil, err
	}

	r.storage.db.QueryRow(
		"UPDATE refresh_tokens SET refresh_token = $1 WHERE userid = $2;",
		encRefreshToken,
		userid,
	)

	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS512, &MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 10)), // 30d
		},
		Ip: ip,
	})
	accessToken, err := accessClaims.SignedString([]byte(SecretKey))
	if err != nil {
		return nil, err
	}

	tokens := &model.Auth{
		RefreshToken: encRefreshToken,
		AccessToken:  accessToken,
	}
	return tokens, nil
}

// FindByUserId ...
func (r *RefreshTokenRepository) FindByUserId(userid int64) (*model.User, error) {
	u := &model.User{}
	if err := r.storage.db.QueryRow(
		"SELECT userid, email, refresh_token FROM refresh_tokens WHERE userid = $1",
		userid,
	).Scan(
		&u.UserId,
		&u.Email,
		&u.RefreshToken,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func sendEmail(to, from, subject, body string) error {
	server := "smtp.yandex.ru"
	port := "587"

	auth := smtp.PlainAuth("", from, "ayhpkxikeubythnm", server)

	err := smtp.SendMail(server+":"+port, auth, from, []string{to}, []byte(body))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
