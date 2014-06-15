package api_http

import ()

//------------------------------------------------------------
// Utils for event handling
//------------------------------------------------------------

// Current app deployemnt version
var _appVersion string


// Sets initial App Version value.
func setAppVersion(ver string) {
    _appVersion = ver
}

// Callback function that sets current App Version variable.
func onAppVersionChanged(ver string) {
    _appVersion = ver
}

// Getter for the current App Version.
func GetVersion() string {
    return _appVersion
}
