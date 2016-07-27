/*
 *@package  nightHawk
 *@file     error.go
 *@author   roshan maskey <roshanmaskey@gmail.com>
 *
 *@description  nightHawk Response error code
 */


 package nightHawk

const (
    NO_ERROR                                = 0
    ERROR_NO_TRIAGE_FILE                    = 1
    ERROR_UNSUPPORTED_FILE_TYPE             = 2
    ERROR_UNSUPPORTED_TRIAGE_FILE           = 3
    ERROR_AUDIT_COMPUTERNAME_REQUIRED       = 4
    ERROR_CONFIG_FILE_READ                  = 5
    ERROR_READING_COMPUTERNAME              = 6
    ERROR_SAME_CASE_AND_COMPUTERNAME		= 7
    ERROR_READING_TRIAGE_FILE				= 8
    ERROR_EXTRACTING_REDLINE_ARCHIVE		= 9
    ERROR_ACCESS_REDLINE_DIRECTORY          = 10
)