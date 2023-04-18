package backup

import "context"

// IAutoBackup - интерфейс для реализации автоматической выгрузки и восстановления состояния
type IAutoBackup interface {
	// Start -
	Start(ctx context.Context) chan struct{}
	// Restore -
	Restore(ctx context.Context) error
	// Backup -
	Backup(ctx context.Context) error
}
