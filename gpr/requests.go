package gpr

import (
	"net/http"
)

func (g *GPR) PerformSubmitLevelRequest(JSON []byte, expectedStatusCode int, target interface{}) {
	g.PerformRequest("/submit", http.MethodPost, JSON, expectedStatusCode, target)
}
