package error

type Error string

func (e Error) Error() string { return string(e) }

const ErrNodeOutOfRange Error = "sequencer: node id out of range"
