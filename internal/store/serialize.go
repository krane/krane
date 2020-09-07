package store

import "encoding/json"

func Serialize(value interface{}) ([]byte, error) { return json.Marshal(value) }

func Deserialize(bytes []byte, struct0 interface{}) error { return json.Unmarshal(bytes, struct0) }
