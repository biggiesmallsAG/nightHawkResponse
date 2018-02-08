/* nighthawk.nhstruct.disk.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Datastructure for Disk and Volume
 */

package nhstruct


/* __start_of_w32disks__ */
type PartitionItem struct {
    PartitionNumber     int 
    PartitionOffset     int
    PartitionLength     int
    PartitionType       string
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}

type DiskItem struct {
    JobCreated          string `xml:"created,attr"`
    DiskName            string
    DiskSize            string
    PartitionList       []PartitionItem `xml:"PartitionList>Partition"`
    IsWhitelisted       bool
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}



/* __start_of_w32volumes__ */
type VolumeItem struct {
    JobCreated                      string `xml:"created,attr"`
    VolumeName                      string 
    DevicePath                      string
    Type                            string
    Name                            string
    SerialNumber                    string
    FileSystemFlags                 string
    FileSystemName                  string
    ActualAvailableAllocationUnits  int
    TotalAllocationUnits            int
    BytesPerSector                  int
    SectorsPerAllocationUnit        int
    CreationTime                    string
    IsMounted                       bool
    IsWhitelisted                   bool
    Tag                             string
    NhComment                       NHComment `json:"Comment"`
}


