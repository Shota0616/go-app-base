package auth

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"os"
	"log"
)

// JWTの署名に使用する秘密鍵を環境変数から取得
var jwtKey = []byte(os.Getenv("JWT_SECRET"))
// リフレッシュトークン用の秘密鍵
var jwtrefreshKey = []byte(os.Getenv("JWT_REFRESH_SECRET"))


// JWTのペイロードに含まれるクレーム(情報)を定義
type Claims struct {
	ID uint `json:"id"` // ユーザーID
	jwt.StandardClaims // 標準のクレーム(例: exp、iatなど)
}

// JWTトークンを生成する関数
func GenerateJWT(id uint) (string, error) {
	// トークンの有効期限を24時間後に設定
	expirationTime := time.Now().Add(24 * time.Hour)
	// claimsオブジェクトを生成
	claims := &Claims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // 有効期限
		},
	}

	// HS256アルゴリズムを使ってヘッダーとペイロードを作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey) // トークンを署名して返す
}

// JWTトークンを検証する関数
func ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{} // 検証結果を格納するためのClaims構造体
	log.Println(tokenStr)
	log.Println(claims)
	// トークンの解析と署名の検証を行い、結果をclaimsに格納
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil // トークンの署名を検証するための秘密鍵を返す
	})
	log.Println(token)
	log.Println(jwtKey)
	if err != nil {
		log.Printf("JWT parse error: %v", err)
		return nil, err // エラーがあればログに出力し、エラーを返す
	}
	if !token.Valid {
		log.Printf("Invalid token: %v", err)
		return nil, err // トークンが無効な場合はエラーを返す
	}
	return claims, nil // 有効な場合はクレームを返す
}

// リフレッシュトークンを生成する関数
func GenerateRefreshToken(id uint) (string, error) {
	// リフレッシュトークンの有効期限を7日後に設定
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // 有効期限
		},
	}

	// HS256アルゴリズムを使ってリフレッシュトークンを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtrefreshKey) // トークンを署名して返す
}

// リフレッシュトークンを使用して新しいJWTを生成する関数
func RefreshJWT(refreshTokenStr string) (string, error) {
	claims := &Claims{} // 検証結果を格納するためのClaims構造体
	// リフレッシュトークンの解析と署名の検証を行い、結果をclaimsに格納
	token, err := jwt.ParseWithClaims(refreshTokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtrefreshKey, nil // トークンの署名を検証するための秘密鍵を返す
	})
	if err != nil {
		return "", err // エラーがあればエラーを返す
	}
	if !token.Valid {
		return "", err // トークンが無効な場合はエラーを返す
	}

	// トークンの有効期限が30秒以内の場合に新しいJWTを発行
	if time.Until(time.Unix(claims.ExpiresAt, 0)) < 30*time.Second {
		return GenerateJWT(claims.ID)
	}
	return "", nil // 新しいトークンが必要ない場合は空文字を返す
}
