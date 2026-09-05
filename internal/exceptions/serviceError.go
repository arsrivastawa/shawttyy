package exceptions

const ErrNodeOutOfRange Error = "sequencer: node id out of range"
const ErrInvalidOriginalURL Error = "service: original_url is invalid"
const ErrInvalidCustomAlias Error = "service: custom_alias must be alphanumeric and <= 11 chars"
const ErrCustomAliasTaken Error = "service: custom_alias is already in use"
const ErrShortURLNotFound Error = "service: short_url not found"
const ErrURLExpired Error = "service: short_url has expired"
