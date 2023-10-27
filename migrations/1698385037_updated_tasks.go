package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("by7d8mdliry8zi5")
		if err != nil {
			return err
		}

		// update
		edit_owner := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "ihs8ldz8",
			"name": "owner",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "_pb_users_auth_",
				"cascadeDelete": true,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), edit_owner)
		collection.Schema.AddField(edit_owner)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("by7d8mdliry8zi5")
		if err != nil {
			return err
		}

		// update
		edit_owner := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "ihs8ldz8",
			"name": "owner",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "_pb_users_auth_",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), edit_owner)
		collection.Schema.AddField(edit_owner)

		return dao.SaveCollection(collection)
	})
}
