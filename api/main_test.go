package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/simple_bank/config"
	"github.com/simple_bank/database"
	"github.com/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, db database.Database) *Server {
	testConfig := config.Config{
		Token: config.TokenConfig{
			SymmetricKey: util.RandomString(32),
			Duration:     time.Minute,
		},
	}

	server, err := NewServer(db, nil, testConfig)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
