package yaorm_test

import (
	"context"
	"testing"

	"github.com/geoffreybauduin/yaorm"
	"github.com/geoffreybauduin/yaorm/testdata"
	"github.com/stretchr/testify/assert"
)

func TestSqlExecutor_Exec(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.NoError(t, err)
	dbp, err := yaorm.NewDBProvider(context.TODO(), "test")
	assert.NoError(t, err)
	_, err = dbp.DB().Exec("SELECT 1")
	assert.NoError(t, err)
}
