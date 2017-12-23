/* nighthawk.nhstruct.modules.go
 * author: roshan maskey <roshanmaskey@gmail.com>
 *
 * Datastructure for Modules
 * Collected as w32drivers-modulelist
 */
package nhstruct

type ModuleItem struct {
	ModuleAddress int64
	ModuleInit int64 
	ModuleBase int64
	ModuleSize int
	ModulePath string
	ModuleName string
	IsGoodEntry string
	IsWhitelisted bool
	NHScore	int
	Tag string
	NhComment NHComment `json:"Comment"`
}