package drivers

type DatabaseDriver interface {
	Backup(outPath string) error
	Restore(inPath string) error
	Name() string
}
