package request

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"encoding/json"
)

func TestUpdateJson(t *testing.T) {
	body := "{\"id\":\"sessionId\",\"data\":{\"key\":\"value\"}}"
	data := Update{}
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		t.Fatal("Json decoding failed: " + err.Error())
	}

	dataFixture := make(map[string]string)
	dataFixture["key"] = "value"
	assert.Equal(t, "sessionId", data.Id, "Id was not correct parsed from JSON")
	assert.Equal(t, dataFixture, data.Data, "Data was not correct parsed from JSON")
}
