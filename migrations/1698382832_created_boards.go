package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "d5tpeg8jtjkexq1",
			"created": "2023-10-27 05:00:32.425Z",
			"updated": "2023-10-27 05:00:32.425Z",
			"name": "boards",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "3j6larjv",
					"name": "name",
					"type": "text",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				}
			],
			"indexes": [],
			"listRule": null,
			"viewRule": null,
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("d5tpeg8jtjkexq1")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
