package vroom

import (
    "github.com/deze333/vroom/api_http"
)

//------------------------------------------------------------
// Initialize
//------------------------------------------------------------

// Initializes API HTTP dispatcher.
func Initialize(ctx api_http.Ctx) (err error) {
    return api_http.Initialize(ctx)
}

