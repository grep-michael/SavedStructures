package savedstructures

type SaveableInterface interface {
	Load() error
	Save() error
}
