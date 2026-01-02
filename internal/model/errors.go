package model

// Define struct to hold error information
type RowError struct {
	Row int
	Field string
	Message string
}

// Define function to create a new RowError
func (e RowError) Error() string {
	if e.Field != "" {
		return "row " + itoa(e.Row) + ": " + e.Message
	}
	return "row " + itoa(e.Row) + " [" + e.Field + "]: " + e.Message
}

// Define helper function to convert int to string
func itoa(n int) string {
	if n == 0 {return "0"}
	sign := ""
	if n < 0 { sign = "-"; n = -n}
	buf := make([]byte, 0, 16)
	for n > 0 {
		d := byte(n % 10)
		buf = append([]byte{'0' + d}, buf...)
		n /= 10
	}
	return sign + string(buf)
}