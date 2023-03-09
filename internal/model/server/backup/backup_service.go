package backup

// IAutoBackup - интерфейс для реализации автоматической выгрузки и восстановления состояния
type IAutoBackup interface {
	// Start -
	Start() chan struct{}
	// Restore -
	Restore() error
	// Backup -
	Backup() error
}
