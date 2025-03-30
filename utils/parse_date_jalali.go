// +build jalali

package utils
import (
	"fmt"
	"time"

	ptime "github.com/yaa110/go-persian-calendar"
)
func parseYmd(s string, loc *time.Location) (time.Time, error){
  var y,m,d int
  _, err := fmt.Sscanf(s, "%d-%d-%d", &y, &m, &d)
  
  pdate := ptime.Date(y, ptime.Month(m), d, 0, 0, 0, 0, loc)
  return pdate.Time(), err
}
