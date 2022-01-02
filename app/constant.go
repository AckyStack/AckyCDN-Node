package app

const PREFIX_VHOST = "<vhostutils>"

const (
	ErrSystemInternal            = 1001 //System Internal Error
	ErrRestricted                = 1002 //Restricted
	ErrAccessDeniedDIANA         = 1003 //Access Denied: Direct IP Access Not Allowed
	ErrHostNotConfigured         = 1004 //Host Not Configured to Serve Web Traffic
	ErrAccessDeniedIP            = 1005 //Access Denied: Your IP address has been banned
	ErrAccessDeniedCountryRegion = 1006 //Access Denied: Country or region banned
	ErrAccessDeniedOBUA          = 1007 //Access Denied: The owner of this website has banned your access based on your browser's signature (User Agent)
)
