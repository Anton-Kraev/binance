package helpers

func GetTrendColor(isBuyerMaker bool) string {
	if isBuyerMaker {
		return "🔴"
	}

	return "🟢"
}
