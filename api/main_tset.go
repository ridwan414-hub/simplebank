package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// func newTestServer(t *testing.T, store db.Store) *Server {
// 	config := utils.Config{
// 		TokenSymmetricKey:   utils.RandomString(32),
// 		AccessTokenDuration: time.Minute,
// 	}
// 	server, err := NewServer(config, &store)
// 	require.NoError(t, err)
// 	return server
// }

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}