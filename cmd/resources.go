package cmd

type ResourceInterface interface {
	Validate() error
}
