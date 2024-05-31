package common

type TicketInterface interface {
	TicketFilename() string
	TicketToYaml() []byte
}
