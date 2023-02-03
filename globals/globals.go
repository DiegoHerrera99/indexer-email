package globals

const TEMPDIR string = "./temp"
const ZINC_INDEX string = "enron"
const ZINC_USER string = "admin"
const ZINC_PWD string = "Complexpass#123"
const ZINC_ENDPOINT string = "http://localhost:4080/api/_bulkv2"
const ZINC_CRTIDX_ENDPOINT string = "http://localhost:4080/api/index"
const ZINC_IDXMAP string = `{
    "name": "enron",
    "storage_type": "disk",
    "shard_num": 3,
    "mappings": {
        "properties": {
            "date": {
                "type": "date",
                "index": true,
                "store": false,
                "sortable": true,
                "aggregatable": true,
                "highlightable": false
            },
            "body": {
                "type": "text",
                "index": true,
                "store": false,
                "sortable": false,
                "aggregatable": false,
                "highlightable": false
            },
            "cc": {
                "type": "text",
                "index": true,
                "store": false,
                "sortable": false,
                "aggregatable": false,
                "highlightable": false
            },
            "from": {
                "type": "text",
                "index": true,
                "store": false,
                "sortable": false,
                "aggregatable": false,
                "highlightable": false
            },
            "subject": {
                "type": "text",
                "index": true,
                "store": false,
                "sortable": false,
                "aggregatable": false,
                "highlightable": false
            },
            "to": {
                "type": "text",
                "index": true,
                "store": false,
                "sortable": false,
                "aggregatable": false,
                "highlightable": false
            }
        }
    }
}`
