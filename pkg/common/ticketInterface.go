package common

// TicketInterface is an interface that all tickets must implement, used in
// commiting with the repo package.
type TicketInterface interface {
	TicketFilename() string
	TicketToYaml() []byte
}
