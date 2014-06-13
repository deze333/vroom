package api_http

import (
    "strings"
)

//------------------------------------------------------------
// Utils for routes
//------------------------------------------------------------

// Generates all logical equals for given route.
func routeSynonyms(r string) (rs []string) {

    // Ignore root path
    if r == "/" {
        return []string{r}
    }

    // If: /hello/there/
    if strings.HasSuffix(r, "/") {
        return []string{r[0:len(r)], r}
    }

    // If: /hello/there
    return []string{r, r + "/"}
}

