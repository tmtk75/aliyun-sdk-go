package rds

const (
	MySQL_NetworkTraffic       = "MySQL_NetworkTraffic"
	MySQL_QPSTPS               = "MySQL_QPSTPS"
	MySQL_Sessions             = "MySQL_Sessions"
	MySQL_InnoDBBufferRatio    = "MySQL_InnoDBBufferRatio"
	MySQL_InnoDBDataReadWriten = "MySQL_InnoDBDataReadWriten"
	MySQL_InnoDBLogRequests    = "MySQL_InnoDBLogRequests"
	MySQL_InnoDBLogWrites      = "MySQL_InnoDBLogWrites"
	MySQL_TempDiskTableCreates = "MySQL_TempDiskTableCreates"
	MySQL_MyISAMKeyBufferRatio = "MySQL_MyISAMKeyBufferRatio"
	MySQL_MyISAMKeyReadWrites  = "MySQL_MyISAMKeyReadWrites"
	MySQL_COMDML               = "MySQL_COMDML"
	MySQL_RowDML               = "MySQL_RowDML"
	MySQL_MemCpuUsage          = "MySQL_MemCpuUsage"
	MySQL_IOPS                 = "MySQL_IOPS"
	MySQL_SpaceUsage           = "MySQL_SpaceUsage"
)

var DBInstancePerformanceKeys = []string{
	MySQL_NetworkTraffic,
	MySQL_QPSTPS,
	MySQL_Sessions,
	MySQL_InnoDBBufferRatio,
	MySQL_InnoDBDataReadWriten,
	MySQL_InnoDBLogRequests,
	MySQL_InnoDBLogWrites,
	MySQL_TempDiskTableCreates,
	MySQL_MyISAMKeyBufferRatio,
	MySQL_MyISAMKeyReadWrites,
	MySQL_COMDML,
	MySQL_RowDML,
	MySQL_MemCpuUsage,
	MySQL_IOPS,
	MySQL_SpaceUsage,
}
