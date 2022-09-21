package validator

type Validatable interface {
	Validate() error
}
