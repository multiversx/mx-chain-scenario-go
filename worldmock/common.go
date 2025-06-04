package worldmock

// ConvertTimeStampSecToMs will convert unix timestamp from seconds to milliseconds
// TODO: this has to be handled properly when round timestamp granularity will be changed to milliseconds
func ConvertTimeStampSecToMs(timeStamp uint64) uint64 {
	return timeStamp * 1000
}
