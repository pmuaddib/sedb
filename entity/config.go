package entity

// Values type used for Config interface
type Values map[string]string

// Config interface for storing settings
type Config interface {
    GetConfig() Values
    SetConfig(Values)
}
