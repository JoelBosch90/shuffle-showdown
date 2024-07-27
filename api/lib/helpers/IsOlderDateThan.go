package helpers

func IsOlderDateThan(yearA uint, monthA uint, dayA uint, yearB uint, monthB uint, dayB uint) bool {
	if yearA < yearB {
		return true
	}

	if yearA == yearB && monthA < monthB {
		return true
	}

	if yearA == yearB && monthA == monthB && dayA < dayB {
		return true
	}

	return false
}
