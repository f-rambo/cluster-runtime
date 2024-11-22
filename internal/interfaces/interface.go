package interfaces

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewLogInterface, NewClusterInterface, NewAppInterface, NewProjectInterface, NewServiceInterface, NewUserInterface)