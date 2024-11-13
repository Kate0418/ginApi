package migrations

import (
	"back/api/models"
)

var Models = map[string]interface{}{
    "test": &models.Test{},
}

func Migrate() {

}
