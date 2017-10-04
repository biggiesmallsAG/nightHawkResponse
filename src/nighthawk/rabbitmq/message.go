/* nigthhawk.rabbitmq.message
 * author: 0xredskull
 *
 * Contains message structures passed using RabbitMQ
 * among nighthawk worker and other components
 */

package rabbitmq

type JobMessage struct {
	CaseName     string
	CaseDate     string
	ComputerName string
	CaseAnalyst  string
	TriageFile   string
}
