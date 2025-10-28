package exchange

import (
	"fmt"
	"testing"
)

func TestTimeRange_Minutes(t *testing.T) {
	tr := ExchangeTime(Trading, "09:30:00", "11:30:00")
	minutes := tr.Minutes()
	fmt.Println(minutes)
}
