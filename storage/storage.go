package storage

// Storage ...
type Storage interface {
	Auth() RefreshTokenRepository
}
