package auth

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
)

// パッケージ読み込み時にint64型をgobに登録する
// セッションにint64型のユーザーIDを保存するために必要
func init() {
	gob.Register(int64(0))
}

const (
	SessionName = "go_todo_session"
	UserKey     = "user_id"
)

type SessionManager struct {
	store *redisstore.RedisStore
}

func NewSessionManager() (*SessionManager, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
		),
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	store, err := redisstore.NewRedisStore(context.Background(), client)
	if err != nil {
		return nil, fmt.Errorf("failed to create redis store: %w", err)
	}

	// session_　プリフィックスを設定
	store.KeyPrefix("session_")

	// 環境変数でSecureフラグを制御（本番環境ではtrue）
	isSecure := os.Getenv("COOKIE_SECURE") == "true"

	// Cookieのオプションを設定
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	})

	return &SessionManager{store: store}, nil
}

func (sm *SessionManager) Get(r *http.Request) (*sessions.Session, error) {
	return sm.store.Get(r, SessionName)
}

func (sm *SessionManager) Store() *redisstore.RedisStore {
	return sm.store
}

func (sm *SessionManager) GetUserID(r *http.Request) (int64, error) {
	session, err := sm.Get(r)
	if err != nil {
		return 0, err
	}

	userID, ok := session.Values[UserKey].(int64)
	if !ok {
		return 0, fmt.Errorf("user not authenticated")
	}

	return userID, nil
}

func (sm *SessionManager) SetUserID(w http.ResponseWriter, r *http.Request, userID int64) error {
	session, err := sm.Get(r)
	if err != nil {
		return err
	}

	session.Values[UserKey] = userID
	return session.Save(r, w)
}

func (sm *SessionManager) Clear(w http.ResponseWriter, r *http.Request) error {
	session, err := sm.Get(r)
	if err != nil {
		return err
	}

	session.Values[UserKey] = nil
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
