package dao

func buildFindQuery(query string, conditions ...interface{}) (string, []interface{}) {
	var queryArgs []interface{}
	where := "WHERE "
	argLength := len(conditions)
	switch {
	case argLength == 1:
		where = "id = ?"
		queryArgs = conditions
	case argLength%2 == 0:
		for i, a := range conditions {
			if i%2 == 0 {
				if i > 0 {
					where += "AND "
				}
				where += a.(string) + " = ? "
			} else {
				queryArgs = append(queryArgs, a)
			}
		}
	default:
		var m any = "invalid number of arguments are passed"
		panic(m)
	}
	return query + " " + where, queryArgs
}
