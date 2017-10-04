/* nighthawk.nhstruct.route.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Datastructure for Network route
 */

package nhstruct


type RouteIpv6Item struct {
    ValidLifeTime       string
    PreferredLifeTime   string
    IsLoopback          bool
    IsAutoconfiguredAddress bool
    IsPublished         bool
    IsImortal           bool
    Origin              string
} 

type RouteEntryItem struct {
    JobCreated          string `xml:"created, attr"`
    IsIpv6              bool `xml:"IsIPv6"`
    //Ipv4Route
    Interface           string 
    Destination         string
    Netmask             string
    RouteType           string
    //Common for Ipv4 and Ipv6 route
    Gateway             string
    Protocol            string
    RouteAge            string
    Metric              int
    // IPv6 Route
    ValidLifeTime       string
    PreferredLifeTime   string
    IsLoopback          bool
    IsAutoconfiguredAddress bool
    IsPublished         bool
    IsImortal           bool
    Origin              string
    IsWhitelisted       bool
    NHScore             int
    Tag                 string
    NhComment           NHComment `json:"Comment"`
}

