package interfaces

type (
	Backupable interface {
		Equal(other any) bool
		Parse(data []byte) error
		ParseNew(data []byte) (Backupable, error)
		Save(fileName string) error
		BackupPrefix() string
	}
)
