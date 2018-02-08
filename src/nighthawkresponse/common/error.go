/*  *@package  nightHawk  *@file     error.go  *@author   roshan maskey
<roshanmaskey@gmail.com>  *  *@description  nightHawk Response error code  */


 //package nighthawk
package common

const (
    NO_ERROR                                = 0x0000
    ERROR_NO_TRIAGE_FILE                    = 0x0001
    ERROR_UNSUPPORTED_FILE_TYPE             = 0x0002
    ERROR_UNSUPPORTED_TRIAGE_FILE           = 0x0003
    ERROR_AUDIT_COMPUTERNAME_REQUIRED       = 0x0004
    ERROR_CONFIG_FILE_READ                  = 0x0005
    ERROR_READING_COMPUTERNAME              = 0x0006
    ERROR_SAME_CASE_AND_COMPUTERNAME		= 0x0007
    ERROR_READING_TRIAGE_FILE				= 0x0008
    ERROR_EXTRACTING_REDLINE_ARCHIVE		= 0x0009
    ERROR_ACCESS_REDLINE_DIRECTORY          = 0x000A
    ERROR_AUDITTYPE_INFO_PARSE              = 0x000B
    ERROR_READING_AUDIT_FILE                = 0x000C
    ERROR_WRITING_OUTPUT_FILE               = 0x000D
    ERROR_ELASTIC_INDEX_QUERY               = 0x000E
    ERROR_ELASTIC_INDEX_POST                = 0x000F

    ERROR_CHANNEL_CONNECT                   = 0x1001
    ERROR_ELASTIC_CONNECT                   = 0x1002
    ERROR_ELASTIC_CREATE_PARENT             = 0x1003
    ERROR_SPLUNK_UPLOAD                     = 0x1004
    ERROR_SPLUNK_AUTHENTICATION             = 0x0105

    // Error for Hunt
    ERROR_INVALID_MD5                       = 0x2000
    ERROR_STACK_CERT_QUERY                  = 0x2001
    ERROR_STACK_NO_AUDITTYPE                = 0x2002
    ERROR_STACK_ITEM_QUERY                  = 0x2003


)