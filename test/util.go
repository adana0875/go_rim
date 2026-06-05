package test

import (
	"fmt"
	"gorim/internal/types"
	"math/rand/v2"

	"github.com/google/uuid"
)

func CreateRandomMod() types.InternalMod {

	name := fmt.Sprintf("mod_%d", rand.IntN(100))
	packid := uuid.New()
	package_name := fmt.Sprintf("mod.package_%d", packid)
	return types.InternalMod{Name: name, PackageId: package_name}
}
