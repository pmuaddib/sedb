package entity

// Data type used for Entrypoint
type Data []string

// Entrypoint type used for handlers
type Entrypoint interface {
    Execute(Data) error
}
